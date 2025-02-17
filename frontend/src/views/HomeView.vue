<!--
Website landing page, handles the login and registration of users
Displays the chat room and the music player. Allows users to communicate with each other
Displays navigation / legal footer
-->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import SiteFooter from '@/components/common/SiteFooter.vue';
import AuthBar from '@/components/auth/AuthBar.vue';
import AuthSwapper from '@/components/auth/AuthSwapper.vue';
import ChatWidget from '@/components/chat/ChatWidget.vue';
import SoundWave from '@/components/common/SoundWave.vue';

export default defineComponent({
  components: {
    AuthBar,
    AuthSwapper,
    ChatWidget,
    SiteFooter,
    SoundWave,
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
    <!-- Sidebar -->
    <aside class="fixed left-0 top-0 h-screen z-10 overflow-y-auto"
      style="max-width: 500px; min-width: 300px; width: 25vw">
      <ChatWidget />
    </aside>

    <!-- Main Content -->
    <main class="flex flex-col w-3/4 min-h-screen" style="margin-left: calc(max(300px, min(25vw, 500px)))">
      <!-- Auth Bar -->
      <AuthBar @updateVisibility="updateVisibility" />
      <!-- Auth Swapper -->
      <AuthSwapper :isVisible="isAuthVisible" :currentForm="currentForm"
        @updateVisibility="(payload) => updateVisibility(payload)" @updateForm="currentForm = $event" />
      <!-- SoundWave -->
      <div class="flex-grow flex flex-col items-center text-center justify-center">
        <SoundWave />
      </div>
      <!-- Footer -->
      <SiteFooter />
    </main>
  </div>
</template>
