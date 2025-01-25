package db_model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Sender     *User     `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender"`
	SenderID   int       `gorm:"type:INTEGER;not null" json:"-"`
	Content    string    `gorm:"type:TEXT;not null" json:"content"`
	Flagged    bool      `gorm:"type:BOOLEAN;default:false" json:"flagged"`
	Removed    bool      `gorm:"type:BOOLEAN;default:false" json:"removed"`
	Censored   bool      `gorm:"type:BOOLEAN;default:false" json:"censored"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt time.Time `gorm:"autoUpdateTime:milli" json:"modified_at"`
}

// ==================== Requests parameters ====================

// MessagesPostRequestParams is the struct for the request body of the POST messages endpoint
type MessagesPostRequestParams struct {
	Message string `json:"message"`
	Sender  *User  `json:"sender"`
}

// MessagesGetRequestParams is the struct for the request body of the GET messages endpoint
type MessagesGetRequestParams struct {
	Order    string   `json:"order"`
	Limit    int      `json:"limit"`
	Page     int      `json:"page"`
	Offset   int      `json:"offset"`
	ID       []int    `json:"id"`
	SenderID []int    `json:"sender_id"`
	Flagged  []bool   `json:"flagged"`
	Censored []bool   `json:"censored"`
	Removed  []bool   `json:"removed"`
	Contains []string `json:"contains"`
}

// MessagesPatchRequestParams is the struct for the request body of the PATCH messages endpoint
type MessagesPatchRequestParams struct {
	Order    string `json:"order"`
	Limit    int    `json:"limit"`
	Page     int    `json:"page"`
	Offset   int    `json:"offset"`
	ID       int    `json:"id"`
	Message  string `json:"message"`
	Removed  []bool `json:"removed"`
	Flagged  []bool `json:"flagged"`
	Censored []bool `json:"censored"`
}

// MessagesDeleteRequestParams is the struct for the request body of the DELETE messages endpoint
type MessagesDeleteRequestParams struct {
	Order    string   `json:"order"`
	Limit    int      `json:"limit"`
	Page     int      `json:"page"`
	Offset   int      `json:"offset"`
	ID       []int    `json:"id"`
	SenderID []int    `json:"sender_id"`
	Flagged  []bool   `json:"flagged"`
	Censored []bool   `json:"censored"`
	Removed  []bool   `json:"removed"`
	Contains []string `json:"contains"`
}

// ================ CRUD Operations ================

// ================ Create ================

// CreateMessage creates a new message in the database
func (message *Message) CreateMessage(db *gorm.DB) error {
	return db.Create(message).Error
}

// CreateMessages creates multiple messages in the database
func CreateMessages(db *gorm.DB, messages []*Message) error {
	return db.Create(messages).Error
}

// ================ Read ================

// GetMessageByID retrieves a message from the database by ID
func GetMessageByID(db *gorm.DB, id int) (*Message, error) {
	message := &Message{}
	err := db.Preload("Sender").First(&message, id).Error
	return message, err
}

// GetMessages retrieves all messages sent by a user from the database
// Does not retrieve the removed and censored messages
func (user *User) GetMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("sender_id = ? AND removed = ? AND censored = ?", user.ID, false, false).Find(&messages).Error
	return messages, err
}

// GetAllMessages retrieves all messages sent by a user from the database, including the removed and censored messages
func (user *User) GetAllMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("sender_id = ?", user.ID).Find(&messages).Error
	return messages, err
}

// GetMessagesByContent retrieves all messages containing a specific content from the database
// Does not retrieve the removed and censored messages
func GetMessagesByContent(db *gorm.DB, content string) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("content LIKE ? AND removed = ? AND censored = ?", "%"+content+"%", false, false).Find(&messages).Error
	return messages, err
}

// GetAllMessagesByContent retrieves all messages containing a specific content from the database, including the removed and censored messages
func GetAllMessagesByContent(db *gorm.DB, content string) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("content LIKE ?", "%"+content+"%").Find(&messages).Error
	return messages, err
}

// GetAllVisibleMessages retrieves all messages from the database, sorted by creation date
// Does not retrieve the removed and censored messages
func GetAllVisibleMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Order("created_at desc").Where("removed = ? AND censored = ?", false, false).Find(&messages).Error
	return messages, err
}

// GetAllMessages retrieves all messages from the database, sorted by creation date
// Includes the removed and censored messages
func GetAllMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Order("created_at desc").Find(&messages).Error
	return messages, err
}

// GetFlaggedMessages retrieves all flagged messages from the database
func GetFlaggedMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("flagged = ?", true).Find(&messages).Error
	return messages, err
}

// GetUserFlaggedMessages retrieves all flagged messages sent by a user from the database
func (user *User) GetFlaggedMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("sender_id = ? AND flagged = ?", user.ID, true).Find(&messages).Error
	return messages, err
}

// GetRemovedMessages retrieves all removed messages from the database
func GetRemovedMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("removed = ?", true).Find(&messages).Error
	return messages, err
}

// GetUserRemovedMessages retrieves all removed messages sent by a user from the database
func (user *User) GetRemovedMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("sender_id = ? AND removed = ?", user.ID, true).Find(&messages).Error
	return messages, err
}

// GetCensoredMessages retrieves all censored messages from the database
func GetCensoredMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("censored = ?", true).Find(&messages).Error
	return messages, err
}

// GetUserCensoredMessages retrieves all censored messages sent by a user from the database
func (user *User) GetCensoredMessages(db *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := db.Preload("Sender").Where("sender_id = ? AND censored = ?", user.ID, true).Find(&messages).Error
	return messages, err
}

func GetMessages(db *gorm.DB, query_params *MessagesGetRequestParams) ([]*Message, error) {
	query := db
	for _, s := range query_params.Contains {
		query = query.Or("content LIKE ?", "%"+s+"%")
	}

	if len(query_params.SenderID) > 0 {
		query = query.Where("sender_id IN ?", query_params.SenderID)
	}
	if len(query_params.Flagged) > 0 {
		query = query.Where("flagged = ?", query_params.Flagged[0])
	}
	if len(query_params.Censored) > 0 {
		query = query.Where("censored = ?", query_params.Censored[0])
	}
	if len(query_params.Removed) > 0 {
		query = query.Where("removed = ?", query_params.Removed[0])
	}

	if len(query_params.ID) > 0 {
		query = query.Or("id IN ?", query_params.ID)
	}

	// Apply order, limit, page, and offset
	query = AddQueryParamsToDB(query, query_params.Order, query_params.Limit, query_params.Page, query_params.Offset)

	var messages []*Message
	err := query.Find(&messages).Error
	return messages, err
}

// ================ Update ================

// UpdateMessage updates a message in the database
func (message *Message) UpdateMessage(db *gorm.DB) error {
	return db.Save(message).Error
}

// ================ Delete ================

// DeleteMessage deletes a message from the database
func (message *Message) DeleteMessage(db *gorm.DB) error {
	return db.Delete(message).Error
}

func DeleteMessage(db *gorm.DB, id int) error {
	message, err := GetMessageByID(db, id)
	if err != nil {
		return err
	}
	return message.DeleteMessage(db)
}

// DeleteMessages deletes multiple messages from the database
func DeleteMessages(db *gorm.DB, messages []*Message) error {
	return db.Delete(messages).Error
}

// DeleteAllMessages deletes all messages from the database
func DeleteAllMessages(db *gorm.DB) error {
	return db.Delete(Message{}).Error
}

// DeleteUserMessages deletes all messages sent by a user from the database
func (user *User) DeleteUserMessages(db *gorm.DB) error {
	return db.Where("sender_id = ?", user.ID).Delete(Message{}).Error
}

// DeleteMessagesByContent deletes all messages containing a specific content from the database
func DeleteMessagesByContent(db *gorm.DB, content string) error {
	return db.Where("content LIKE ?", "%"+content+"%").Delete(Message{}).Error
}
