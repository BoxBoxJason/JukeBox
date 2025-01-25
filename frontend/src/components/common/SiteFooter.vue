<!--
This component is the footer of the site
-->

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'SiteFooter',

  methods: {
    currentYear() {
      return new Date().getFullYear();
    },
  }
});

</script>
<template>
  <footer class="w-full flex justify-between items-center px-12 py-4 bg-gray-800 text-white bg-slate-800">
    <span>©{{ currentYear() }} JukeBox</span>
    <a href="/help">Help ?</a>
    <a href="/about">About Us</a>
  </footer>
</template>
