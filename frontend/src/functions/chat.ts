import type { APIMessage, WebsocketDisplayMessage } from '@/constants/types'
import { WEBSOCKET_MESSAGE_TYPES } from '@/constants/types'

export function apiMessageToWebsocketMessage(apiMessage: APIMessage): WebsocketDisplayMessage {
  return {
    type: WEBSOCKET_MESSAGE_TYPES.DISPLAY,
    content: apiMessage.content,
    sender: {
      id: apiMessage.sender.id,
      username: apiMessage.sender.username,
      avatar: apiMessage.sender.avatar,
      subscriber_tier: apiMessage.sender.subscriber_tier,
      admin: apiMessage.sender.admin,
    },
    created_at: new Date(apiMessage.created_at),
    modified_at: new Date(apiMessage.modified_at),
    message_id: apiMessage.message_id,
  }
}

export function contentToRawIncomingMessage(content: string): string {
  return JSON.stringify({
    type: WEBSOCKET_MESSAGE_TYPES.RAW_INCOMING_MESSAGE,
    content,
  });
}
