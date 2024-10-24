package db_model

import (
	"os"

	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initialize the database connection and create the tables
func CreateTables() {
	db, err := OpenConnection()
	if err != nil {
		logger.Fatal("Failed to open the database connection")
	}
	defer CloseConnection(db)
	logger.Info("Creating tables")
	if !db.Migrator().HasTable(&User{}) || !db.Migrator().HasTable(&AuthToken{}) || !db.Migrator().HasTable(&Message{}) {
		err = db.AutoMigrate(&User{}, &AuthToken{}, &Message{})
		if err != nil {
			logger.Fatal("Failed to create tables:", err)
		} else {
			logger.Info("Tables created successfully")
		}
	}
}

// OpenConnection opens a connection to the SQLite database
func OpenConnection() (*gorm.DB, error) {
	if _, err := os.Stat(constants.DB_DIR); os.IsNotExist(err) {
		err = os.Mkdir(constants.DB_DIR, 0700)
		if err != nil {
			logger.Critical("Failed to create the database directory", err)
			return nil, httputils.NewDatabaseError("Failed to create the database directory")
		}
	} else if err != nil {
		logger.Critical("Failed to check if the database file exists:", err)
		return nil, httputils.NewDatabaseError("Failed to check if the database file exists")
	}

	db, err := gorm.Open(sqlite.Open(constants.DB_FILE), &gorm.Config{})
	if err != nil {
		logger.Critical("Failed to open the database connection", err)
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
