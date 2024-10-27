package api

import (
	"net/http"
	"strconv"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_controller "github.com/boxboxjason/jukebox/internal/controller"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

const (
	USERS_PREFIX = "/users"
)

func SetUsersRoutes(r chi.Router) {
	users_subrouter := chi.NewRouter()

	// Unauthenticated routes
	users_subrouter.Post("/", CreateUser)

	// Authenticated routes
	users_subrouter.Group(func(auth_router chi.Router) {
		auth_router.Use(middlewares.AuthMiddleware)
		auth_router.Get("/", GetUsers)
		auth_router.Get(ID_PARAM_ENDPOINT, GetUser)
		auth_router.Put(ID_PARAM_ENDPOINT, UpdateUser)
		auth_router.Delete(ID_PARAM_ENDPOINT, DeleteUser)
		auth_router.Delete("/", DeleteUsers)
		auth_router.Patch(ID_PARAM_ENDPOINT, UpdateUser)
		auth_router.Get(ID_PARAM_ENDPOINT+"/ban", GetUserBans)
	})

	// Admin routes
	users_subrouter.Group(func(admin_router chi.Router) {
		admin_router.Use(middlewares.AdminAuthMiddleware)
		admin_router.Post(ID_PARAM_ENDPOINT+"/ban", CreateUserBan)
	})

	r.Mount(USERS_PREFIX, users_subrouter)
}

// ==================== CRUD operations ====================

// ==================== Create ====================
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Check that the user is not already authenticated
	_, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if ok {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user already authenticated"))
		return
	}

	// Retrieve the user input
	username, err := httputils.RetrievePostFormStringParameter(r, constants.USERNAME_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	email, err := httputils.RetrievePostFormStringParameter(r, constants.EMAIL_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	password, err := httputils.RetrievePostFormStringParameter(r, constants.PASSWORD_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Create the user
	db_user, err := db_controller.CreateUser(nil, username, email, password)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the user back to the client
	httputils.SendJSONResponse(w, db_user)
}

func CreateUserBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	issuer, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	reason, err := httputils.RetrievePostFormStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	duration, err := httputils.RetrievePostFormIntParameter(r, constants.DURATION_PARAMETER, false)
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
	defer db_model.CloseConnection(db)

	user_to_ban, err := db_controller.GetUser(db, user_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ban, err := db_controller.BanUser(db, issuer, user_to_ban, duration, reason, ban_type)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, ban)
}

// ==================== Read ====================
func GetUsers(w http.ResponseWriter, r *http.Request) {
	usernames, err := httputils.RetrieveStringListValueParameter(r, constants.USERNAME_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	partial_username, err := httputils.RetrieveStringListValueParameter(r, constants.PARTIAL_USERNAME_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAM, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	minimum_subscriber_tier, err := httputils.RetrieveIntParameter(r, constants.SUBSCRIBER_TIER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	admin, err := httputils.RetrieveBoolParameter(r, constants.ADMIN_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	users, err := db_controller.GetUsers(nil, ids, usernames, partial_username, nil, admin, minimum_subscriber_tier)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	user, err := db_controller.GetUser(nil, user_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, user)
}

func GetUserBans(w http.ResponseWriter, r *http.Request) {
	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ends_after, err := httputils.RetrieveIntParameter(r, constants.ENDS_AFTER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ban_types, err := httputils.RetrieveStringListValueParameter(r, constants.BAN_TYPE, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	issuer_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ISSUER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	bans, err := db_controller.GetBans(nil, nil, []int{user_id}, issuer_ids, ban_types, "", ends_after)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, bans)
}

// ==================== Update ====================
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	if user.ID != user_id && !user.Admin {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user does not have permission to update user"))
		return
	}

	// User reserved fields
	email, err := httputils.RetrievePostFormStringParameter(r, constants.EMAIL_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	} else if len(email) > 0 && user.ID != user_id {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user does not have permission to update email for user: "+strconv.Itoa(user_id)))
		return
	}
	password, err := httputils.RetrievePostFormStringParameter(r, constants.PASSWORD_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	} else if len(password) > 0 && user.ID != user_id {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user does not have permission to update password for user: "+strconv.Itoa(user_id)))
		return
	}

	// Shared fields
	username, err := httputils.RetrievePostFormStringParameter(r, constants.USERNAME_PARAMETER, true)
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

	user_to_update, err := db_controller.GetUser(db, user_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	avatar_file, _, err := httputils.RetrieveImageFile(r, constants.AVATAR_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	user_to_update, err = db_controller.UpdateUser(db, user_to_update, username, email, password, avatar_file)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, user_to_update)
}

// ==================== Delete ====================
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	user_to_delete_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAM)
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

	user_to_delete, err := db_controller.GetUser(db, user_to_delete_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	if !db_controller.UserHasPermissionToDeleteUser(user, user_to_delete) {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user does not have permission to delete user"))
		return
	}

	err = db_controller.DeleteUser(db, user_to_delete)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "user deleted")
}

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	requester, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	usernames, err := httputils.RetrieveStringListValueParameter(r, constants.USERNAME_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAM, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	emails, err := httputils.RetrieveStringListValueParameter(r, constants.EMAIL_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
	reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	err = db_controller.DeleteUsers(nil, requester, ids, usernames, emails, reason)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "users deleted")
}
