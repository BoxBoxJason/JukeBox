<!--
This component is responsible for displaying error notifications
It will display an error message for a specified duration
-->

<script lang="ts">
import { defineComponent, } from 'vue';

export default defineComponent({
  name: 'ErrorNotification',
  props: {
    message: {
      type: String,
      required: true,
    },
  },
});
</script>

<template>
  <span class="bg-red-900 text-white text-center p-4 border-2 border-red-500 rounded-lg mx-auto">
    {{ message }}
  </span>
</template>
