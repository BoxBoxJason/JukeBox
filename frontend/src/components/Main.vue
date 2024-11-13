<template>
  <div>
    <button @click="showRegisterWidget">Register</button>
    <button @click="showSignInWidget">Sign In</button>

    <RegisterWidget v-if="isRegisterWidgetVisible" />
    <SignInWidget v-if="isSignInWidgetVisible" />
  </div>

</template>

<script>
import RegisterWidget from './RegisterWidget.vue'
import SignInWidget from './SignInWidget.vue'

export default {
  components: {
    RegisterWidget,
    SignInWidget
  },
  data() {
    return {
      isRegisterWidgetVisible: false,
      isSignInWidgetVisible: false
    }
  },
  methods: {
    showRegisterWidget() {
      this.isRegisterWidgetVisible = true
    },
    showSignInWidget() {
      this.isSignInWidgetVisible = true
    }
  },
}
</script>

