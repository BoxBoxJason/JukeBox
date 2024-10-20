package db_model

import (
	"github.com/boxboxjason/jukebox/pkg/customerrors"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"gorm.io/gorm"
)

type User struct {
	ID                 int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username           string `gorm:"type:TEXT;unique;not null" json:"username"`
	Hashed_Password    string `gorm:"type:TEXT;not null" json:"hashed_password"`
	Email              string `gorm:"type:TEXT;unique;not null" json:"email"`
	Admin              bool   `gorm:"type:BOOLEAN;not null;default:false" json:"admin"`
	Banned             bool   `gorm:"type:BOOLEAN;not null;default:false" json:"banned"`
	TotalContributions int    `gorm:"type:INTEGER;not null;default:0" json:"total_contributions"`
	MinutesListened    int    `gorm:"type:INTEGER;not null;default:0" json:"minutes_listened"`
	Subscriber_Tier    int    `gorm:"type:INTEGER;not null;default:0" json:"subscriber_tier"`
	CreatedAt          int    `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt         int    `gorm:"autoUpdateTime:milli" json:"modified_at"`
}

// ================ CRUD Operations ================
// ================ Create ================

// CreateUser creates a new user in the database
func (user *User) CreateUser(db *gorm.DB) error {
	return db.Create(&user).Error
}

// CreateUsers creates multiple users in the database
func CreateUsers(db *gorm.DB, users []*User) error {
	return db.Create(&users).Error
}

// ================ Read ================

// GetUserByID retrieves a user from the database by ID
func GetUserByID(db *gorm.DB, id int) (*User, error) {
	user := &User{}
	err := db.First(&user, id).Error
	return user, err
}

// GetUserByUsername retrieves a user from the database by username
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	user := &User{}
	err := db.Where("username = ?", username).First(&user).Error
	return user, err
}

// GetUserByEmail retrieves a user from the database by email
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := &User{}
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

// GetUserByUsernameOREmail retrieves a user from the database by username or email
func GetUserByUsernameOREmail(db *gorm.DB, username_or_email string) (*User, error) {
	user := &User{}
	err := db.Where("username = ? OR email = ?", username_or_email, username_or_email).First(&user).Error
	return user, err
}

// GetUsersByUsername retrieves all users which username partially matches the given username
func GetUsersByUsername(db *gorm.DB, username string) ([]*User, error) {
	users := []*User{}
	err := db.Where("username LIKE ?", "%"+username+"%").Find(&users).Error
	return users, err
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]*User, error) {
	users := []*User{}
	err := db.Find(&users).Error
	return users, err
}

// IsAdmin checks if a user is an admin
func (user *User) IsAdmin() bool {
	return user.Admin
}

// IsBanned checks if a user is banned
func (user *User) IsBanned() bool {
	return user.Banned
}

// IsSubscriber checks if a user is a subscriber
func (user *User) IsSubscriber() bool {
	return user.Subscriber_Tier > 0
}

// GetMinutesListened retrieves the number of minutes a user has listened to music
func (user *User) GetMinutesListened() int {
	return user.MinutesListened
}

// GetSubscriberTier retrieves the subscriber tier of a user
func (user *User) GetSubscriberTier() int {
	return user.Subscriber_Tier
}

// GetTotalContributions retrieves the total number of contributions a user has made
func (user *User) GetTotalContributions() int {
	return user.TotalContributions
}

func (user *User) CheckPasswordMatches(password string) bool {
	return cryptutils.CompareHashAndString(user.Hashed_Password, password)
}

func (user *User) CheckAuthTokenMatchesByType(db *gorm.DB, raw_token string, token_type string) (AuthToken, error) {
	tokens, err := user.GetUserTokensByType(db, token_type)
	if err != nil {
		return AuthToken{}, err
	}
	for _, token := range tokens {
		if cryptutils.CompareHashAndString(token.Hashed_Token, raw_token) {
			return token, nil
		}
	}
	return AuthToken{}, customerrors.NewUnauthorizedError("Invalid token")
}

// ToJSON converts a user to a JSON object
func (user *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":                  user.ID,
		"username":            user.Username,
		"admin":               user.Admin,
		"banned":              user.Banned,
		"total_contributions": user.TotalContributions,
		"minutes_listened":    user.MinutesListened,
		"subscriber_tier":     user.Subscriber_Tier,
	}
}

// ================ Update ================

// UpdateUser updates a user in the database
func (user *User) UpdateUser(db *gorm.DB) error {
	return db.Save(&user).Error
}

// UpdateUsers updates multiple users in the database
func UpdateUsers(db *gorm.DB, users []*User) error {
	return db.Save(&users).Error
}

// ================ Delete ================

// DeleteUser deletes a user from the database
func (user *User) DeleteUser(db *gorm.DB) error {
	return db.Delete(&user).Error
}

// DeleteUsers deletes multiple users from the database
func DeleteUsers(db *gorm.DB, users []*User) error {
	return db.Delete(&users).Error
}
