package constants

import (
	"fmt"
	"os"
	"path"
)

type contextKey string

var (
	JUKEBOX_PATH         = path.Join(os.Getenv("HOME"), ".jukebox")
	DB_DIR               = path.Join(JUKEBOX_PATH, "db")
	DB_FILE              = path.Join(DB_DIR, "jukebox.db")
	DB_BACKUP_DIR        = path.Join(DB_DIR, "backup")
	LOG_DIR              = path.Join(JUKEBOX_PATH, "logs")
	TOKEN_EXPIRATION_MAP = map[string]int64{
		ACCESS_TOKEN:  ACCESS_TOKEN_EXPIRATION,
		REFRESH_TOKEN: REFRESH_TOKEN_EXPIRATION,
	}
)

const (
	AUTH_SCHEME                          = "Bearer"
	ACCESS_TOKEN                         = "access"
	ACCESS_TOKEN_COOKIE_NAME             = "access_token"
	ACCESS_TOKEN_COOKIE_PATH             = "/"
	ACCESS_TOKEN_EXPIRATION              = 4
	ACCESS_TOKEN_CONTEXT_KEY  contextKey = "access_token"
	REFRESH_TOKEN                        = "refresh"
	REFRESH_TOKEN_COOKIE_NAME            = "refresh_token"
	REFRESH_TOKEN_COOKIE_PATH            = "/api/auth"
	REFRESH_TOKEN_EXPIRATION             = 7 * 24
	USER_CONTEXT_KEY          contextKey = "user"
)

func init() {
	if _, err := os.Stat(JUKEBOX_PATH); os.IsNotExist(err) {
		os.Mkdir(JUKEBOX_PATH, os.ModePerm)
	} else if err != nil {
		fmt.Println("Failed to create Jukebox directory:", err)
		os.Exit(1)
	}
}
