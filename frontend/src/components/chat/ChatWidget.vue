<!--
This widget is responsible for the chat functionality
Displays the chat messages
-->

<script lang="ts">
import { defineComponent } from 'vue';
import SendButtonIcon from '@/components/icons/SendIcon.vue';

export default defineComponent({
  name: 'ChatWidget',
  components: {
    SendButtonIcon,
  },
});
</script>

<template>
  <div class="flex flex-col h-full bg-neutral-900">
    <div id="chat-display" class="flex flex-col gap-2 flex-grow overflow-y-auto">

    </div>
    <div class="flex gap-2 w-full my-6 h-7 px-2">
      <input id="chat-input" class="flex-grow rounded-lg" placeholder="Type your message here" />
      <button id="chat-send" class="h-7 w-7 align-middle cursor-pointer bg-neutral-500 text-neutral-800 rounded-r-lg">
        <SendButtonIcon />
      </button>
    </div>
  </div>
</template>
