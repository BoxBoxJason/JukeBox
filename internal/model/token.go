package db_model

import (
	"gorm.io/gorm"
)

type AuthToken struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	User         User   `gorm:"foreignKey:UserID" json:"user"`
	UserID       int    `gorm:"type:INTEGER;not null" json:"-"`
	Hashed_Token string `gorm:"type:TEXT;unique;not null" json:"hashed_token"`
	Expiration   int64  `gorm:"type:INTEGER;not null" json:"expiration"`
	Type         string `gorm:"type:TEXT;default:access" json:"type"`
	LinkedToken  int    `gorm:"type:INTEGER;default:-1" json:"linked_token"`
	CreatedAt    int    `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt   int    `gorm:"autoUpdateTime:milli" json:"modified_at"`
}

// ================ CRUD Operations ================
// ================ Create ================
// CreateAuthToken creates a new auth token in the database
func (auth_token *AuthToken) CreateAuthToken(db *gorm.DB) error {
	return db.Create(auth_token).Error
}

// CreateAuthTokens creates multiple auth tokens in the database
func CreateAuthTokens(db *gorm.DB, auth_tokens []*AuthToken) error {
	return db.Create(auth_tokens).Error
}

// ================ Read ================
// GetAuthTokenByID retrieves an auth token from the database by ID
func GetAuthTokenByID(db *gorm.DB, id int) (*AuthToken, error) {
	auth_token := &AuthToken{}
	err := db.First(&auth_token, id).Error
	return auth_token, err
}

// GetUserTokens retrieves all auth tokens for a user from the database
func (user *User) GetUserTokensByType(db *gorm.DB, token_type string) ([]*AuthToken, error) {
	var tokens []*AuthToken
	err := db.Where("user_id = ? AND type = ?", user.ID, token_type).Find(&tokens).Error
	return tokens, err
}

func (auth_token *AuthToken) GetLinkedToken(db *gorm.DB) (*AuthToken, error) {
	linked_token := &AuthToken{}
	err := db.First(&linked_token, auth_token.LinkedToken).Error
	return linked_token, err
}

// ================ Update ================
// UpdateAuthToken updates an auth token in the database
func (auth_token *AuthToken) UpdateAuthToken(db *gorm.DB) error {
	return db.Save(auth_token).Error
}

// ================ Delete ================
// DeleteAuthToken deletes an auth token from the database
func (auth_token *AuthToken) DeleteAuthToken(db *gorm.DB) error {
	return db.Delete(auth_token).Error
}

// DeleteAuthTokens deletes multiple auth tokens from the database
func DeleteAuthTokens(db *gorm.DB, auth_tokens []*AuthToken) error {
	return db.Delete(auth_tokens).Error
}

// DeleteUserTokens deletes all auth tokens for a user from the database
func (user *User) DeleteUserTokens(db *gorm.DB) error {
	return db.Where("user_id = ?", user.ID).Delete(AuthToken{}).Error
}

// DeleteAllTokens deletes all auth tokens from the database
func DeleteAllTokens(db *gorm.DB) error {
	return db.Delete(AuthToken{}).Error
}
