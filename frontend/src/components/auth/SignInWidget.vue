<!--
This component is responsible for logging in the user
The component sends a POST request to the server to log in a user
The component emits an event to the parent component when the user is successfully signed in
 -->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import ArrowRightIcon from '@/components/icons/ArrowRightIcon.vue';
import ShowPasswordIcon from '@/components/icons/ShowPasswordIcon.vue';
import HidePasswordIcon from '@/components/icons/HidePasswordIcon.vue';
import { setIdentity } from '@/functions/auth';

export default defineComponent({
  name: 'SignInWidget',

  components: {
    ArrowRightIcon,
    ShowPasswordIcon,
    HidePasswordIcon,
  },

  emits: ['loginSuccess'],

  setup(_, { emit }) {
    const passwordVisible = ref(false);
    const isSubmitting = ref(false);

    // Toggle password visibility
    const togglePasswordVisibility = () => {
      passwordVisible.value = !passwordVisible.value;
    };

    const clearForm = () => {
      // Clear form fields
      let form = document.getElementById('signin-form') as HTMLFormElement;
      form.reset();
    };

    // Handle form submission
    const handleSubmit = async (event: Event) => {
      event.preventDefault(); // Prevent default form submission behavior
      isSubmitting.value = true;

      const form = event.target as HTMLFormElement;
      const formData = new FormData(form);

      try {
        const response = await fetch('/api/auth/login', {
          method: 'POST',
          body: JSON.stringify(Object.fromEntries(formData.entries())),
          headers: { 'Content-Type': 'application/json' },
        });

        const data = await response.json();

        if (!response.ok) {
          emit('loginSuccess', { success: false, message: data.error });
        } else {
          // Emit success event if login is successful
          setIdentity(data.user_id, data.username);
          emit('loginSuccess', { success: true });
          location.reload();
          clearForm();
        }
      } catch (error: any) {
        emit('loginSuccess', { success: false, message: error.message });
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
    <form id="signin-form" class="flex flex-col gap-2" @submit="handleSubmit">
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
            class="absolute right-2 top-2/4 transform -translate-y-2/4 text-slate-100 hover:text-primary-500 w-6 h-6"
            @click="togglePasswordVisibility" aria-label="Toggle password visibility">
            <HidePasswordIcon v-if="passwordVisible" />
            <ShowPasswordIcon v-else />
          </button>
        </div>
      </div>
      <button type="submit" :disabled="isSubmitting"
        class="flex items-center justify-center w-64 p-2 bg-primary-500 text-white rounded-lg">
        <div class="w-6 h-6 mr-2">
          <ArrowRightIcon />
        </div>
        <span v-if="!isSubmitting">Sign In</span>
        <span v-else>Submitting...</span>
      </button>
    </form>
  </div>
</template>
