package collector

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type StartupEntry struct {
	Name      string
	Command   string
	Source    string
	Location  string
	Timestamp time.Time
}

func CollectStartupEntries() ([]StartupEntry, error) {
	entries := make([]StartupEntry, 0)

	for _, item := range []struct {
		source string
		key    string
	}{
		{source: "HKCU Run", key: `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`},
		{source: "HKLM Run", key: `HKLM\Software\Microsoft\Windows\CurrentVersion\Run`},
	} {
		output, err := exec.Command("reg", "query", item.key).CombinedOutput()
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(output), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "HKEY_") {
				continue
			}

			if strings.Contains(trimmed, "REG_") {
				name, command := parseRegistryLine(trimmed)
				if name == "" {
					continue
				}
				entries = append(entries, StartupEntry{
					Name:      name,
					Command:   command,
					Source:    item.source,
					Location:  item.key,
					Timestamp: time.Now(),
				})
			}
		}
	}

	startupDir, err := os.UserConfigDir()
	if err == nil {
		startupPath := filepath.Join(startupDir, "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
		files, readErr := os.ReadDir(startupPath)
		if readErr == nil {
			for _, entry := range files {
				if entry.IsDir() {
					continue
				}
				entries = append(entries, StartupEntry{
					Name:      entry.Name(),
					Command:   filepath.Join(startupPath, entry.Name()),
					Source:    "Startup Folder",
					Location:  startupPath,
					Timestamp: time.Now(),
				})
			}
		}
	}

	return entries, nil
}

func parseRegistryLine(line string) (string, string) {
	segments := strings.SplitN(line, "REG_", 2)
	if len(segments) < 2 {
		return "", ""
	}

	name := strings.TrimSpace(segments[0])
	rest := "REG_" + segments[1]
	parts := strings.Fields(rest)
	if len(parts) < 2 {
		return "", ""
	}

	command := strings.Join(parts[1:], " ")
	return name, command
}
