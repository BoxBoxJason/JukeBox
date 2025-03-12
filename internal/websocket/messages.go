package websocket

import (
	"encoding/json"
	"time"

	db_controller "github.com/boxboxjason/jukebox/internal/controller"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/coder/websocket"
)

const (
	MESSAGE_TYPE_DISPLAY = "display"
	RAW_INCOMING_MESSAGE = "raw_incoming_message"
)

type SenderWebSocket struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Avatar         string `json:"avatar"`
	SubscriberTier int    `json:"subscriber_tier"`
	Admin          bool   `json:"admin"`
}

type WebSocketMessage struct {
	Type       string          `json:"type"`
	Content    string          `json:"content"`
	Sender     SenderWebSocket `json:"sender"`
	CreatedAt  time.Time       `json:"created_at"`
	ModifiedAt time.Time       `json:"modified_at"`
	MessageID  int             `json:"message_id"`
}

type WebsocketRawIncomingMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func processWebsocketMessage(message_type websocket.MessageType, message []byte, sender *db_model.User) ([]byte, error) {
	var processed_message []byte
	if message_type == websocket.MessageText {
		var incoming_message WebsocketRawIncomingMessage
		// Unmarshal the (json) message
		err := json.Unmarshal(message, &incoming_message)
		if err != nil {
			return nil, err
		}
		if incoming_message.Type == RAW_INCOMING_MESSAGE {
			db_message, err := db_controller.CreateMessage(nil, &db_model.MessagesPostRequestParams{
				Sender:  sender,
				Message: incoming_message.Content,
			})
			if err != nil {
				return nil, err
			} else {
				db_message.Sender = sender
				go addMessage(db_message.Content)
			}

			// Marshal the display message
			processed_message, err = json.Marshal(WebSocketMessage{
				Type:    MESSAGE_TYPE_DISPLAY,
				Content: db_message.Content,
				Sender: SenderWebSocket{
					ID:             sender.ID,
					Username:       sender.Username,
					Avatar:         sender.Avatar,
					SubscriberTier: sender.Subscriber_Tier,
					Admin:          sender.Admin,
				},
				CreatedAt:  db_message.CreatedAt,
				ModifiedAt: db_message.ModifiedAt,
				MessageID:  db_message.ID,
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, httputils.NewBadRequestError("invalid message type")
	}
	return processed_message, nil
}
