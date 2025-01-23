<!--
This component is responsible for the registration of the user
The component sends a POST request to the server to create a new user
The component emits an event to the parent component when the user is successfully registered
 -->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import ShowPasswordIcon from '@/components/icons/ShowPasswordIcon.vue';
import HidePasswordIcon from '@/components/icons/HidePasswordIcon.vue';
import ArrowRightIcon from '../icons/ArrowRightIcon.vue';

export default defineComponent({
  name: 'RegisterWidget',

  components: {
    ShowPasswordIcon,
    HidePasswordIcon,
    ArrowRightIcon,
  },

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
    <form class="flex flex-col gap-2" action="/api/user" method="POST">
      <div class="auth-input-container">
        <label for="username" class="auth-input-label">Username</label>
        <input type="text" class="auth-input" placeholder="Enter your username" id="username" name="username"
          required />
      </div>
      <div class="auth-input-container">
        <label for="email" class="auth-input-label">Email</label>
        <input type="email" class="auth-input" placeholder="Enter your email" id="email" name="email" required />
      </div>
      <div class="auth-input-container">
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
        <ArrowRightIcon class="w-6 h-6 mr-2" />
        Register
      </button>
    </form>
  </div>
</template>
