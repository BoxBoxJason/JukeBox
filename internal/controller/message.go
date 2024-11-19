package db_controller

import (
	"strings"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/gorm"
)

// ================= CRUD Operations =================

// ================= Create =================

// CreateMessage creates a new message in the database
func CreateMessage(db *gorm.DB, message string, user *db_model.User) (*db_model.Message, error) {
	db_message := db_model.Message{
		Sender:  user,
		Content: strings.TrimSpace(message),
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

	err := db_message.CreateMessage(db)
	if err != nil {
		return &db_message, err
	}

	err = user.IncreaseContributionsCount(db)

	return &db_message, err
}

// ================= Read =================
func GetMessages(db *gorm.DB, ids []int, sender_ids []int, flagged []bool, censored []bool, removed []bool, contains []string) ([]*db_model.Message, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	for _, id := range ids {
		if id < 0 {
			return nil, httputils.NewBadRequestError("id must be a positive integer")
		}
	}

	for _, sender_id := range sender_ids {
		if sender_id < 0 {
			return nil, httputils.NewBadRequestError("sender_id must be a positive integer")
		}
	}

	return db_model.GetMessages(db, ids, sender_ids, flagged, censored, removed, contains)
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
func UpdateMessage(db *gorm.DB, db_message *db_model.Message, message string, flagged []bool, censored []bool, removed []bool) (*db_model.Message, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	if len(message) > 0 {
		db_message.Content = message
	}

	if len(flagged) > 0 {
		db_message.Flagged = flagged[0]
	}

	if len(censored) > 0 {
		db_message.Censored = censored[0]
	}

	if len(removed) > 0 {
		db_message.Removed = removed[0]
	}

	// Update the message
	err := db_message.UpdateMessage(db)
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

func DeleteMessages(db *gorm.DB, ids []int, sender_ids []int, flagged []bool, censored []bool, removed []bool, contains []string) error {
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
	messages, err := db_model.GetMessages(db, ids, sender_ids, flagged, censored, removed, contains)
	if err != nil {
		return err
	}

	// Delete all messages
	return db_model.DeleteMessages(db, messages)
}
