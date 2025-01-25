<!--
This component is responsible for displaying error notifications
It will display an error message for a specified duration
-->

<script lang="ts">
import { defineComponent, ref, watch } from 'vue';

export default defineComponent({
  name: 'ErrorNotification',
  props: {
    message: {
      type: String,
      required: true,
    },
    duration: {
      type: Number,
      default: 15000, // Default to 15 seconds
    },
  },
  setup(props) {
    const visible = ref(false);

    watch(
      () => props.message,
      (newMessage) => {
        if (newMessage) {
          visible.value = true;
          setTimeout(() => {
            visible.value = false;
          }, props.duration);
        }
      }
    );

    return {
      visible,
    };
  },
});
</script>

<template>
  <span v-if="visible" class="bg-red-900 text-white text-center p-4 border-2 border-red-500 rounded-lg mx-auto">
    {{ message }}
  </span>
</template>
