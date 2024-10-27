package db_controller

import (
	"sync"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"gorm.io/gorm"
)

// ================= CRUD Operations =================

// ================= Create =================

// CreateMessage creates a new message in the database
func CreateMessage(db *gorm.DB, message string, user *db_model.User) (*db_model.Message, error) {
	db_message := db_model.Message{
		Sender:  user,
		Content: message,
	}

	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return &db_message, err
		}
		defer db_model.CloseConnection(db)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Create the message and increase the user's contributions count
	var err error
	go func() {
		defer wg.Done()
		err = db_message.CreateMessage(db)
	}()
	go func() {
		defer wg.Done()
		user.IncreaseContributionsCount(db)
	}()

	return &db_message, err
}

// ================= Read =================
func GetMessages(db *gorm.DB) ([]*db_model.Message, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	return db_model.GetAllVisibleMessages(db)
}

func GetMessage(id int) (*db_model.Message, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)
	return db_model.GetMessageByID(db, id)
}

// ================= Update =================

// UpdateMessage updates a message in the database
func UpdateMessage(db *gorm.DB, id int, message string) (*db_model.Message, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	// Retrieve the message
	db_message, err := db_model.GetMessageByID(db, id)
	if err != nil {
		return nil, err
	}

	// Update the message
	db_message.Content = message
	err = db_message.UpdateMessage(db)
	return db_message, err
}

// ================= Delete =================

// DeleteMessage deletes a message from the database
func DeleteMessage(db *gorm.DB, id int) error {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	// Retrieve the message
	db_message, err := db_model.GetMessageByID(db, id)
	if err != nil {
		return err
	}

	// Delete the message
	return db_message.DeleteMessage(db)
}

func DeleteMessages(db *gorm.DB) error {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	// Retrieve all messages
	messages, err := db_model.GetAllMessages(db)
	if err != nil {
		return err
	}

	// Delete all messages
	return db_model.DeleteMessages(db, messages)
}

func UserHasPermissionToDeleteMessage(user *db_model.User, message *db_model.Message) bool {
	return message.Sender.ID == user.ID || user.Admin
}
