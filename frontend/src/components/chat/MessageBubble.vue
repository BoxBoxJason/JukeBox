<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'MessageBubble',
  props: {
    messageId: {
      type: Number,
      required: true,
    },
    senderId: {
      type: Number,
      required: true,
    },
    senderName: {
      type: String,
      required: true,
    },
    content: {
      type: String,
      required: true,
    },
    date: {
      type: String,
      required: true,
    },
  },
});
</script>

<template>
  <div class="flex flex-col gap-1">
    <div class="flex justify-between">
      <span class="font-bold">{{ senderName }}</span>
      <span class="text-sm">{{ date }}</span>
    </div>
    <div class="p-2 rounded-lg">
      <p>{{ content }}</p>
    </div>
  </div>
</template>
