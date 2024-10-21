package constants

import (
	"fmt"
	"os"
	"path"
)

var (
	JUKEBOX_PATH  = path.Join(os.Getenv("HOME"), ".jukebox")
	DB_DIR        = path.Join(JUKEBOX_PATH, "db")
	DB_FILE       = path.Join(DB_DIR, "jukebox.db")
	DB_BACKUP_DIR = path.Join(DB_DIR, "backup")
)

const (
	AUTH_SCHEME = "Identity"
)

func init() {
	if _, err := os.Stat(JUKEBOX_PATH); os.IsNotExist(err) {
		os.Mkdir(JUKEBOX_PATH, os.ModePerm)
	} else if err != nil {
		fmt.Println("Failed to create Jukebox directory:", err)
		os.Exit(1)
	}
}

// BuildJukeboxPath returns the absolute path to a file in the Jukebox directory.
func BuildJukeboxPath(relative_path string) string {
	return path.Join(JUKEBOX_PATH, relative_path)
}
