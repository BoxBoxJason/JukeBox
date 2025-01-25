export const WEBSOCKET_MESSAGE_TYPES = {
  // Display message (sent by the server to all clients)
  DISPLAY: 'display',
  // Raw incoming message (sent to the server as is)
  RAW_INCOMING_MESSAGE: 'raw_incoming_message',
}

// WebSocket broadcast message types
export type WebsocketDisplayMessage = {
  type: string,
  content: string,
  sender: {
    id: number,
    username: string,
    avatar: string,
    subscriber_tier: number,
    admin: boolean,
  }
  created_at: Date | string,
  modified_at: Date | string,
  message_id: number,
}

// Websocket raw incoming message
export type WebsocketRawIncomingMessage = {
  type: string,
  content: string,
}

export type APIMessage = {
  censored: boolean,
  flagged: boolean,
  removed: boolean,
  content: string,
  message_id: number,
  sender: APIUser,
  created_at: string,
  modified_at: string,
}

export type APIUser = {
  admin: boolean,
  avatar: string,
  created_at: string,
  modified_at: string,
  id: number,
  banned: boolean,
  minutes_listened: number,
  subscriber_tier: number,
  total_contributions: number,
  username: string,
}
