<template>
  <div>
    <h3>Chat Messages</h3>
    <ul>
      <li v-for="message in messages" :key="message.message_id">
        <strong>{{ message.sender_name }}:</strong> {{ message.content }}
      </li>
    </ul>
    <input v-model="newMessage" placeholder="Type your message..." />
    <button @click="sendMessage">Send</button>
  </div>
</template>

<script>
import { connectWebSocket, sendMessage, disconnectWebSocket } from '@/services/websocket';

export default {
  data() {
    return {
      messages: [],
      newMessage: '',
    };
  },
  methods: {
    handleWebSocketMessage(message) {
      if (message.action === "create") {
        this.messages = [...this.messages, message].sort((a, b) => a.message_id - b.message_id);
      } else if (message.action === "delete") {
        this.messages = this.messages.filter(msg => msg.message_id !== message.message_id);
      }
    },
    sendMessage() {
      if (this.newMessage.trim() !== '') {
        sendMessage({ action: "create", content: this.newMessage });
        this.newMessage = '';
      }
    },
  },
  mounted() {
    connectWebSocket('ws://localhost:8080/ws/chat', this.handleWebSocketMessage);
  },
  beforeDestroy() {
    disconnectWebSocket();
  },
};
</script>
