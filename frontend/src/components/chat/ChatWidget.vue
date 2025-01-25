<!--
This widget is responsible for the chat functionality
Displays the chat messages
Establishes a WebSocket connection to the server
Sends messages to the server
Displays error messages
-->

<script lang="ts">
import { defineComponent, ref, onMounted, onBeforeUnmount } from 'vue';
import SendButtonIcon from '@/components/icons/SendIcon.vue';
import ErrorNotification from '@/components/common/ErrorNotification.vue';
import MessageBubble from '@/components/chat/MessageBubble.vue';
import type { WebsocketDisplayMessage } from '@/constants/types';
import { WEBSOCKET_MESSAGE_TYPES } from '@/constants/types';
import { apiMessageToWebsocketMessage, contentToRawIncomingMessage } from '@/functions/chat';

export default defineComponent({
  name: 'ChatWidget',
  components: {
    SendButtonIcon,
    ErrorNotification,
    MessageBubble,
  },
  setup() {
    const isSubmitting = ref(false);
    const websocketConnection = ref<WebSocket | null>(null);
    const messages = ref<WebsocketDisplayMessage[]>([]);
    const errorMessage = ref<string | null>(null);

    const establishConnection = () => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      const wsUrl = `${protocol}//${host}/chat/ws`;

      websocketConnection.value = new WebSocket(wsUrl);

      websocketConnection.value.onopen = () => {
        console.log('WebSocket connection established.');
      };

      websocketConnection.value.onmessage = (event) => {
        const data: WebsocketDisplayMessage = JSON.parse(event.data);
        if (data.type === WEBSOCKET_MESSAGE_TYPES.DISPLAY) {
          messages.value.push(data);
        } else {
          console.error('Unknown message type:', data.type);
        }
      };

      websocketConnection.value.onclose = (event) => {
        console.log('WebSocket connection closed:', event);
        displayError(`Connection closed: ${event.reason} \nPlease log in again.`);
        setTimeout(establishConnection, 3000);
      };

      websocketConnection.value.onerror = (error) => {
        console.error('WebSocket error:', error);
        displayError('WebSocket encountered an error.');
      };
    };

    const displayError = (message: string, timeout: number = 15000) => {
      errorMessage.value = message;
      setTimeout(() => {
        errorMessage.value = null;
      }, timeout); // Hide after timeout
    };

    const handleSubmit = (event: Event) => {
      event.preventDefault();
      isSubmitting.value = true;

      const input = document.getElementById('chat-input') as HTMLTextAreaElement;
      const message = input.value.trim();

      if (websocketConnection.value && websocketConnection.value.readyState === WebSocket.OPEN && message) {
        websocketConnection.value.send(contentToRawIncomingMessage(message));
        input.value = '';
      } else {
        console.error('WebSocket is not connected.');
        displayError('Unable to send message: WebSocket is not connected.');
      }

      isSubmitting.value = false;
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault(); // Prevent default Enter behavior
        handleSubmit(new Event('submit'));
      }
    };

    onMounted(() => {
      fetch('/api/messages?limit=100')
        .then((response) => response.json())
        .then((data) => {
          messages.value.push(...data.map(apiMessageToWebsocketMessage));
        })
        .catch((error) => {
          console.error('Failed to fetch messages:', error);
          displayError('Failed to fetch messages.');
        });
      establishConnection();
    });

    onBeforeUnmount(() => {
      if (websocketConnection.value) {
        websocketConnection.value.close();
      }
    });

    return {
      isSubmitting,
      messages,
      errorMessage,
      handleSubmit,
      handleKeyDown,
    };
  },
});
</script>

<template>
  <div class="flex flex-col h-full bg-neutral-900">
    <!-- Error Notification -->
    <div v-if="errorMessage" class="flex justify-center mt-2">
      <ErrorNotification v-if="errorMessage" :message="errorMessage" />
    </div>

    <!-- Chat display -->
    <div id="chat-display" class="flex flex-col gap-2 flex-grow overflow-y-auto py-4 px-2">
      <MessageBubble v-for="(message, index) in messages" :key="index" :message="message" />
    </div>

    <!-- Input form -->
    <div class="bg-neutral-950">
      <form class="flex gap-2 w-full my-4 px-2" id="chat-form" @submit="handleSubmit">
        <textarea id="chat-input"
          class="flex-grow rounded-lg py-1 px-2 bg-neutral-800 text-neutral-50 resize-none overflow-y-auto my-auto"
          placeholder="What are we playing, boss?" rows="2" @keydown="handleKeyDown"></textarea>
        <button type="submit" :disabled="isSubmitting"
          class="h-7 w-7 align-middle cursor-pointer bg-neutral-500 text-neutral-800 rounded-r-lg my-auto">
          <SendButtonIcon />
        </button>
      </form>
    </div>
  </div>
</template>
