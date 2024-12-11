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

	websocketConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errMsg := "Failed to upgrade connection to websocket"
		logger.Error(errMsg, err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError(errMsg))
		return
	}
	defer websocketConnection.Close()

	clients[websocketConnection] = user
	logger.Info("New client connected", user)

	go handleBroadcast()
	handleMessages(websocketConnection, user)
}

func handleMessages(websocketConnection *websocket.Conn, user *db_model.User) {
	defer func() {
		delete(clients, websocketConnection)
		logger.Info("Client disconnected", user)
	}()

	websocketConnection.SetReadDeadline(time.Now().Add(60 * time.Second))
	websocketConnection.SetPongHandler(func(string) error {
		websocketConnection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message WebSocketMessage
		if err := websocketConnection.ReadJSON(&message); err != nil {
			logger.Error("Failed to read message", err)
			break
		}

		// Determine the action and call the appropriate handler
		switch message.Action {
		case "create":
			handleCreateMessage(message, user, websocketConnection)
		case "delete":
			handleDeleteMessage(message, user, websocketConnection)
		default:
				websocketConnection.WriteJSON(map[string]interface{}{
					"status":  "error",
					"message": "Unknown action",
				})
		}
	}
}

func handleCreateMessage(message WebSocketMessage, sender *db_model.User, websocketConnection *websocket.Conn) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Error("Failed to open database connection", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to open database connection",
		})
		return
	}
	defer db_model.CloseConnection(db)

	// Create the message
	newMessage := db_model.Message{
		Sender:  sender,
		Content: message.Content,
	}

	if err := newMessage.CreateMessage(db); err != nil {
		logger.Error("Failed to create message", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to create message in database",
		})
		return
	}

	// Confirm the message was created
	websocketConnection.WriteJSON(map[string]interface{}{
		"status":     "success",
		"message":    "Message created",
		"message_id": newMessage.ID,
	})

	// Fill in the message details
	message.SenderID = sender.ID
	message.MessageID = newMessage.ID
	message.SenderSubscriberTier = sender.Subscriber_Tier
	message.SendTime = time.Now().Format(time.RFC3339)
	message.Action = "create"

	// Broadcast the message
	broadcast <- message
}

func handleDeleteMessage(message WebSocketMessage, sender *db_model.User, websocketConnection *websocket.Conn) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Error("Failed to open database connection", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to open database connection",
		})
		return
	}
	defer db_model.CloseConnection(db)

	// Retrieve the message
	messageToDelete, err := db_model.GetMessageByID(db, message.MessageID)
	if err != nil {
		logger.Error("Failed to retrieve message", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve message from database",
		})
		return
	}

	// Check if the user is the sender of the message OR an admin
	if messageToDelete.Sender.ID != sender.ID && !sender.Admin {
		logger.Error("User is not the sender of the message or an admin")
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "User is not the sender of the message or an admin",
		})
		return
	}

	// Delete the message
	if err := messageToDelete.DeleteMessage(db); err != nil {
		logger.Error("Failed to delete message", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to delete message from database",
		})
		return
	}

	// Confirm the message was deleted
	websocketConnection.WriteJSON(map[string]interface{}{
		"status":  "success",
		"message": "Message deleted",
	})

	// Broadcast the message
	message.Action = "delete"
	broadcast <- message
}

func handleBroadcast() {
	for {
		message := <-broadcast
		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				logger.Error("Failed to send message to client", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
