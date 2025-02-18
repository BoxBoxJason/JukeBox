<!--
This widget is responsible for the chat functionality
Displays the chat messages
Establishes a WebSocket connection to the server
Sends messages to the server
Displays error messages
-->

<script lang="ts">
import { defineComponent, ref, onMounted, onBeforeUnmount } from 'vue'
import SendButtonIcon from '@/components/icons/SendIcon.vue'
import ErrorNotification from '@/components/common/ErrorNotification.vue'
import MessageBubble from '@/components/chat/MessageBubble.vue'
import type { WebsocketDisplayMessage } from '@/constants/types'
import { WEBSOCKET_MESSAGE_TYPES } from '@/constants/types'
import { apiMessageToWebsocketMessage, contentToRawIncomingMessage } from '@/functions/chat'
import { isUserConnected } from '@/functions/auth'

const PLEASE_LOG_IN_ERROR = 'Please log in to use the chat.'

export default defineComponent({
  name: 'ChatWidget',
  components: {
    SendButtonIcon,
    ErrorNotification,
    MessageBubble
  },
  setup() {
    const isSubmitting = ref(false)
    const websocketConnection = ref<WebSocket | null>(null)
    const messages = ref<WebsocketDisplayMessage[]>([])
    const errorMessage = ref<string | null>(null)

    const displayError = (message: string, duration: number = 15000) => {
      errorMessage.value = message
      if (duration > 0) {
        setTimeout(() => {
          errorMessage.value = null
        }, duration)
      }
    }

    const establishConnection = () => {
      if (isUserConnected()) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const host = window.location.host
        const wsUrl = `${protocol}//${host}/chat/ws`

        websocketConnection.value = new WebSocket(wsUrl)

        websocketConnection.value.onopen = () => {
          console.log('WebSocket connection established.')
        }

        websocketConnection.value.onmessage = (event) => {
          const data: WebsocketDisplayMessage = JSON.parse(event.data)
          if (data.type === WEBSOCKET_MESSAGE_TYPES.DISPLAY) {
            messages.value.push(data)
          } else {
            console.error('Unknown message type:', data.type)
          }
        }

        websocketConnection.value.onclose = (event) => {
          console.log('WebSocket connection closed:', event)
          // If the user is not authenticated, show "Please log in" permanently.
          if (!isUserConnected()) {
            displayError(PLEASE_LOG_IN_ERROR, 0)
          } else {
            // Otherwise, display the error (for a short time) and retry the connection.
            displayError(`Connection closed: ${event.reason}\nPlease log in again.`)
            setTimeout(establishConnection, 3000)
          }
        }

        websocketConnection.value.onerror = (error) => {
          console.error('WebSocket error:', error)
          displayError('WebSocket encountered an error.')
        }
      } else {
        displayError(PLEASE_LOG_IN_ERROR, 0)
        setTimeout(establishConnection, 3000)
      }
    }

    const handleSubmit = (event: Event) => {
      event.preventDefault()
      isSubmitting.value = true

      const input = document.getElementById('chat-input') as HTMLTextAreaElement
      const message = input.value.trim()

      if (
        websocketConnection.value &&
        websocketConnection.value.readyState === WebSocket.OPEN &&
        message
      ) {
        websocketConnection.value.send(contentToRawIncomingMessage(message))
        input.value = ''
      } else {
        console.error('WebSocket is not connected.')
        displayError('Unable to send message: WebSocket is not connected.')
      }

      isSubmitting.value = false
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        handleSubmit(new Event('submit'))
      }
    }

    onMounted(() => {
      fetch('/api/messages?limit=100')
        .then((response) => response.json())
        .then((data) => {
          messages.value.push(...data.map(apiMessageToWebsocketMessage))
        })
        .catch((error) => {
          console.error('Failed to fetch messages:', error)
          displayError('Failed to fetch messages.')
        })
      establishConnection()
    })

    onBeforeUnmount(() => {
      if (websocketConnection.value) {
        websocketConnection.value.close()
      }
    })

    return {
      isSubmitting,
      messages,
      errorMessage,
      handleSubmit,
      handleKeyDown
    }
  }
})
</script>

<template>
  <div class="flex flex-col h-full bg-[var(--color-background-soft)] z-50">
    <!-- Title -->
    <h1 class="text-3xl text-center p-4">
      <span class="text-[var(--color-heading)]">Juke</span><span class="text-[var(--color-heading-2)]">Box</span>
    </h1>

    <!-- Error Notification -->
    <div v-if="errorMessage" class="flex justify-center mt-2">
      <ErrorNotification v-if="errorMessage" :message="errorMessage" />
    </div>

    <!-- Chat display -->
    <div id="chat-display" class="flex flex-col gap-2 grow overflow-y-auto py-4 px-2">
      <MessageBubble v-for="(message, index) in messages" :key="index" :message="message" />
    </div>

    <!-- Input form -->
    <div class="bg-neutral-950">
      <form class="flex gap-2 w-full my-4 px-2" id="chat-form" @submit="handleSubmit">
        <textarea id="chat-input"
          class="flex-grow rounded-lg border border-[var(--color-border)] bg-[var(--color-background-mute)] resize-none overflow-y-auto my-auto placeholder-[var(--color-chat)] py-1 px-2"
          placeholder="What are we playing, boss?" rows="2" @keydown="handleKeyDown"></textarea>
        <button type="submit" :disabled="isSubmitting"
          class="h-7 w-7 align-middle cursor-pointer bg-transparent text-[var(--color-background-soft)] my-auto">
          <SendButtonIcon />
        </button>
      </form>
    </div>
  </div>
</template>
