<!--
Website help page
Explanation of how to use the website
Displays navigation / legal footer
-->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import SiteFooter from '@/components/common/SiteFooter.vue';
import AuthBar from '@/components/auth/AuthBar.vue';
import AuthSwapper from '@/components/auth/AuthSwapper.vue';

export default defineComponent({
  components: {
    AuthBar,
    AuthSwapper,
    SiteFooter,
  },

  setup() {
    const isAuthVisible = ref(false);
    const currentForm = ref('signin');

    const updateVisibility = ({ visible, form }: { visible: boolean; form?: string }) => {
      if (visible !== undefined) {
        isAuthVisible.value = visible;
      }
      if (form) {
        currentForm.value = form;
      }
    };

    return {
      isAuthVisible,
      currentForm,
      updateVisibility,
    };
  },
});
</script>

<template>
  <div class="flex">
    <!-- Main Content -->
    <main class="flex flex-col min-h-screen">
      <!-- Auth Bar -->
      <AuthBar @updateVisibility="updateVisibility" />
      <!-- Auth Swapper -->
      <AuthSwapper :isVisible="isAuthVisible" :currentForm="currentForm"
        @updateVisibility="(payload) => updateVisibility(payload)" @updateForm="currentForm = $event" />
      <!-- Push Footer to Bottom -->
      <div class="flex-grow"></div>
      <!-- Footer -->
      <SiteFooter />
    </main>
  </div>
</template>
