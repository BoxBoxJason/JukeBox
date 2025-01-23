<!--
This component is responsible for toggling between the sign in and register forms.
It uses the SignInWidget and RegisterWidget components.
-->

<script lang="ts">
import { defineComponent, onMounted, onBeforeUnmount } from 'vue';
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
    const closeWidget = () => {
      emit('updateVisibility', { visible: false });
    };

    const handleOutsideClick = (event: MouseEvent) => {
      // if (props.isVisible) {
      //   const target = event.target as HTMLElement;
      //   if (!target.closest('.auth-swapper')) {
      //     closeWidget();
      //   }
      // }
    };

    onMounted(() => {
      document.addEventListener('click', handleOutsideClick);
    });

    onBeforeUnmount(() => {
      document.removeEventListener('click', handleOutsideClick);
    });

    return {
      closeWidget,
    };
  },
});
</script>

<template>
  <div v-if="isVisible"
    class="auth-swapper z-20 flex flex-col items-center px-12 fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 p-4 rounded-lg shadow-lg">
    <button class="absolute top-1 right-1" @click="closeWidget">
      <CrossIcon class="h-6 w-6" />
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
    <SignInWidget v-if="currentForm === 'signin'" />
    <RegisterWidget v-else />
  </div>
</template>
