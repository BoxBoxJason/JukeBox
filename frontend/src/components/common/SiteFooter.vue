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
  <footer class="w-full flex justify-between items-center px-12 py-4 text-[var(--color-text)] text-sm">
    <span>©{{ currentYear() }} <a href="/">JukeBox</a></span>
    <a href="/help"><u>Help ?</u></a>
    <a href="/about"><u>About Us</u></a>
  </footer>
</template>
