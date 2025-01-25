package db_controller

import (
	"time"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"gorm.io/gorm"
)

// ================= CRUD Operations =================

// ================= Create =================
func BanUsers(db *gorm.DB, query_params *db_model.BansPostRequestParams) ([]*db_model.Ban, error) {
	bans := make([]*db_model.Ban, len(query_params.Target))
	for i, target := range query_params.Target {
		bans[i] = &db_model.Ban{
			IssuerID: query_params.Issuer.ID,
			TargetID: target.ID,
			Type:     query_params.Type,
			EndsAt:   time.Now().Add(time.Duration(query_params.Duration) * time.Second),
			Reason:   query_params.Reason,
		}
	}
	err := db_model.CreateBans(db, bans)
	return bans, err
}

func GetBans(db *gorm.DB, query_params *db_model.BansGetRequestParams) ([]*db_model.Ban, error) {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}
	return db_model.GetBansByFilters(db, query_params)
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
func UpdateBan(db *gorm.DB, query_params *db_model.BansPatchRequestParams) (*db_model.Ban, error) {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	ban, err := db_model.GetBanByID(db, query_params.ID)
	if err != nil {
		return nil, err
	}

	ban.EndsAt = time.Now().Add(time.Duration(query_params.Duration) * time.Second)
	ban.Reason = query_params.Reason

	return ban, ban.UpdateBan(db)
}

// ================= Delete =================
func DeleteBan(db *gorm.DB, ban_id int) error {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}
	return db_model.DeleteBanById(db, ban_id)
}

func DeleteBans(db *gorm.DB, query_params *db_model.BansDeleteRequestParams) error {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	bans, err := db_model.GetBansByFilters(db, &db_model.BansGetRequestParams{
		ID:       query_params.ID,
		TargetID: query_params.TargetID,
		IssuerID: query_params.IssuerID,
		Type:     query_params.Type,
		Order:    query_params.Order,
		Limit:    query_params.Limit,
		Page:     query_params.Page,
		Offset:   query_params.Offset,
	})

	if err != nil {
		return err
	}

	return db_model.DeleteBans(db, bans)
}
