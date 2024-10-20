package db_model

import (
	"os"

	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/pkg/customerrors"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	CreateTables()
}

// Initialize the database connection and create the tables
func CreateTables() {
	db, err := OpenConnection()
	if err != nil {
		logger.Fatal("Failed to open the database connection")
	}
	defer CloseConnection(db)
	if !db.Migrator().HasTable(&User{}) {
		logger.Info("Creating tables")
		err := db.AutoMigrate(User{}, AuthToken{})
		if err != nil {
			logger.Info("Tables created successfully")
		} else {
			logger.Fatal("Failed to create tables:", err)
		}
	}
}

// OpenConnection opens a connection to the SQLite database
func OpenConnection() (*gorm.DB, error) {
	if _, err := os.Stat(constants.DB_FILE); os.IsNotExist(err) {
		err = os.Mkdir(constants.DB_DIR, 0600)
		if err != nil {
			logger.Critical("Failed to create the database directory", err)
			return nil, customerrors.NewDatabaseError("Failed to create the database directory")
		}
	} else if err != nil {
		logger.Critical("Failed to check if the database file exists:", err)
		return nil, customerrors.NewDatabaseError("Failed to check if the database file exists")
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
