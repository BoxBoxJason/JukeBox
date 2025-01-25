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

func SetupBansRoutes(r chi.Router) {
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

// PostBan creates a new ban in the database
func PostBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the issuer from the context (middleware generated)
	issuer, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user not authenticated"))
		return
	}

	// Retrieve the targets ids from query parameters
	target_ids, err := httputils.RetrieveIntListValueParameter(r, constants.TARGET_ID_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the duration from query parameters
	duration, err := httputils.RetrieveIntParameter(r, constants.DURATION_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the reason from query parameters
	reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ban type from query parameters
	ban_type, err := httputils.RetrieveStringParameter(r, constants.BAN_TYPE, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	// Retrieve the targets from the database
	targets, err := db_model.GetUsersByFilters(db, &db_model.UsersGetRequestParams{ID: target_ids})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Ban the targets
	bans, err := db_controller.BanUsers(db, &db_model.BansPostRequestParams{
		Target:   targets,
		Issuer:   issuer,
		Reason:   reason,
		Duration: duration,
		Type:     ban_type,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, bans)
}

// ==================== Read ====================

// GetBan retrieves a ban from the database by ID
func GetBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ban id from the query parameters
	ban_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ban from the database
	ban, err := db_controller.GetBanByID(nil, ban_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendJSONResponse(w, ban)
}

// GetBans retrieves all bans from the database by the given query parameters
func GetBans(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ban ids from the query parameters
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the target ids from the query parameters
	target_ids, err := httputils.RetrieveIntListValueParameter(r, constants.TARGET_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the issuer ids from the query parameters
	issuer_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ISSUER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ban types from the query parameters
	types, err := httputils.RetrieveStringListValueParameter(r, constants.BAN_TYPE, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the reason from the query parameters
	reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ends_after from the query parameters
	ends_after, err := httputils.RetrieveTimeStampParameter(r, constants.ENDS_AFTER_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the base parameters for the request
	order, limit, page, offset := retrieveBaseParams(r)

	// Retrieve the bans from the database
	bans, err := db_controller.GetBans(nil, &db_model.BansGetRequestParams{
		ID:        ids,
		TargetID:  target_ids,
		IssuerID:  issuer_ids,
		Type:      types,
		Reason:    reason,
		EndsAfter: ends_after,
		Order:     order,
		Limit:     limit,
		Page:      page,
		Offset:    offset,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendJSONResponse(w, bans)
}

// ==================== Update ====================

// PatchBan updates a ban in the database
func PatchBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ban id from the query parameters
	ban_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the new duration from the query parameters
	new_duration, err := httputils.RetrieveIntParameter(r, constants.DURATION_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the new reason from the query parameters
	new_reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	// Update the ban
	ban, err := db_controller.UpdateBan(db, &db_model.BansPatchRequestParams{
		ID:       ban_id,
		Duration: new_duration,
		Reason:   new_reason,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendJSONResponse(w, ban)
}

// ==================== Delete ====================

// DeleteBan deletes a ban from the database
func DeleteBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ban id from the query parameters
	ban_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	// Delete the ban
	err = db_controller.DeleteBan(db, ban_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "Ban deleted successfully")
}

// DeleteBans deletes multiple bans from the database
func DeleteBans(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ban ids from the query parameters
	ban_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the issuer ids from the query parameters
	issuer_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ISSUER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ban types from the query parameters
	types, err := httputils.RetrieveStringListValueParameter(r, constants.BAN_TYPE, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the target ids from the query parameters
	target_ids, err := httputils.RetrieveIntListValueParameter(r, constants.TARGET_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the base parameters for the request
	order, limit, page, offset := retrieveBaseParams(r)

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	err = db_controller.DeleteBans(db, &db_model.BansDeleteRequestParams{
		ID:       ban_ids,
		IssuerID: issuer_ids,
		Type:     types,
		TargetID: target_ids,
		Order:    order,
		Limit:    limit,
		Page:     page,
		Offset:   offset,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "Bans deleted successfully")
}
