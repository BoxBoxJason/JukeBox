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

  emits: ['registerSuccess'],

  components: {
    ShowPasswordIcon,
    HidePasswordIcon,
    ArrowRightIcon,
  },

  setup(_, { emit }) {
    const passwordVisible = ref(false);
    const isSubmitting = ref(false);

    // Toggle password visibility
    const togglePasswordVisibility = () => {
      passwordVisible.value = !passwordVisible.value;
    };

    const clearForm = () => {
      // Clear form fields
      let form = document.getElementById('register-form') as HTMLFormElement;
      form.reset();
    };

    // Handle form submission
    const handleSubmit = async (event: Event) => {
      event.preventDefault(); // Prevent default form submission behavior
      isSubmitting.value = true;

      const form = event.target as HTMLFormElement;
      const formData = new FormData(form);

      try {
        const response = await fetch('/api/users', {
          method: 'POST',
          body: JSON.stringify(Object.fromEntries(formData.entries())),
          headers: { 'Content-Type': 'application/json' },
        });

        if (!response.ok) {
          const errorData = await response.json();
          emit('registerSuccess', { success: false, message: errorData.error });
        }

        // Emit success event if the request is successful
        emit('registerSuccess', { success: true, message: 'User registered successfully ! Please log in' });
        clearForm();
      } catch (error: any) {
        emit('registerSuccess', { success: false, message: error.message });
      } finally {
        isSubmitting.value = false;
      }
    };

    return {
      passwordVisible,
      togglePasswordVisibility,
      handleSubmit,
      isSubmitting,
    };
  },
});
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full">
    <form id="register-form" class="flex flex-col gap-2" @submit="handleSubmit">
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
            class="absolute right-2 top-2/4 transform -translate-y-2/4 hover:text-primary-500 w-6 h-6"
            @click="togglePasswordVisibility" aria-label="Toggle password visibility">
            <HidePasswordIcon class="cursor-pointer" v-if="passwordVisible" />
            <ShowPasswordIcon class="cursor-pointer" v-else />
          </button>
        </div>
      </div>
      <button type="submit" :disabled="isSubmitting"
        class="flex items-center justify-center w-64 p-2 text-[var(--color-chat)] rounded-lg cursor-pointer">
        <div class="w-6 h-6 mr-2">
          <ArrowRightIcon />
        </div>
        <span v-if="!isSubmitting">Register</span>
        <span v-else>Submitting...</span>
      </button>
    </form>
  </div>
</template>
