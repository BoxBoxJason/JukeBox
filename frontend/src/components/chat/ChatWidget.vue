<!--
This widget is responsible for the chat functionality
Displays the chat messages
-->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import SendButtonIcon from '@/components/icons/SendIcon.vue';

export default defineComponent({
  name: 'ChatWidget',
  components: {
    SendButtonIcon,
  },

  setup() {
    const isSubmitting = ref(false);
    const websocketConnection = ref<WebSocket | null>(null);

    const establishConnection = () => {
      // Establish a websocket connection
    };

    // Handle form submission
    const handleSubmit = async (event: Event) => {
      event.preventDefault(); // Prevent default form submission behavior
      isSubmitting.value = true;

      const form = event.target as HTMLFormElement;

      try {
        // Send the message to the server through websocket connection
      } catch (error) {
        console.error('An error occurred while sending the message:', error);
      } finally {
        isSubmitting.value = false;
      }
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        // Prevent a newline from being added
        event.preventDefault();
        // Trigger form submission
        const form = document.getElementById('chat-form') as HTMLFormElement;
        if (form) {
          form.dispatchEvent(new Event('submit'));
        }
      }
    };

    return {
      isSubmitting,
      handleSubmit,
      handleKeyDown,
    };

  },

});
</script>

<template>
  <div class="flex flex-col h-full bg-neutral-900">
    <div id="chat-display" class="flex flex-col gap-2 flex-grow overflow-y-auto">
      <!-- Chat messages will be displayed here -->
    </div>
    <div class="bg-neutral-950">
      <form class="flex gap-2 w-full my-4 px-2" id="chat-form" @submit="handleSubmit">
        <textarea id="chat-input"
          class="flex-grow rounded-lg py-1 px-2 bg-neutral-800 text-neutral-50 resize-none overflow-y-auto my-auto"
          placeholder="What are we playing boss ?" rows="2" @keydown="handleKeyDown"></textarea>
        <button type="submit" :disabled="isSubmitting"
          class="h-7 w-7 align-middle cursor-pointer bg-neutral-500 text-neutral-800 rounded-r-lg my-auto">
          <SendButtonIcon />
        </button>
      </form>
    </div>
  </div>
</template>
