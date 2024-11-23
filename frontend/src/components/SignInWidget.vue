<script lang="ts">
import { defineComponent, ref } from 'vue'
import IconConnect from "./icons/Icon_connect.vue"
import IconQuit from "./icons/Icon_quit.vue"
import IconPassword from './icons/Icon_password.vue';

export default defineComponent({
  components: {
    IconConnect,
    IconQuit,
    IconPassword
  },
  emits: ['close'],

  data() {
    return{
      nomIcon: "show",
      typeInput: "password",
    }
  },

  methods: {
    closepopup() {
      // Émet l'événement `close` pour le composant parent
      this.$emit('close');
    },

    togglePasswordButton() {
      if (this.nomIcon == "show"){
        this.nomIcon = "hide"
        this.typeInput = "text"
      }
      else{
        this.nomIcon = "show"
        this.typeInput = "password"
      }
    },
  },
})
</script>

<template>
  <div class="widget-container">
    <div class="widget-sign-in">
      <button class="quit" @click="closepopup">
        <IconQuit />
      </button>
      <h1 class="title">Welcome back on JukeBox !</h1>
      <h2 class="username-text">Username</h2>
      <input class="username-input" placeholder="Type your username here !"/>
      <h3 class="password-text">Password</h3>
      <input :type="typeInput" class="password-input" placeholder="Type your password here !"/>
      <h4 class="forget-pass"><u>Forget password ?</u></h4>
      <button class="password-button" @click="togglePasswordButton">
        <IconPassword :name="nomIcon"/>
      </button>
      <button class="connect-button">
        <IconConnect />
      </button>
    </div>
  </div>
</template>

<style scoped>
.widget-container {
  position: fixed;
  top: 0%;
  left: 0%;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center; /* Centre verticalement */
  justify-content: center; /* Centre horizontalement, optionnel */
  background-color: rgba(0,0,0, 0.8);
}

.widget-sign-in {
  position: absolute;
  height: 65%;
  width: 30%;
  border-radius: 7px;
  display: flex;
  justify-content: center;
  background-color: var(--color-background);
}

.title{
  position: relative;
  top: 12%;
  font-family: 'Open Sans';
  font-size: 30px;
}

.username-input{
  position: absolute;
  width: 90%;
  height: 8%;
  top: 38%;
  background-color: var(--color-background-mute);
  border: none;
  border-radius: 7px;
  padding-left: 2%;
  color: var(--color-text);
}

.username-text{
  position: absolute;
  top: 31%;
  font-family: 'Open Sans';
  font-size: 14px;
  left: 6%;
}

.password-input{
  position: absolute;
  width: 90%;
  height: 8%;
  top: 57%;
  background-color: var(--color-background-mute);
  border: none;
  border-radius: 7px;
  padding-left: 2%;
  color: var(--color-text);
}

.password-text{
  position: absolute;
  top: 50%;
  font-family: 'Open Sans';
  font-size: 14px;
  left: 6%;
}

.forget-pass{
  position: absolute;
  top: 69%;
  left: 6%;
  font-family: 'Open Sans';
  font-size: 14px;
}

.connect-button{
  position: absolute;
  width: 10%;
  height: 10%;
  top: 82%;
  background-color: var(--color-background);
  border: none;
}

.quit{
  position: absolute;
  width: 6%;
  height: 5%;
  top: 2%;
  right: 3%;
  background-color: var(--color-background);
  border: none;
}

.password-button{
  position: absolute;
  width: 8%;
  height: 5%;
  top: 58.5%;
  right: 8%;
  background-color: transparent;
  border: none;
}
</style>

