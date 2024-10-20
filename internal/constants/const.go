package constants

import (
	"os"
	"path"
)

var (
	JUKEBOX_PATH  = path.Join(os.Getenv("HOME"), ".jukebox")
	DB_DIR        = path.Join(JUKEBOX_PATH, "db")
	DB_FILE       = path.Join(DB_DIR, "jukebox.db")
	DB_BACKUP_DIR = path.Join(DB_DIR, "backup")
)
