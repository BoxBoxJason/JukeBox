package chatwebsocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/gorilla/websocket"
)

// Clients struct pour gérer les connexions de manière thread-safe
type Clients struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]*db_model.User
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[*websocket.Conn]*db_model.User),
	}
}

func (c *Clients) Add(conn *websocket.Conn, user *db_model.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[conn] = user
}

func (c *Clients) Remove(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, conn)
}

func (c *Clients) Broadcast(message WebSocketMessage) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for client := range c.clients {
		if err := client.WriteJSON(message); err != nil {
			logger.Error("Failed to send message to client", err)
			client.Close()
			// Retrait du client après erreur d'écriture
			go func(cl *websocket.Conn) {
				clients.Remove(cl)
			}(client)
		}
	}
}

var (
	upgrader = websocket.Upgrader{
		// CheckOrigin doit être ajusté pour plus de sécurité en production
		CheckOrigin: func(r *http.Request) bool {
			// Exemple : autoriser uniquement certaines origines
			// origin := r.Header.Get("Origin")
			// return origin == "https://votre-domaine.com"
			return true
		},
	}

	clients   = NewClients() // Structure thread-safe pour les clients
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

// ChatWebSocket gère l'upgrade HTTP vers WebSocket et l'inscription du client
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

	clients.Add(websocketConnection, user)
	logger.Info("New client connected", user)

	// Configuration du PongHandler pour prolonger la connexion
	websocketConnection.SetReadDeadline(time.Now().Add(60 * time.Second))
	websocketConnection.SetPongHandler(func(string) error {
		websocketConnection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Gestion des messages reçus
	go handleMessages(websocketConnection, user)
}

// handleMessages lit les messages entrants du client et les traite
func handleMessages(websocketConnection *websocket.Conn, user *db_model.User) {
	defer func() {
		clients.Remove(websocketConnection)
		websocketConnection.Close()
		logger.Info("Client disconnected", user)
	}()

	for {
		var message WebSocketMessage
		if err := websocketConnection.ReadJSON(&message); err != nil {
			// Erreur de lecture : déconnexion du client
			logger.Error("Failed to read message", err)
			break
		}

		// Analyse de l'action
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

	newMessage := db_model.Message{
		Sender:  sender,
		Content: message.Content,
	}

	if err := newMessage.CreateMessage(db); err != nil {
		logger.Error("Failed to create message in DB", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to create message in database",
		})
		return
	}

	// Confirmation côté client
	websocketConnection.WriteJSON(map[string]interface{}{
		"status":     "success",
		"message":    "Message created",
		"message_id": newMessage.ID,
	})

	// Préparation du message à diffuser
	message.SenderID = sender.ID
	message.MessageID = newMessage.ID
	message.SenderSubscriberTier = sender.Subscriber_Tier
	message.SendTime = time.Now().Format(time.RFC3339)
	message.Action = "create"

	broadcast <- message
}

func handleDeleteMessage(message WebSocketMessage, sender *db_model.User, websocketConnection *websocket.Conn) {
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

	messageToDelete, err := db_model.GetMessageByID(db, message.MessageID)
	if err != nil {
		logger.Error("Failed to retrieve message", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve message from database",
		})
		return
	}

	// Vérification autorisation
	if messageToDelete.Sender.ID != sender.ID && !sender.Admin {
		logger.Error("User not allowed to delete message")
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "User is not the sender of the message or an admin",
		})
		return
	}

	if err := messageToDelete.DeleteMessage(db); err != nil {
		logger.Error("Failed to delete message", err)
		websocketConnection.WriteJSON(map[string]interface{}{
			"status":  "error",
			"message": "Failed to delete message from database",
		})
		return
	}

	// Confirmation de suppression
	websocketConnection.WriteJSON(map[string]interface{}{
		"status":  "success",
		"message": "Message deleted",
	})

	message.Action = "delete"
	broadcast <- message
}

// HandleBroadcast écoute le channel broadcast et envoie à tous les clients connectés
func HandleBroadcast() {
	for {
		message := <-broadcast
		clients.Broadcast(message)
	}
}
