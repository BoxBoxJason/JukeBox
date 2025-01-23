<!--
This component is responsible for logging in the user
The component sends a POST request to the server to log in a user
The component emits an event to the parent component when the user is successfully signed in
 -->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import ArrowRigtIcon from '@/components/icons/ArrowRightIcon.vue';
import ShowPasswordIcon from '@/components/icons/ShowPasswordIcon.vue';
import HidePasswordIcon from '@/components/icons/HidePasswordIcon.vue';

export default defineComponent({
  name: 'SignInWidget',

  components: {
    ArrowRigtIcon,
    ShowPasswordIcon,
    HidePasswordIcon,
  },

  emits: ['loginSuccess'],

  setup(_, { emit }) {
    const passwordVisible = ref(false);

    // Toggle password visibility
    const togglePasswordVisibility = () => {
      passwordVisible.value = !passwordVisible.value;
    };

    return {
      passwordVisible,
      togglePasswordVisibility,
    };
  },
});
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full">
    <form class="flex flex-col gap-2" action="/api/auth/login" method="POST">
      <div class="auth-input-container">
        <label for="username_or_email" class="auth-input-label">Username or Email</label>
        <input type="text" class="auth-input" placeholder="Enter your username or email" id="username_or_email"
          name="username_or_email" required />
      </div>
      <div class="auth-input-container relative">
        <label for="password" class="auth-input-label">Password</label>
        <div class="relative">
          <input :type="passwordVisible ? 'text' : 'password'" class="auth-input pr-10" placeholder="Password"
            id="password" name="password" required />
          <button type="button"
            class="absolute right-2 top-2/4 transform -translate-y-2/4 text-gray-500 hover:text-primary-500 w-6 h-6"
            @click="togglePasswordVisibility" aria-label="Toggle password visibility">
            <HidePasswordIcon v-if="passwordVisible" />
            <ShowPasswordIcon v-else />
          </button>
        </div>
      </div>
      <button type="submit" class="flex items-center justify-center w-64 p-2 bg-primary-500 text-white rounded-lg">
        <ArrowRigtIcon class="w-6 h-6 mr-2" />
        Sign In
      </button>
    </form>
  </div>
</template>
