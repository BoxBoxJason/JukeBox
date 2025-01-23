<!--
This component is responsible for the authentication bar.
Manages user authentication state internally using localStorage.
Displays the Sign In and Register buttons if the user is not connected.
Displays the user's username and a logout button if the user is connected.
-->

<script lang="ts">
import { defineComponent, ref, onMounted, onUnmounted } from 'vue';
import LogoutIcon from '@/components/icons/LogoutIcon.vue';
import { LOCAL_STORAGE_KEYS } from '@/constants/storage';
import { autoLogin, fullLogout } from '@/functions/auth';

export default defineComponent({
  name: 'AuthBar',

  components: {
    LogoutIcon,
  },

  emits: ['updateVisibility'],

  setup(_, { emit }) { // Access emit in setup
    const username = ref<string | null>(null);
    const userId = ref<string | null>(null);

    const syncFromLocalStorage = () => {
      username.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
      userId.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
    };

    const handleStorageEvent = (event: CustomEvent) => {
      syncFromLocalStorage();
    };

    // Load user data from localStorage on mount
    onMounted(() => {
      // Load username and userId from localStorage
      syncFromLocalStorage();

      // If nothing was loaded, attempt to request auth from refresh token
      autoLogin().then((authSuccess) => {
        if (authSuccess) {
          username.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
          userId.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
        }
      });

      window.addEventListener('localStorageChange', handleStorageEvent as unknown as EventListener);
    });

    onUnmounted(() => {
      // Cleanup listener
      window.removeEventListener('localStorageChange', handleStorageEvent as unknown as EventListener);
    });

    const logout = async () => {
      const success = await fullLogout();
      if (success) {
        username.value = null;
        userId.value = null;
      }
    };

    const showSignIn = () => {
      emit('updateVisibility', { visible: true, form: 'signin' });
    };

    const showRegister = () => {
      emit('updateVisibility', { visible: true, form: 'register' });
    };

    return {
      username,
      logout,
      showSignIn,
      showRegister,
    };
  },
});
</script>

<template>
  <nav class="w-full p-4">
    <div class="w-full flex justify-between items-center" v-if="username">
      <span>{{ username }}</span>
      <button @click="logout">
        <LogoutIcon class="w-6 h-6" />
      </button>
    </div>
    <div v-else class="w-full flex justify-end items-center gap-2">
      <button class="auth-button" @click="showSignIn">Sign In</button>
      <button class="auth-button" @click="showRegister">Register</button>
    </div>
  </nav>
</template>
