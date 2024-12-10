TODO: gérer le zoom
TODO: Ajouter l'element qui change au son de la musique
TODO: Gérer les langues ?
TODO: Boite message

<script lang="ts">
import { defineComponent, ref } from 'vue'
import SignIn from '../components/SignInWidget.vue'
import Register from '../components/RegisterWidget.vue'
import IconSendButton from '../components/icons/Icon_send_button.vue'
import ChatWidget from '../components/ChatWidget.vue'

const text = ref<string>("")

export default defineComponent({
  components: {
    SignIn,
    Register,
    IconSendButton,
    ChatWidget,
  },
  setup() {
    const isSignInVisible = ref<boolean>(false) // Typage de isSignInVisible en booléen
    const isRegisterVisible = ref<boolean>(false) // Typage de isSignInVisible en booléen
    const passwordInput = ref<string>("")

    function toggleSignIn() {
      // Inverse la valeur de isSignInVisible
      isSignInVisible.value = !isSignInVisible.value
    }

    function toggleRegister() {
      // Inverse la valeur de isSignInVisible
      isRegisterVisible.value = !isRegisterVisible.value
    }

    function togglePassword() {
      if (passwordInput.value == 'password') {
        passwordInput.value = 'text'
      } else {
        passwordInput.value = 'password'
      }
    }

    return {
      isSignInVisible,
      isRegisterVisible,
      toggleSignIn,
      toggleRegister,
      togglePassword
    }
  }
})
</script>

<template>

  <!-- Page principale -->
  <div class="main-container">

    <div class="message-input-container">
      <input class="input-prompt" placeholder="What do you want to play ?"/>
      <button class="send-button">
        <IconSendButton />
      </button>
    </div>

    <div class="login-container">       
      <button class="sign-in" @click="toggleSignIn">Sign In</button>
      <button class="register" @click="toggleRegister">Register</button>
    </div>

    <div class="info-container">
      <p class="copyright">©2024 JukeBox</p>
      <p class="help"><u>Help ?</u></p>
      <p class="about_us"><u>About us</u></p>
    </div>
  </div>

  <ChatWidget />
  
  <!-- Barre chat -->
  <div class="left-bar">
      <p class="title">Jukebox</p>
  </div>

  <!-- Pop Up -->
  <SignIn v-if="isSignInVisible" @close="toggleSignIn"/>
  <Register v-if="isRegisterVisible" @close="toggleRegister"/>
  
</template>

<style scoped>

.left-bar {
  position: fixed; /* Fixé au côté gauche de l'écran */
  top: 0%;
  left: 0%;
  width: 20%; /* 20% de la largeur de la fenêtre */
  height: 100%; /* 100% de la hauteur de la fenêtre */
  background-color: var(--color-background-soft); /* Couleur de la barre */
}

.main-container {
  position: fixed;
  top: 0;
  left: 20%;
  width: 80%; /* Largeur de 20% de la fenêtre */
  height: 100%; /* Hauteur de 100% de la fenêtre */
  background-color: var(--color-background);
  display: flex; /* Active Flexbox */
  align-items: center; /* Centre verticalement */
  justify-content: center;
}

.message-input-container {
  position: absolute;
  display: flex;
  align-items: center;
  width: 65%;
  height: 6%;
}

.input-prompt {
  width: 100%;
  height: 100%;
  border-radius: 100px;
  border: 1px solid;
  border-color: var(--color-background-mute);
  background-color: var(--color-background-mute);
  font-size: 16px;
  color: var(--color-text);
  padding-left: 3%;
}

.send-button {
  display: flex;
  align-items: center;
  justify-content: center;
  position: absolute; /* Positionnement absolu pour superposer */
  right:  2%; /* Distance du bord droit */
  height: 4vh;
  width: 4vh;
  color: var(--color-background);
  border: none;
  border-radius: 100%;
  cursor: pointer;
  background-color: var(--color-bouton-color);
}

.send-button svg {
  width: 100%; /* Ajuste la taille pour qu'elle soit proportionnelle */
  height: 100%;
}

.title{
  margin-top: 5%;
  text-align: center;
  font-size: 25px;
  font-family: 'Gagalin';
  color: var(--color-heading);
}

.info-container{
  position: absolute;
  font-family: 'Open Sans';
  font-size: 14px;
  width: 100%;
  height: 5%;
  bottom: 0%;
  left: 0%;
}

.copyright{
  position: absolute;
  left: 2%;
}

.help{
  position: absolute;
  right: 2%;

}

.about_us{
  position: absolute;
  right: 9%;
}

.login-container{
  position: absolute;
  width: 20%;
  height: 5%;
  top: 2%;
  right: 2%;
}

.sign-in{
  position: absolute;
  width: 46%;
  height: 100%;
  left: 0%;
  border-radius: 100px;
  background-color: var(--color-login);
  font-family: 'Open Sans';
  font-size: 16px;
  color:var(--color-text);
  border: 1px solid;
  cursor: pointer;
  border-color: black;
}

.register{
  position: absolute;
  width: 46%;
  height: 100%;
  right: 0%;
  border-radius: 100px;
  background-color: var(--color-register);
  font-family: 'Open Sans';
  font-size: 16px;
  color:var(--color-text);
  border: 1px solid;
  cursor: pointer;
  border-color: black;
}

</style>

