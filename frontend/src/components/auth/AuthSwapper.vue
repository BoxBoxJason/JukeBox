<!--
This component is responsible for toggling between the sign in and register forms.
It uses the SignInWidget and RegisterWidget components.
-->

<script lang="ts">
import { defineComponent, onMounted, onBeforeUnmount, ref } from 'vue';
import SignInWidget from '@/components/auth/SignInWidget.vue';
import RegisterWidget from '@/components/auth/RegisterWidget.vue';
import CrossIcon from '@/components/icons/CrossIcon.vue';

export default defineComponent({
  name: 'AuthSwapper',
  components: {
    SignInWidget,
    RegisterWidget,
    CrossIcon,
  },
  props: {
    currentForm: {
      type: String,
      required: true,
    },
    isVisible: {
      type: Boolean,
      required: true,
    },
  },
  emits: ['updateForm', 'updateVisibility'],
  setup(props, { emit }) {
    const message = ref('');
    const messageType = ref<'success' | 'error' | ''>('');
    const messageTimeout = ref<ReturnType<typeof setTimeout> | null>(null);

    const closeWidget = () => {
      emit('updateVisibility', { visible: false });
    };

    const handleOutsideClick = (event: MouseEvent) => {
      // Logic for handling outside click, (should close the auth-swapper widget)
    };

    const showMessage = (payload: { success: boolean; message: string }) => {
      messageType.value = payload.success ? 'success' : 'error';
      message.value = payload.message;

      if (payload.success) {
        emit('updateForm', 'signin'); // Switch to sign in form after successful registration
      }

      if (messageTimeout.value) {
        clearTimeout(messageTimeout.value);
      }

      messageTimeout.value = setTimeout(() => {
        message.value = '';
        messageType.value = '';
      }, 10000); // Clear message after 10 seconds
    };

    const handleLoginSuccess = (payload: { success: boolean; message?: string }) => {
      if (payload.success) {
        closeWidget(); // Close the widget on successful login
      } else {
        showMessage({ success: false, message: payload.message || 'Login failed' });
      }
    };

    onMounted(() => {
      document.addEventListener('click', handleOutsideClick);
    });

    onBeforeUnmount(() => {
      document.removeEventListener('click', handleOutsideClick);
      if (messageTimeout.value) {
        clearTimeout(messageTimeout.value);
      }
    });

    return {
      closeWidget,
      message,
      messageType,
      showMessage,
      handleLoginSuccess,
    };
  },
});
</script>

<template>
  <div v-if="isVisible"
    class="auth-swapper z-20 flex flex-col items-center px-12 fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 p-4 bg-slate-800 rounded-lg shadow-lg">
    <button class="absolute top-1 right-1 text-slate-100 h-6 w-6" @click="closeWidget">
      <CrossIcon />
    </button>
    <div class="flex gap-4 mb-4">
      <button class="px-4 py-2 rounded focus:outline-none" :class="{
        'bg-blue-500 text-white': currentForm === 'signin',
        'bg-gray-200 text-gray-800 hover:bg-blue-300 hover:text-white': currentForm !== 'signin',
      }" @click="$emit('updateForm', 'signin')">
        Sign In
      </button>
      <button class="px-4 py-2 rounded focus:outline-none" :class="{
        'bg-blue-500 text-white': currentForm === 'register',
        'bg-gray-200 text-gray-800 hover:bg-blue-300 hover:text-white': currentForm !== 'register',
      }" @click="$emit('updateForm', 'register')">
        Register
      </button>
    </div>
    <!-- SignInWidget and RegisterWidget -->
    <SignInWidget v-if="currentForm === 'signin'" @loginSuccess="handleLoginSuccess" />
    <RegisterWidget v-else @registerSuccess="showMessage" />
    <!-- Message Display -->
    <div class="h-7 mt-2">
      <div v-if="message" :class="{
        'text-green-600': messageType === 'success',
        'text-red-600': messageType === 'error',
      }">
        {{ message }}
      </div>
    </div>
  </div>
</template>
