package db_model

import "gorm.io/gorm"

type Message struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Sender     *User  `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"-"`
	SenderID   int    `gorm:"type:INTEGER;not null" json:"sender_id"`
	Content    string `gorm:"type:TEXT;not null" json:"content"`
	Flagged    bool   `gorm:"type:BOOLEAN;default:false" json:"flagged"`
	Removed    bool   `gorm:"type:BOOLEAN;default:false" json:"removed"`
	Censored   bool   `gorm:"type:BOOLEAN;default:false" json:"censored"`
	CreatedAt  int    `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt int    `gorm:"autoUpdateTime:milli" json:"modified_at"`
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
