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
import IconLogout from '../components/icons/Icon_logout.vue'

const text = ref<string>("")

export default defineComponent({
  components: {
    SignIn,
    Register,
    IconSendButton,
    ChatWidget,,
    IconLogout
  },
  setup() {
    const isSignInVisible = ref<boolean>(false) // Typage de isSignInVisible en booléen
    const isRegisterVisible = ref<boolean>(false) // Typage de isSignInVisible en booléen
    const passwordInput = ref<string>("")
    const isConnected = ref(false)
    const errorMessage = ref<string | null>(null)
    const isLogout = ref(false)

    const toggleSignIn = () => {
      // Inverse la valeur de isSignInVisible
      isSignInVisible.value = !isSignInVisible.value
    }

    const toggleRegister = () => {
      // Inverse la valeur de isSignInVisible
      isRegisterVisible.value = !isRegisterVisible.value
    }

    const togglePassword = () => {
      if (passwordInput.value == 'password') {
        passwordInput.value = 'text'
      } else {
        passwordInput.value = 'password'
      }
    }

    const toggleConnection = () => {
      isConnected.value = !isConnected.value
    }

    const logout = async () => {
      isLogout.value = true

      try {
        const response = await fetch('https://localhost:3000/api/auth/logout', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
        })

        if (response.ok) {
          // Succès : Traitez les données de la réponse ou redirigez
          const data = await response.json()
          errorMessage.value = 'Compte correctement deconnecté !'
          toggleConnection()
        } else {
          errorMessage.value = await response.text()
        }
      } catch (err) {
        // Gérer les erreurs réseau
        if (err instanceof TypeError && err.message.includes('Failed to fetch')) {
          errorMessage.value = 'Impossible de se connecter au serveur.'
        } else {
          errorMessage.value = (err as Error).message || 'Erreur réseau inattendue.'
        }
      } finally {
        isLogout.value = false // Débloque le bouton
      }
    }

    return {
      isSignInVisible,
      isRegisterVisible,
      isConnected,
      toggleSignIn,
      toggleRegister,
      togglePassword,
      toggleConnection,
      logout
    }
  }
})
</script>

<template>

  <!-- Page principale -->
  <div class="main-container">

    <!-- Barre de message -->
    <div class="message-input-container">
      <input class="input-prompt" placeholder="What do you want to play ?"/>
      <button class="send-button">
        <IconSendButton />
      </button>
    </div>

    <!-- Barre de connexion si pas connecté -->
    <div v-if="!isConnected" class="login-container">
      <button class="sign-in" @click="toggleSignIn">Sign in</button>
      <button class="register" @click="toggleRegister">Register</button>
    </div>

    <!-- Barre de connexion si connecté -->
    <div v-if="isConnected" class="logged-container">
      <button class="logout" @click="logout">
        <IconLogout />
      </button>
      <h1 class="username"><b>Mathisadi</b></h1>
    </div>

    <!-- Barre d'informations -->
    <div class="info-container">
      <p class="copyright">©2024 JukeBox</p>
      <p class="help"><u>Help ?</u></p>
      <p class="about_us"><u>About us</u></p>
    </div>
  </div>

  <ChatWidget />

  <!-- Barre chat -->
  <div class="left-bar">
      <div class="title-container">
        <p class="title">Jukebox</p>
      </div>
  </div>

  <!-- Pop Up -->
  <SignIn v-if="isSignInVisible" @close="toggleSignIn" @success="toggleConnection"/>
  <Register v-if="isRegisterVisible" @close="toggleRegister" @success="toggleConnection"/>

</template>

<style scoped>

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

.info-container{
  position: absolute;
  font-family: 'Roboto';
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
  font-family: 'Roboto';
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
  font-family: 'Roboto';
  font-size: 16px;
  color:var(--color-text);
  border: 1px solid;
  cursor: pointer;
  border-color: black;
}

.logged-container{
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 5%;
  top: 2%;
}

.logout{
  position: absolute;
  height: 100%;
  width: 3%;
  right: 1%;
  background-color: var(--color-background);
  border: none;
  cursor: pointer;
}

.username{
  position: absolute;
  left: 2%;
  font-family: 'Roboto';
  font-size: 16px;
  color:var(--color-text);
}

.left-bar {
  position: fixed; /* Fixé au côté gauche de l'écran */
  display: flex;
  align-items: center;
  top: 0%;
  left: 0%;
  width: 20%; /* 20% de la largeur de la fenêtre */
  height: 100%; /* 100% de la hauteur de la fenêtre */
  background-color: var(--color-background-soft); /* Couleur de la barre */
}

.title-container{
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  top: 2%;
  height: 5%;
  width: 100%;
}

.title{
  position: absolute;
  padding-bottom: 1.5%;
  font-size: 26px;
  font-family: 'Titan one';
  color: var(--color-heading);
}
</style>

