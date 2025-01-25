<!--
The message bubble component is used to display messages in the chat room.
It is a simple component that displays the sender's username, the message content, and the time the message was sent.
It also includes a method to format the display time and a method to determine the margin of the message bubble based on the sender's ID.
-->

<script lang="ts">
import { defineComponent } from 'vue';
import type { WebsocketDisplayMessage } from '@/constants/types';
import { formatDisplayTime, fullFormatDisplayTime } from '@/functions/time';
import { getUserId } from '@/functions/auth';

export default defineComponent({
  name: 'MessageBubble',

  props: {
    message: {
      type: Object as () => WebsocketDisplayMessage,
      required: true,
    },
  },

  methods: {
    formatDisplayTime,
    fullFormatDisplayTime,
    senderIsUser() {
      const user_id = getUserId();
      return this.message.sender.id === user_id;
    },
  },
});
</script>

<template>
  <div v-if="senderIsUser()" class="flex flex-col gap-1 w-[90%] p-2 rounded-lg bg-emerald-600 ml-auto">
    <div class="flex justify-between">
      <span class="font-bold text-emerald-900">{{ message.sender.username }}</span>
      <span class="text-sm text-emerald-900" :title="fullFormatDisplayTime(message.created_at)">{{
        formatDisplayTime(message.created_at) }}</span>
    </div>
    <div class="px-2 text-neutral-50">
      <p>{{ message.content }}</p>
    </div>
  </div>

  <div v-else class="flex flex-col gap-1 w-[90%] p-2 rounded-lg bg-neutral-500 mr-auto">
    <div class="flex justify-between">
      <span class="font-bold text-neutral-100">{{ message.sender.username }}</span>
      <span class="text-sm text-neutral-100" :title="fullFormatDisplayTime(message.created_at)">{{
        formatDisplayTime(message.created_at) }}</span>
    </div>
    <div class="px-2 text-neutral-50">
      <p>{{ message.content }}</p>
    </div>
  </div>
</template>
