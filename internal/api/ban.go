package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_controller "github.com/boxboxjason/jukebox/internal/controller"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

const (
	BANS_PREFIX = "/bans"
)

func SetBansRoutes(r chi.Router) {
	bans_subrouter := chi.NewRouter()

	// Authenticated routes
	bans_subrouter.Group(func(auth_router chi.Router) {
		auth_router.Use(middlewares.AdminAuthMiddleware)
		auth_router.Post("/", PostBan)
		auth_router.Get("/", GetBans)
		auth_router.Get(ID_PARAM_ENDPOINT, GetBan)
		auth_router.Patch(ID_PARAM_ENDPOINT, PatchBan)
		auth_router.Delete(ID_PARAM_ENDPOINT, DeleteBan)
		auth_router.Delete("/", DeleteBans)
	})

	r.Mount(BANS_PREFIX, bans_subrouter)
}

// ==================== CRUD operations ====================

// ==================== Create ====================
func PostBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the issuer
	issuer, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user not authenticated"))
		return
	}

	target_ids, err := httputils.RetrieveIntListValueParameter(r, constants.TARGET_ID_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	duration, err := httputils.RetrieveIntParameter(r, constants.DURATION_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ban_type, err := httputils.RetrieveStringParameter(r, constants.BAN_TYPE, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	targets, err := db_model.GetUsersByFilters(db, target_ids, nil, nil, nil, nil, 0)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	bans, err := db_controller.BanUsers(db, issuer, targets, duration, reason, ban_type)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, bans)
}

// ==================== Read ====================
func GetBan(w http.ResponseWriter, r *http.Request) {
	ban_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ban, err := db_controller.GetBanByID(nil, ban_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendJSONResponse(w, ban)
}

func GetBans(w http.ResponseWriter, r *http.Request) {
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAM, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	target_ids, err := httputils.RetrieveIntListValueParameter(r, constants.TARGET_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	issuer_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ISSUER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	types, err := httputils.RetrieveStringListValueParameter(r, constants.BAN_TYPE, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ends_after, err := httputils.RetrieveIntParameter(r, constants.ENDS_AFTER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	bans, err := db_controller.GetBans(nil, ids, target_ids, issuer_ids, types, reason, ends_after)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendJSONResponse(w, bans)
}

// ==================== Update ====================
func PatchBan(w http.ResponseWriter, r *http.Request) {
	ban_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	new_duration, err := httputils.RetrieveIntParameter(r, constants.DURATION_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	new_reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	ban, err := db_controller.GetBanByID(db, ban_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	err = db_controller.UpdateBan(db, ban, new_duration, new_reason)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendJSONResponse(w, ban)
}

// ==================== Delete ====================
func DeleteBan(w http.ResponseWriter, r *http.Request) {
	ban_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	ban, err := db_controller.GetBanByID(db, ban_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	err = db_controller.DeleteBan(db, ban)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "Ban deleted successfully")
}

func DeleteBans(w http.ResponseWriter, r *http.Request) {
	ban_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAM, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	issuer_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ISSUER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	types, err := httputils.RetrieveStringListValueParameter(r, constants.BAN_TYPE, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	target_ids, err := httputils.RetrieveIntListValueParameter(r, constants.TARGET_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	bans, err := db_controller.GetBans(db, ban_ids, target_ids, issuer_ids, types, "", 0)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	err = db_controller.DeleteBans(db, bans)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "Bans deleted successfully")
}
