package chatwebsocket

import (
	"net/http"
	"time"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins, but customize this for security in production
	}
	clients   = make(map[*websocket.Conn]*db_model.User)
	broadcast = make(chan WebSocketMessage)
)

type WebSocketMessage struct {
	SenderName           string `json:"sender_name"`
	SenderID             int    `json:"sender_id"`
	SenderSubscriberTier int    `json:"sender_subscriber_tier"`
	MessageID            int    `json:"message_id"`
	Content              string `json:"message"`
	Action               string `json:"action"`
	SendTime             string `json:"send_time"`
	EditTime             string `json:"edit_time"`
}

// WebSocket for chat
func ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	websocket_connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		err_msg := "Failed to upgrade connection to websocket"
		logger.Error(err_msg, err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError(err_msg))
		return
	}
	defer websocket_connection.Close()

	clients[websocket_connection] = user
	logger.Info("New client connected", user)

	handleMessages(websocket_connection, user)
}

func handleMessages(websocket_connection *websocket.Conn, user *db_model.User) {
	for {
		var message WebSocketMessage
		err := websocket_connection.ReadJSON(&message)
		if err != nil {
			logger.Error("Failed to read message", err)
			delete(clients, websocket_connection)
			break
		}

		// Determine the action and call the appropriate handler
		switch message.Action {
		case "create":
			handleCreateMessage(message, user, websocket_connection)
		case "delete":
			handleDeleteMessage(message, user, websocket_connection)
		default:
			// If the action is unknown, send an error back to the client
			websocket_connection.WriteJSON(map[string]interface{}{
				"status":  "error",
				"message": "Unknown action",
			})
		}
	}
}

func handleCreateMessage(message WebSocketMessage, sender *db_model.User, websocket_connection *websocket.Conn) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Error("Failed to open database connection", err)
		websocket_connection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to open database connection",
		})
		return
	}
	defer db_model.CloseConnection(db)

	// Create the message
	new_message := db_model.Message{
		Sender:  *sender,
		Content: message.Content,
	}

	err = new_message.CreateMessage(db)
	if err != nil {
		logger.Error("Failed to create message", err)
		websocket_connection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to create message in database",
		})
		return
	}

	// Confirm the message was created
	websocket_connection.WriteJSON(map[string]interface{}{
		"status":     "success",
		"message":    "Message created",
		"message_id": new_message.ID,
	})

	// Fill in the message details
	message.SenderID = sender.ID
	message.MessageID = new_message.ID
	message.SenderSubscriberTier = sender.Subscriber_Tier
	message.SendTime = time.Now().Format(time.RFC3339)
	message.Action = "create"

	// Broadcast the message
	broadcast <- message
}

func handleDeleteMessage(message WebSocketMessage, sender *db_model.User, websocket_connection *websocket.Conn) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Error("Failed to open database connection", err)
		websocket_connection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to open database connection",
		})
		return
	}
	defer db_model.CloseConnection(db)

	// Retrieve the message
	message_to_delete, err := db_model.GetMessageByID(db, message.MessageID)
	if err != nil {
		logger.Error("Failed to retrieve message", err)
		websocket_connection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve message from database",
		})
		return
	}

	// Check if the user is the sender of the message OR an admin
	if message_to_delete.Sender.ID != sender.ID && sender.IsAdmin() {
		logger.Error("User is not the sender of the message or an admin")
		websocket_connection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "User is not the sender of the message or an admin",
		})
		return
	}

	// Delete the message
	err = message_to_delete.DeleteMessage(db)
	if err != nil {
		logger.Error("Failed to delete message", err)
		websocket_connection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to delete message from database",
		})
		return
	}

	// Confirm the message was deleted
	websocket_connection.WriteJSON(map[string]interface{}{
		"status":  "success",
		"message": "Message deleted",
	})

	// Broadcast the message
	message.Action = "delete"
	broadcast <- message
}
