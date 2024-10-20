package db_tables

import (
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
)

func init() {
	CreateTables()
}

// Initialize the database connection and create the tables
func CreateTables() {
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Fatal("Failed to connect database:", err)
	} else {
		if !db.Migrator().HasTable(&User{}) {
			logger.Info("Creating tables")
			err = db.AutoMigrate(User{}, AuthToken{})
			if err != nil {
				logger.Info("Tables created successfully")
			} else {
				logger.Fatal("Failed to create tables:", err)
			}
		}
	}
}

type User struct {
	ID                 int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username           string `gorm:"type:TEXT;unique;not null" json:"username"`
	Hashed_Password    string `gorm:"type:TEXT;not null" json:"hashed_password"`
	Email              string `gorm:"type:TEXT;unique;not null" json:"email"`
	Admin              bool   `gorm:"type:BOOLEAN;not null;default:false" json:"admin"`
	Banned             bool   `gorm:"type:BOOLEAN;not null;default:false" json:"banned"`
	TotalContributions int    `gorm:"type:INTEGER;not null;default:0" json:"total_contributions"`
	Subscriber_Tier    int    `gorm:"type:INTEGER;not null;default:0" json:"subscriber_tier"`
}

type AuthToken struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	User         User   `gorm:"foreignKey:UserID" json:"user"`
	Hashed_Token string `gorm:"type:TEXT;not null" json:"hashed_token"`
	Expiration   int64  `gorm:"type:INTEGER;not null" json:"expiration"`
}
