<!--
This component is responsible for toggling between the sign in and register forms.
It uses the SignInWidget and RegisterWidget components.
-->

<script lang="ts">
import { defineComponent, onMounted, onBeforeUnmount, ref } from 'vue'
import SignInWidget from '@/components/auth/SignInWidget.vue'
import RegisterWidget from '@/components/auth/RegisterWidget.vue'
import CrossIcon from '@/components/icons/CrossIcon.vue'

export default defineComponent({
  name: 'AuthSwapper',
  components: {
    SignInWidget,
    RegisterWidget,
    CrossIcon
  },
  props: {
    currentForm: {
      type: String,
      required: true
    },
    isVisible: {
      type: Boolean,
      required: true
    }
  },
  emits: ['updateForm', 'updateVisibility'],
  setup(props, { emit }) {
    const message = ref('')
    const messageType = ref<'success' | 'error' | ''>('')
    const messageTimeout = ref<ReturnType<typeof setTimeout> | null>(null)

    const closeWidget = () => {
      emit('updateVisibility', { visible: false })
    }

    const handleOutsideClick = (event: MouseEvent) => {
      // Logic for handling outside click, (should close the auth-swapper widget)
    }

    const showMessage = (payload: { success: boolean; message: string }) => {
      messageType.value = payload.success ? 'success' : 'error'
      message.value = payload.message

      if (payload.success) {
        emit('updateForm', 'signin') // Switch to sign in form after successful registration
      }

      if (messageTimeout.value) {
        clearTimeout(messageTimeout.value)
      }

      messageTimeout.value = setTimeout(() => {
        message.value = ''
        messageType.value = ''
      }, 10000) // Clear message after 10 seconds
    }

    const handleLoginSuccess = (payload: { success: boolean; message?: string }) => {
      if (payload.success) {
        closeWidget() // Close the widget on successful login
      } else {
        showMessage({ success: false, message: payload.message || 'Login failed' })
      }
    }

    onMounted(() => {
      document.addEventListener('click', handleOutsideClick)
    })

    onBeforeUnmount(() => {
      document.removeEventListener('click', handleOutsideClick)
      if (messageTimeout.value) {
        clearTimeout(messageTimeout.value)
      }
    })

    return {
      closeWidget,
      message,
      messageType,
      showMessage,
      handleLoginSuccess
    }
  }
})
</script>

<template>
  <div v-if="isVisible" class="fixed z-20 top-0 left-0 w-screen h-screen bg-black bg-opacity-85">
    <div
      class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-30 flex flex-col items-center px-10 transform p-4 bg-[var(--color-background)] rounded-lg"
    >
      <button class="absolute top-2 right-2 h-6 w-6" @click="closeWidget">
        <CrossIcon />
      </button>
      <div class="flex gap-4 mt-4 mb-4">
        <button
          class="px-4 py-2 rounded focus:outline-none"
          :class="{
            'bg-[var(--color-hover)] border border-[var(--color-border)] text-[var(--color-chat)] hover:bg-[var(--color-hover)]':
              currentForm === 'signin',
            'bg-[var(--color-background)] border border-[var(--color-border)] text-[var(--color-chat)]':
              currentForm !== 'signin'
          }"
          @click="$emit('updateForm', 'signin')"
        >
          Sign In
        </button>
        <button
          class="px-4 py-2 rounded focus:outline-none"
          :class="{
            'bg-[var(--color-hover)] border border-[var(--color-border)] text-[var(--color-chat)]':
              currentForm === 'register',
            'bg-[var(--color-background)] border border-[var(--color-border)] text-[var(--color-chat)]':
              currentForm !== 'register'
          }"
          @click="$emit('updateForm', 'register')"
        >
          Register
        </button>
      </div>
      <!-- SignInWidget and RegisterWidget -->
      <SignInWidget v-if="currentForm === 'signin'" @loginSuccess="handleLoginSuccess" />
      <RegisterWidget v-else @registerSuccess="showMessage" />
      <!-- Message Display -->
      <div class="h-7 mt-2">
        <div
          v-if="message"
          :class="{
            'text-green-600': messageType === 'success',
            'text-red-600': messageType === 'error'
          }"
        >
          {{ message }}
        </div>
      </div>
    </div>
  </div>
</template>
