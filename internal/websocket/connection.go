package websocket

import (
	"net/http"
	"time"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/coder/websocket"
)

// EstablishConnection establishes a websocket connection with the client and listens for incoming messages
// Stores the connection in the Clients struct
func EstablishConnection(w http.ResponseWriter, r *http.Request) {

	// Retrieve the user from the context
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	// Retrieve the access token from the context
	_, ok = r.Context().Value(constants.ACCESS_TOKEN_CONTEXT_KEY).(*db_model.AuthToken)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("access token not found"))
		return
	}

	// Create a websocket connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		logger.Error("failed to upgrade connection to websocket", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("websocket connection upgrade failed"))
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "websocket connection closed")

	// Add the connection to the connection pool
	connectionPool.Add(conn, user)
	defer connectionPool.Remove(conn)

	// Create a ticker to check the access token every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Listen for incoming messages
	for {
		select {
		case <-ticker.C:
			access_token, ok := r.Context().Value(constants.ACCESS_TOKEN_CONTEXT_KEY).(*db_model.AuthToken)
			if !ok || access_token.IsExpired() {
				conn.Close(websocket.StatusPolicyViolation, "access token expired or not found")
				return
			}
		default:
			typ, msg, err := conn.Read(r.Context())
			if err != nil {
				logger.Error("failed to read message from websocket connection", err)
				return
			}
			processed_message, err := processWebsocketMessage(typ, msg, user)
			if err == nil {
				connectionPool.Broadcast(r.Context(), processed_message)
			}
		}
	}
}
