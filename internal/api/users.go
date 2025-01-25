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
		auth_router.Get(ID_PARAM_ENDPOINT+BANS_PREFIX, GetUserBans)
		auth_router.Get(ID_PARAM_ENDPOINT+MESSAGES_PREFIX, GetUserMessages)
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

// CreateUser creates a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the request body
	username, err := httputils.RetrievePostFormStringParameter(r, constants.USERNAME_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the email from the request body
	email, err := httputils.RetrievePostFormStringParameter(r, constants.EMAIL_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the password from the request body
	password, err := httputils.RetrievePostFormStringParameter(r, constants.PASSWORD_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Create the user
	db_user, err := db_controller.CreateUser(nil, &db_model.UsersPostRequestParams{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the user back to the client
	httputils.SendJSONResponse(w, db_user)
}

// CreateUserBan creates a new ban in the database for a user
func CreateUserBan(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	issuer, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	// Retrieve the user id from the request parameters
	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the reason from the request body
	reason, err := httputils.RetrievePostFormStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the duration from the request body
	duration, err := httputils.RetrievePostFormIntParameter(r, constants.DURATION_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ban type from the request body
	ban_type, err := httputils.RetrievePostFormStringParameter(r, constants.BAN_TYPE, false)
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

	user_to_ban, err := db_controller.GetUser(db, user_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	ban, err := db_controller.BanUsers(db, &db_model.BansPostRequestParams{
		Issuer:   issuer,
		Target:   []*db_model.User{user_to_ban},
		Type:     ban_type,
		Duration: duration,
		Reason:   reason,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, ban)
}

// ==================== Read ====================

// GetUsers retrieves users from the database based on the request parameters
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Retrieve user parameters from the request
	usernames, err := httputils.RetrieveStringListValueParameter(r, constants.USERNAME_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the partial username from the request
	partial_username, err := httputils.RetrieveStringListValueParameter(r, constants.PARTIAL_USERNAME_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the user ids from the request
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the minimum subscriber tier from the request
	minimum_subscriber_tier, err := httputils.RetrieveIntParameter(r, constants.SUBSCRIBER_TIER_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the admin parameter from the request
	admin, err := httputils.RetrieveBoolParameter(r, constants.ADMIN_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the base parameters for the request
	order, limit, page, offset := retrieveBaseParams(r)

	users, err := db_controller.GetUsers(nil, &db_model.UsersGetRequestParams{
		Username:        usernames,
		PartialUsername: partial_username,
		ID:              ids,
		SubscriberTier:  minimum_subscriber_tier,
		Admin:           admin,
		Order:           order,
		Limit:           limit,
		Page:            page,
		Offset:          offset,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, users)
}

// GetUser retrieves a user from the database based on the request parameters
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user id from the request parameters
	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the user from the database
	user, err := db_controller.GetUser(nil, user_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the user to the client
	httputils.SendJSONResponse(w, user)
}

// GetUserBans retrieves bans for a user from the database based on the request parameters
func GetUserBans(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user id from the request parameters
	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ends_after parameter from the request
	ends_after, err := httputils.RetrieveTimeStampParameter(r, constants.ENDS_AFTER_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the ban types from the request
	ban_types, err := httputils.RetrieveStringListValueParameter(r, constants.BAN_TYPE, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the issuer ids from the request
	issuer_ids, err := httputils.RetrieveIntListValueParameter(r, constants.ISSUER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the base parameters for the request
	order, limit, page, offset := retrieveBaseParams(r)

	// Retrieve the bans from the database
	bans, err := db_controller.GetBans(nil, &db_model.BansGetRequestParams{
		TargetID:  []int{user_id},
		EndsAfter: ends_after,
		Type:      ban_types,
		IssuerID:  issuer_ids,
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

// GetUserMessages retrieves messages for a user from the database based on the request parameters
func GetUserMessages(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user id from the request parameters parameters
	id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the flagged parameter from the request parameters
	flagged, err := httputils.RetrieveBoolParameter(r, constants.FLAGGED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the censored parameter from the request parameters
	censored, err := httputils.RetrieveBoolParameter(r, constants.CENSORED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the removed parameter from the request parameters
	removed, err := httputils.RetrieveBoolParameter(r, constants.REMOVED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the contains parameter from the request parameters
	contains, err := httputils.RetrieveStringListValueParameter(r, constants.CONTAINS_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the base parameters for the request
	order, limit, page, offset := retrieveBaseParams(r)

	// Retrieve the messages
	messages, err := db_controller.GetMessages(nil, &db_model.MessagesGetRequestParams{
		SenderID: []int{id},
		Flagged:  flagged,
		Censored: censored,
		Removed:  removed,
		Contains: contains,
		Order:    order,
		Limit:    limit,
		Page:     page,
		Offset:   offset,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the messages to the client
	httputils.SendJSONResponse(w, messages)
}

// ==================== Update ====================
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	// Retrieve the user id from the request parameters
	user_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Check if the user has permission to update the user
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

	// Admin reserved fields
	admin, err := httputils.RetrievePostFormBoolParameter(r, constants.ADMIN_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	} else if len(admin) > 0 && !user.Admin {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user does not have permission to update admin status for user: "+strconv.Itoa(user_id)))
		return
	}

	// Shared fields
	username, err := httputils.RetrievePostFormStringParameter(r, constants.USERNAME_PARAMETER, true)
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

	user_to_update, err = db_controller.UpdateUser(db, user_to_update, &db_model.UsersRawPatchRequestParams{
		ID:       user_id,
		Username: username,
		Email:    email,
		Password: password,
		Admin:    admin,
		Avatar:   avatar_file,
	})
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

	// Retrieve the user id from the request parameters
	user_to_delete_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
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

	// Retrieve the user to delete
	user_to_delete, err := db_controller.GetUser(db, user_to_delete_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Check if the user has permission to delete the user
	if !db_controller.UserHasPermissionToDeleteUser(user, user_to_delete) {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user does not have permission to delete user"))
		return
	}

	// Delete the user
	err = db_controller.DeleteUser(db, user_to_delete)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "user deleted")
}

// DeleteUsers deletes users from the database based on the request parameters
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	requester, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	// Retrieve the usernames from the request
	usernames, err := httputils.RetrieveStringListValueParameter(r, constants.USERNAME_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the ids from the request
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the emails from the request
	emails, err := httputils.RetrieveStringListValueParameter(r, constants.EMAIL_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the reason from the request
	reason, err := httputils.RetrieveStringParameter(r, constants.REASON_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}

	// Retrieve the base parameters for the request
	order, limit, page, offset := retrieveBaseParams(r)

	// Delete the users
	err = db_controller.DeleteUsers(nil, requester, &db_model.UsersDeleteRequestParams{
		Username: usernames,
		ID:       ids,
		Email:    emails,
		Reason:   reason,
		Order:    order,
		Limit:    limit,
		Page:     page,
		Offset:   offset,
	})
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	httputils.SendSuccessResponse(w, "users deleted")
}
