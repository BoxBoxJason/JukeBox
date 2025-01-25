package db_model

import (
	"mime/multipart"
	"time"

	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"gorm.io/gorm"
)

type User struct {
	ID                 int          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username           string       `gorm:"type:TEXT;unique;not null" json:"username"`
	Hashed_Password    string       `gorm:"type:TEXT;not null" json:"-"`
	Email              string       `gorm:"type:TEXT;unique;not null" json:"-"`
	Avatar             string       `gorm:"type:TEXT;default:'default_avatar.png'" json:"avatar"`
	Admin              bool         `gorm:"type:BOOLEAN;not null;default:false" json:"admin"`
	Banned             bool         `gorm:"type:BOOLEAN;not null;default:false" json:"banned"`
	TotalContributions int          `gorm:"type:INTEGER;not null;default:0" json:"total_contributions"`
	MinutesListened    int          `gorm:"type:INTEGER;not null;default:0" json:"minutes_listened"`
	Subscriber_Tier    int          `gorm:"type:INTEGER;not null;default:0" json:"subscriber_tier"`
	Messages           []*Message   `gorm:"foreignKey:SenderID" json:"-"`
	Tokens             []*AuthToken `gorm:"foreignKey:UserID" json:"-"`
	Bans               []*Ban       `gorm:"foreignKey:TargetID" json:"-"`
	CreatedAt          time.Time    `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt         time.Time    `gorm:"autoUpdateTime:milli" json:"modified_at"`
}

// ==================== USER ====================

// UsersPostRequestParams is the struct for the request body of the POST users endpoint
type UsersPostRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// UsersGetRequestParams is the struct for the request body of the GET users endpoint
type UsersGetRequestParams struct {
	Order           string   `json:"order"`
	Limit           int      `json:"limit"`
	Page            int      `json:"page"`
	Offset          int      `json:"offset"`
	Username        []string `json:"username"`
	PartialUsername []string `json:"partial_username"`
	ID              []int    `json:"id"`
	SubscriberTier  int      `json:"subscriber_tier"`
	Admin           []bool   `json:"admin"`
	Email           []string `json:"email"`
}

// UsersDeleteRequestParams is the struct for the request body of the DELETE users endpoint
type UsersDeleteRequestParams struct {
	Order    string   `json:"order"`
	Limit    int      `json:"limit"`
	Page     int      `json:"page"`
	Offset   int      `json:"offset"`
	ID       []int    `json:"id"`
	Username []string `json:"username"`
	Email    []string `json:"email"`
	Reason   string   `json:"reason"`
}

// UsersPatchRequestParams is the struct for the request body of the PATCH users endpoint
type UsersPatchRequestParams struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	SubscriberTier int    `json:"subscriber_tier"`
	Admin          []bool `json:"admin"`
	Avatar         string `json:"avatar"`
}

// UsersRawPatchRequestParams is the struct for the request body of the PATCH users endpoint
type UsersRawPatchRequestParams struct {
	ID             int            `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Password       string         `json:"password"`
	SubscriberTier int            `json:"subscriber_tier"`
	Admin          []bool         `json:"admin"`
	Avatar         multipart.File `json:"avatar"`
}

// ================ CRUD Operations ================
// ================ Create ================

// CreateUser creates a new user in the database
func (user *User) CreateUser(db *gorm.DB) error {
	return db.Create(user).Error
}

// CreateUsers creates multiple users in the database
func CreateUsers(db *gorm.DB, users []*User) error {
	return db.Create(&users).Error
}

// ================ Read ================

// GetUserByID retrieves a user from the database by ID
func GetUserByID(db *gorm.DB, id int) (*User, error) {
	user := &User{}
	err := db.First(user, id).Error
	return user, err
}

// GetUserByUsername retrieves a user from the database by username
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	user := &User{}
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

// GetUserByEmail retrieves a user from the database by email
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := &User{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

// GetUserByUsernameOREmail retrieves a user from the database by username or email
func GetUserByUsernameOREmail(db *gorm.DB, username_or_email string) (*User, error) {
	user := &User{}
	err := db.Where("username = ? OR email = ?", username_or_email, username_or_email).First(user).Error
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

func GetUsersByFilters(db *gorm.DB, query_params *UsersGetRequestParams) ([]*User, error) {
	users := []*User{}
	query := db

	// Build the "OR" conditions for ids, usernames, and partial usernames
	orConditions := db // Start an empty query for OR conditions

	// IDs
	if len(query_params.ID) > 0 {
		orConditions = orConditions.Or("id IN ?", query_params.ID)
	}

	// Usernames
	if len(query_params.Username) > 0 {
		orConditions = orConditions.Or("username IN ?", query_params.Username)
	}

	// Emails
	if len(query_params.Email) > 0 {
		orConditions = orConditions.Or("email IN ?", query_params.Email)
	}

	// Partial usernames
	if len(query_params.PartialUsername) > 0 {
		for _, partial_username := range query_params.PartialUsername {
			orConditions = orConditions.Or("username LIKE ?", "%"+partial_username+"%")
		}
	}

	// Combine the OR conditions into the main query
	query = query.Where(orConditions)

	// Apply the remaining "AND" filters
	if len(query_params.Admin) == 1 {
		query = query.Where("admin = ?", query_params.Admin[0])
	}
	if query_params.SubscriberTier > 0 {
		query = query.Where("subscriber_tier >= ?", query_params.SubscriberTier)
	}

	// Apply order, limit, page, and offset
	query = AddQueryParamsToDB(query, query_params.Order, query_params.Limit, query_params.Page, query_params.Offset)

	err := query.Find(&users).Error
	return users, err
}

// IsSubscriber checks if a user is a subscriber
func (user *User) IsSubscriber() bool {
	return user.Subscriber_Tier > 0
}

func (user *User) CheckPasswordMatches(password string) bool {
	return cryptutils.CompareHashAndString(user.Hashed_Password, password)
}

// ================ Update ================

// UpdateUser updates a user in the database
func (user *User) UpdateUser(db *gorm.DB) error {
	return db.Save(user).Error
}

func (user *User) IncreaseContributionsCount(db *gorm.DB) error {
	user.TotalContributions++
	return db.Save(user).Error
}

// UpdateUsers updates multiple users in the database
func UpdateUsers(db *gorm.DB, users []*User) error {
	return db.Save(users).Error
}

// ================ Delete ================

// DeleteUser deletes a user from the database
func (user *User) DeleteUser(db *gorm.DB) error {
	return db.Delete(user).Error
}

// DeleteUsers deletes multiple users from the database
func DeleteUsers(db *gorm.DB, users []*User) error {
	return db.Delete(users).Error
}
