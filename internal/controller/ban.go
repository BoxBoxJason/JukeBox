package db_controller

import (
	"time"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"gorm.io/gorm"
)

// ================= CRUD Operations =================

// ================= Create =================
func BanUser(db *gorm.DB, issuer *db_model.User, target *db_model.User, duration int, reason string, ban_type string) (*db_model.Ban, error) {
	ban := &db_model.Ban{
		IssuerID: issuer.ID,
		TargetID: target.ID,
		Type:     ban_type,
		EndsAt:   calculateBanEnd(duration),
		Reason:   reason,
	}
	err := ban.CreateBan(db)
	return ban, err
}

func BanUsers(db *gorm.DB, issuer *db_model.User, targets []*db_model.User, duration int, reason string, ban_type string) ([]*db_model.Ban, error) {
	bans := make([]*db_model.Ban, len(targets))
	for i, target := range targets {
		bans[i] = &db_model.Ban{
			IssuerID: issuer.ID,
			TargetID: target.ID,
			Type:     ban_type,
			EndsAt:   calculateBanEnd(duration),
			Reason:   reason,
		}
	}
	err := db_model.CreateBans(db, bans)
	return bans, err
}

// ================= Read =================
func calculateBanEnd(duration_hour int) int {
	return int(time.Now().Unix()) + duration_hour*3600
}

func GetBans(db *gorm.DB, ids []int, target_ids []int, issuer_ids []int, types []string, reason string, ends_after int) ([]*db_model.Ban, error) {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}
	return db_model.GetBansByFilters(db, ids, target_ids, issuer_ids, types, reason, ends_after)
}

func GetBanByID(db *gorm.DB, id int) (*db_model.Ban, error) {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}
	return db_model.GetBanByID(db, id)
}

// ================= Update =================
func UpdateBan(db *gorm.DB, ban *db_model.Ban, new_duration int, new_reason string) error {
	ban.EndsAt = calculateBanEnd(new_duration)
	ban.Reason = new_reason

	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	return ban.UpdateBan(db)
}

// ================= Delete =================
func DeleteBan(db *gorm.DB, ban *db_model.Ban) error {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}
	return ban.DeleteBan(db)
}

func DeleteBans(db *gorm.DB, bans []*db_model.Ban) error {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}
	return db_model.DeleteBans(db, bans)
}
