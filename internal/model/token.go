package db_model

import (
	"fmt"

	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/gorm"
)

type AuthToken struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int        `gorm:"type:INTEGER;not null" json:"-"`
	User          *User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Hashed_Token  string     `gorm:"type:TEXT;unique;not null" json:"hashed_token"`
	Expiration    int64      `gorm:"type:INTEGER;not null" json:"expiration"`
	Type          string     `gorm:"type:TEXT;not null" json:"type"`
	LinkedTokenID *int       `gorm:"type:INTEGER;default:null" json:"linked_token"`
	LinkedToken   *AuthToken `gorm:"foreignKey:LinkedTokenID;constraint:OnDelete:SET NULL" json:"-"`
	CreatedAt     int        `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt    int        `gorm:"autoUpdateTime:milli" json:"modified_at"`
}

func (token *AuthToken) BeforeSave(tx *gorm.DB) (err error) {
	if len(token.Type) == 0 {
		token.Type = constants.ACCESS_TOKEN
	}
	validTypes := map[string]bool{constants.ACCESS_TOKEN: true, constants.REFRESH_TOKEN: true}
	if !validTypes[token.Type] {
		return fmt.Errorf("invalid type; must be either 'access' or 'refresh'")
	}
	return nil
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

func (user *User) CheckAuthTokenMatchesByType(db *gorm.DB, raw_token string, token_type string) (*AuthToken, error) {
	if len(user.Tokens) == 0 {
		var err error
		user, err = GetUserByID(db.Preload("Tokens"), user.ID)
		if err != nil {
			return &AuthToken{}, err
		}
	}

	for _, token := range user.Tokens {
		if token.Type == token_type && cryptutils.CompareHashAndString(token.Hashed_Token, raw_token) {
			return token, nil
		}
	}
	return &AuthToken{}, httputils.NewUnauthorizedError("Invalid token")
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
