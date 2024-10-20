package db_model

import (
	"os"

	"github.com/boxboxjason/jukebox/internal/constants"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenConnection opens a connection to the SQLite database
func OpenConnection() (*gorm.DB, error) {
	if _, err := os.Stat(constants.DB_FILE); os.IsNotExist(err) {
		os.Mkdir(constants.DB_DIR, os.ModePerm)
	} else if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(constants.DB_FILE), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CloseConnection closes the connection to the SQLite database
func CloseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
