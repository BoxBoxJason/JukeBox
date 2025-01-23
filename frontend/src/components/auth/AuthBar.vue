<!--
This component is responsible for the authentication bar.
Manages user authentication state internally using localStorage.
Displays the Sign In and Register buttons if the user is not connected.
Displays the user's username and a logout button if the user is connected.
-->

<script lang="ts">
import { defineComponent, ref, onMounted, computed, onUnmounted } from 'vue';
import { LOCAL_STORAGE_KEYS } from '@/constants/storage';
import { isUserFullySignedIn, autoLogin, fullLogout } from '@/functions/auth';

export default defineComponent({
  name: 'AuthBar',

  emits: ['updateVisibility'],

  setup(_, { emit }) { // Access emit in setup
    const username = ref<string | null>(null);
    const userId = ref<string | null>(null);

    const syncFromLocalStorage = () => {
      username.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
      userId.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
    };

    const handleStorageEvent = (event: StorageEvent) => {
      if (event.key === LOCAL_STORAGE_KEYS.USERNAME || event.key === LOCAL_STORAGE_KEYS.USER_ID) {
        syncFromLocalStorage();
      }
    };

    // Load user data from localStorage on mount
    onMounted(() => {
      // Load username and userId from localStorage
      syncFromLocalStorage();

      // If nothing was loaded, attempt to request auth from refresh token
      if (!isUserFullySignedIn()) {
        autoLogin().then((authSuccess) => {
          if (authSuccess) {
            username.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
            userId.value = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
          }
        })
      }
      window.addEventListener('storage', handleStorageEvent);
    });

    onUnmounted(() => {
      // Cleanup listener
      window.removeEventListener('storage', handleStorageEvent);
    });

    const isSignedIn = computed(() => {
      return isUserFullySignedIn() && username.value !== null && userId.value !== null;
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
      isSignedIn,
    };
  },
});
</script>

<template>
  <nav class="w-full p-4">
    <div class="w-full flex justify-between items-center" v-if="isSignedIn">
      <span>{{ username }}</span>
      <button @click="logout">Logout</button>
    </div>
    <div v-else class="w-full flex justify-end items-center gap-2">
      <button class="auth-button" @click="showSignIn">Sign In</button>
      <button class="auth-button" @click="showRegister">Register</button>
    </div>
  </nav>
</template>
