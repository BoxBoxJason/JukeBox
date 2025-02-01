package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/coder/websocket"
)

const (
	PING_INTERVAL       = 30 * time.Second
	PONG_TIMEOUT        = 10 * time.Second
	AUTH_CHECK_INTERVAL = 5 * time.Minute
)

// EstablishConnection establishes a websocket connection with the client and listens for incoming messages
func EstablishConnection(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	_, ok = r.Context().Value(constants.ACCESS_TOKEN_CONTEXT_KEY).(*db_model.AuthToken)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("access token not found"))
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		logger.Error("failed to upgrade connection to websocket", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("websocket connection upgrade failed"))
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "websocket connection closed")

	connectionPool.Add(conn, user)
	defer connectionPool.Remove(conn)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	pongReceived := make(chan bool, 1)

	// Start authentication check goroutine
	go monitorAuthToken(ctx, cancel, conn, r)

	// Start ping-pong mechanism goroutine
	go managePingPong(ctx, cancel, conn, pongReceived)

	// Start listening for messages in a separate goroutine
	go listenForMessages(ctx, conn, user)

	// Block until context is canceled
	<-ctx.Done()
}

// monitorAuthToken checks the validity of the access token periodically
func monitorAuthToken(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, r *http.Request) {
	ticker := time.NewTicker(AUTH_CHECK_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			accessToken, ok := r.Context().Value(constants.ACCESS_TOKEN_CONTEXT_KEY).(*db_model.AuthToken)
			if !ok || accessToken.IsExpired() {
				logger.Error("Access token expired or not found, closing connection")
				conn.Close(websocket.StatusPolicyViolation, "access token expired or not found")
				cancel() // Cancel all goroutines
				return
			}
		}
	}
}

// managePingPong handles the ping-pong mechanism to keep the connection alive
func managePingPong(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, pongReceived chan bool) {
	pingTicker := time.NewTicker(PING_INTERVAL)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			if err := conn.Write(ctx, websocket.MessageText, []byte("ping")); err != nil {
				logger.Error("Failed to send ping, closing connection", err)
				conn.Close(websocket.StatusInternalError, "ping failed")
				cancel()
				return
			}

			select {
			case <-time.After(PONG_TIMEOUT):
				logger.Error("Pong timeout, closing connection")
				conn.Close(websocket.StatusGoingAway, "pong not received")
				cancel()
				return
			case <-pongReceived:
				// Pong received, continue
			}
		}
	}
}

// listenForMessages handles incoming messages from the websocket
func listenForMessages(ctx context.Context, conn *websocket.Conn, user *db_model.User) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			typ, msg, err := conn.Read(ctx)
			if err != nil {
				continue // continue the loop to keep the connection alive
			}
			processedMessage, err := processWebsocketMessage(typ, msg, user)
			if err == nil {
				connectionPool.Broadcast(ctx, processedMessage)
			}
		}
	}
}
