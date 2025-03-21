package constants

import (
	"fmt"
	"os"
	"path"
)

type contextKey string

var (
	// Current app version (replaced at build time)
	JUKEBOX_VERSION = "latest-dev"
	// Path to the Jukebox directory
	JUKEBOX_PATH = path.Join(os.Getenv("HOME"), ".jukebox")
	// Path to the Jukebox images directory
	IMAGES_DIR = path.Join(JUKEBOX_PATH, "images")
	// Path to the users avatars directory
	AVATARS_DIR = path.Join(IMAGES_DIR, "avatars")
	// Path to the Jukebox db directory
	DB_DIR = path.Join(JUKEBOX_PATH, "db")
	// Path to the Jukebox db file (if using sqlite)
	DB_FILE = path.Join(DB_DIR, "jukebox.db")
	// Path to the Jukebox db backup directory
	DB_BACKUP_DIR = path.Join(DB_DIR, "backup")
	// Path to the Jukebox logs directory
	LOG_DIR = path.Join(JUKEBOX_PATH, "logs")
	// Auth Token expiration map
	TOKEN_EXPIRATION_MAP = map[string]int64{
		ACCESS_TOKEN:  ACCESS_TOKEN_EXPIRATION,
		REFRESH_TOKEN: REFRESH_TOKEN_EXPIRATION,
	}
)

const (
	// Auth token scheme
	AUTH_SCHEME = "Bearer"
	// ==================== ACCESS TOKEN ====================
	// Access token Type constant
	ACCESS_TOKEN = "access"
	// Access token cookie name
	ACCESS_TOKEN_COOKIE_NAME = "accessToken"
	// Access token cookie path
	ACCESS_TOKEN_COOKIE_PATH = "/"
	// Access token expiration in hours
	ACCESS_TOKEN_EXPIRATION = 4
	// Access token context key (used to store/retrieve the token from the context)
	ACCESS_TOKEN_CONTEXT_KEY contextKey = "accessToken"
	// Refresh token Type constant
	// ==================== REFRESH TOKEN ====================
	REFRESH_TOKEN = "refresh"
	// Refresh token cookie name
	REFRESH_TOKEN_COOKIE_NAME = "refreshToken"
	// Refresh token cookie path
	REFRESH_TOKEN_COOKIE_PATH = "/api/auth"
	// Refresh token expiration in hours
	REFRESH_TOKEN_EXPIRATION = 7 * 24
	// User context key (used to store/retrieve the user from the context)
	USER_CONTEXT_KEY contextKey = "user"
	// ==================== BAN ====================
	// Ban Type constant
	BAN_TYPE  = "ban"
	MUTE_TYPE = "mute"
	// ==================== REQUESTS PARAMETERS ====================
	ID_PARAMETER               = "id"
	USERNAME_PARAMETER         = "username"
	USER_ID_PARAMETER          = "user_id"
	PARTIAL_USERNAME_PARAMETER = "partial_username"
	EMAIL_PARAMETER            = "email"
	PASSWORD_PARAMETER         = "password"
	SUBSCRIBER_TIER_PARAMETER  = "subscriber_tier"
	ADMIN_PARAMETER            = "admin"
	REASON_PARAMETER           = "reason"
	AVATAR_PARAMETER           = "avatar"
	TARGET_ID_PARAMETER        = "target_id"
	DURATION_PARAMETER         = "duration"
	ISSUER_ID_PARAMETER        = "issuer_id"
	ENDS_AFTER_PARAMETER       = "ends_after"
	SENDER_ID_PARAMETER        = "sender_id"
	FLAGGED_PARAMETER          = "flagged"
	CENSORED_PARAMETER         = "censored"
	REMOVED_PARAMETER          = "removed"
	CONTAINS_PARAMETER         = "contains"
	MESSAGE_PARAMETER          = "message"
	ORDER_PARAMETER            = "order"
	LIMIT_PARAMETER            = "limit"
	PAGE_PARAMETER             = "page"
	OFFSET_PARAMETER           = "offset"
)

func init() {
	err := os.MkdirAll(JUKEBOX_PATH, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create Jukebox directory:", err)
		os.Exit(1)
	}
}
