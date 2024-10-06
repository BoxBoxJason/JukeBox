package static

import (
	"fmt"
	"os"
	"path"
)

var (
	JUKEBOX_DIR = path.Join(os.Getenv("HOME"), ".jukebox")
)

func init() {
	if _, err := os.Stat(JUKEBOX_DIR); os.IsNotExist(err) {
		os.Mkdir(JUKEBOX_DIR, os.ModePerm)
	} else if err != nil {
		fmt.Println("Failed to create Jukebox directory:", err)
		os.Exit(1)
	}
}

// BuildJukeboxPath returns the absolute path to a file in the Jukebox directory.
func BuildJukeboxPath(relative_path string) string {
	return path.Join(JUKEBOX_DIR, relative_path)
}
