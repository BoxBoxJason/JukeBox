package db_model

import (
	"time"

	"github.com/boxboxjason/jukebox/internal/constants"
	"gorm.io/gorm"
)

type Ban struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Target     *User  `gorm:"foreignKey:TargetID;constraint:OnDelete:CASCADE" json:"-"`
	TargetID   int    `json:"target_id"`
	Issuer     *User  `gorm:"foreignKey:IssuerID;constraint:OnDelete:CASCADE" json:"-"`
	IssuerID   int    `json:"issuer_id"`
	Reason     string `json:"reason"`
	EndsAt     int    `json:"ends_at"`
	Type       string `json:"type"`
	CreatedAt  int    `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt int    `gorm:"autoUpdateTime:milli" json:"modified_at"`
}

// ================ CRUD Operations ================
// ================ Create ================
// CreateBan creates a new ban in the database
func (ban *Ban) CreateBan(db *gorm.DB) error {
	return db.Create(ban).Error
}

// CreateBans creates multiple bans in the database
func CreateBans(db *gorm.DB, bans []*Ban) error {
	return db.Create(&bans).Error
}

// ================ Read ================
// GetBanByID retrieves a ban from the database by ID
func GetBanByID(db *gorm.DB, id int) (*Ban, error) {
	ban := &Ban{}
	err := db.First(ban, id).Error
	return ban, err
}

// GetBansByIssuerID retrieves all bans issued by a user from the database
func GetBansByIssuerID(db *gorm.DB, issuer_id int) ([]*Ban, error) {
	bans := []*Ban{}
	err := db.Where("issuer_id = ?", issuer_id).Find(&bans).Error
	return bans, err
}

// GetAllBans retrieves all bans from the database
func GetAllBans(db *gorm.DB) ([]*Ban, error) {
	bans := []*Ban{}
	err := db.Find(&bans).Error
	return bans, err
}

// GetActiveBans retrieves all active bans from the database
func GetActiveBans(db *gorm.DB) ([]*Ban, error) {
	current_time := int(time.Now().Unix())
	bans := []*Ban{}
	err := db.Where("ends_at > ?", current_time).Find(&bans).Error
	return bans, err
}

// GetBans retrieves all bans targeting a user from the database
func (user *User) GetBans(db *gorm.DB) ([]*Ban, error) {
	bans := []*Ban{}
	err := db.Where("target_id = ? AND type = ?", user.ID, constants.BAN_TYPE).Find(&bans).Error
	return bans, err
}

// GetMutes retrieves all mutes targeting a user from the database
func (user *User) GetMutes(db *gorm.DB) ([]*Ban, error) {
	bans := []*Ban{}
	err := db.Where("target_id = ? AND type = ?", user.ID, constants.MUTE_TYPE).Find(&bans).Error
	return bans, err
}

// GetActiveBans retrieves all active bans targeting a user from the database
// The current bans are sorted by the end date
func (user *User) GetActiveBans(db *gorm.DB) ([]*Ban, error) {
	current_time := int(time.Now().Unix())
	bans := []*Ban{}
	var err error
	if len(user.Bans) > 0 {
		for _, ban := range user.Bans {
			if ban.EndsAt > current_time {
				bans = append(bans, ban)
			}
		}
	} else {
		err = db.Where("target_id = ? AND ends_at > ? AND type = ?", user.ID, current_time, constants.BAN_TYPE).Order("ends_at desc").Find(&bans).Error
	}
	return bans, err
}

// GetActiveMutes retrieves all active mutes targeting a user from the database
// The current mutes are sorted by the end date
func (user *User) GetActiveMutes(db *gorm.DB) ([]*Ban, error) {
	current_time := int(time.Now().Unix())
	bans := []*Ban{}
	err := db.Where("target_id = ? AND ends_at > ? AND type = ?", user.ID, current_time, constants.MUTE_TYPE).Order("ends_at desc").Find(&bans).Error
	return bans, err
}

// GetIssuedBans retrieves all bans issued by a user from the database
func (user *User) GetIssuedBans(db *gorm.DB) ([]*Ban, error) {
	bans := []*Ban{}
	err := db.Where("issuer_id = ?", user.ID).Find(&bans).Error
	return bans, err
}

// GetActiveIssuedBans retrieves all active bans issued by a user from the database
func (user *User) GetActiveIssuedBans(db *gorm.DB) ([]*Ban, error) {
	current_time := int(time.Now().Unix())
	bans := []*Ban{}
	err := db.Where("issuer_id = ? AND ends_at > ?", user.ID, current_time).Find(&bans).Error
	return bans, err
}

func GetBansByFilters(db *gorm.DB, ids []int, target_ids []int, issuer_ids []int, types []string, reason string, ends_after int) ([]*Ban, error) {
	query := db.Where("ends_at > ?", ends_after)
	if len(target_ids) > 0 {
		query = query.Where("target_id IN ?", target_ids)
	}
	if len(issuer_ids) > 0 {
		query = query.Where("issuer_id IN ?", issuer_ids)
	}
	if len(types) > 0 {
		query = query.Where("type IN ?", types)
	}

	if reason != "" {
		query = query.Where("reason LIKE ?", "%"+reason+"%")
	}

	query.Or("id IN ?", ids)

	bans := []*Ban{}
	err := query.Find(&bans).Error
	return bans, err
}

// ================ Update ================
// UpdateBan updates a ban in the database
func (ban *Ban) UpdateBan(db *gorm.DB) error {
	return db.Save(ban).Error
}

// ================ Delete ================
// DeleteBan deletes a ban from the database
func (ban *Ban) DeleteBan(db *gorm.DB) error {
	return db.Delete(ban).Error
}

// DeleteBans deletes multiple bans from the database
func DeleteBans(db *gorm.DB, bans []*Ban) error {
	return db.Delete(&bans).Error
}
