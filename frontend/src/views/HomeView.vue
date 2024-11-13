TODO: gérer le zoom
TODO: Ajouter l'element qui change au son de la musique
TODO: Gérer les langues ?
TODO: Transformer contribuer et help pour pouvoir aller sur une autre page
TODO: Ajouter bouton sign in et register
TODO: Boite message

<script lang="ts">
import { defineComponent, ref } from 'vue'
import SignIn from '../components/SignInWidget.vue'

const text = ref<string>("")

export default defineComponent({
  components: {
    SignIn
  },
  setup() {
    const isSignInVisible = ref<boolean>(false) // Typage de isSignInVisible en booléen

    function toggleSignIn() {
      // Inverse la valeur de isSignInVisible
      isSignInVisible.value = !isSignInVisible.value
    }

    return {
      isSignInVisible,
      toggleSignIn
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
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="3.12 1 17 24" fill="none">
          <path d="M7 14.5L12 9.5L17 14.5" stroke=var(--color-background-mute) stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>  
        </svg>
      </button>
    </div>

    <div class="login-container">       
      <button class="sign-in" @click="toggleSignIn">Sign In</button>
      <button class="register">Register</button>
    </div>

    <div class="info-container">
      <p class="copyright">©2024 JukeBox</p>
      <p class="help"><u>Help ?</u></p>
      <p class="contribuer"><u>Contribuer</u></p>
    </div>
  </div>

  <!-- Barre chat -->
  <div class="left-bar">
      <p class="title">Jukebox</p>
  </div>

  <!-- Pop Up -->
  <SignIn v-if="isSignInVisible" @close="toggleSignIn"/>
  
  
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

.contribuer{
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

