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
func CreateMessage(db *gorm.DB, query_params *db_model.MessagesPostRequestParams) (*db_model.Message, error) {
	db_message := db_model.Message{
		Sender:  query_params.Sender,
		Content: strings.TrimSpace(query_params.Message),
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

	err = query_params.Sender.IncreaseContributionsCount(db)

	return &db_message, err
}

// ================= Read =================
func GetMessages(db *gorm.DB, query_params *db_model.MessagesGetRequestParams) ([]*db_model.Message, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	for _, id := range query_params.ID {
		if id < 0 {
			return nil, httputils.NewBadRequestError("id must be a positive integer")
		}
	}

	for _, sender_id := range query_params.SenderID {
		if sender_id < 0 {
			return nil, httputils.NewBadRequestError("sender_id must be a positive integer")
		}
	}

	return db_model.GetMessages(db.Preload("Sender"), query_params)
}

func GetMessage(id int) (*db_model.Message, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)
	return db_model.GetMessageByID(db.Preload("Sender"), id)
}

// ================= Update =================

// UpdateMessage updates a message in the database
func UpdateMessage(db *gorm.DB, query_params *db_model.MessagesPatchRequestParams) (*db_model.Message, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	db_message, err := db_model.GetMessageByID(db.Preload("Sender"), query_params.ID)
	if err != nil {
		return nil, err
	}

	if len(query_params.Message) > 0 {
		db_message.Content = query_params.Message
	}

	if len(query_params.Flagged) > 0 {
		db_message.Flagged = query_params.Flagged[0]
	}

	if len(query_params.Censored) > 0 {
		db_message.Censored = query_params.Censored[0]
	}

	if len(query_params.Removed) > 0 {
		db_message.Removed = query_params.Removed[0]
	}

	// Update the message
	err = db_message.UpdateMessage(db)
	return db_message, err
}

// UpdateExistingMessage updates an existing message in the database
func UpdateExistingMessage(db *gorm.DB, message *db_model.Message, query_params *db_model.MessagesPatchRequestParams) error {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	if len(query_params.Message) > 0 {
		message.Content = query_params.Message
	}

	if len(query_params.Flagged) > 0 {
		message.Flagged = query_params.Flagged[0]
	}

	if len(query_params.Censored) > 0 {
		message.Censored = query_params.Censored[0]
	}

	if len(query_params.Removed) > 0 {
		message.Removed = query_params.Removed[0]
	}

	// Update the message
	return message.UpdateMessage(db)
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

	return db_model.DeleteMessage(db, id)
}

func DeleteMessages(db *gorm.DB, query_params *db_model.MessagesDeleteRequestParams) error {
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
	messages, err := db_model.GetMessages(db, (*db_model.MessagesGetRequestParams)(query_params))
	if err != nil {
		return err
	}

	// Delete all messages
	return db_model.DeleteMessages(db, messages)
}
