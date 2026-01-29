import sqlite3

def get_connection(path="syspulse.db"):
    return sqlite3.connect(path)
