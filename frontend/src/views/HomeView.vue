<script lang="ts">
import { defineComponent, ref, onBeforeUnmount } from 'vue';
import { connectWebSocket, sendMessage, disconnectWebSocket } from '../services/WebSocketService'; // Assurez-vous d'importer WebSocketService
import SignIn from '../components/SignInWidget.vue'
import Register from '../components/RegisterWidget.vue'
import IconSendButton from '../components/icons/Icon_send_button.vue'
import IconLogout from '../components/icons/Icon_logout.vue'

export default defineComponent({
  components: {
    SignIn,
    Register,
    IconSendButton,
    IconLogout
  },
  setup() {
    const isSignInVisible = ref<boolean>(false)
    const isRegisterVisible = ref<boolean>(false)
    const isConnected = ref(false)
    const errorMessage = ref<string | null>(null)
    const isLogout = ref(false)
    const messages = ref<{ username: string; text: string; date: string }[]>([])
    const newMessage = ref<string>("")
    const username = ref<string>("")
    const socketUrl = 'wss://localhost:3000/ws/chat'; // Utiliser wss:// pour WebSocket sécurisé

    // Connexion WebSocket et gestion des messages
    const handleIncomingMessage = (message: any) => {
      messages.value.push({
        username: message.sender_name,
        text: message.message,
        date: new Date(message.send_time).toLocaleString(),
      });
    };

    // Connexion WebSocket après l'authentification
    const toggleConnection = (user: string) => {
      if (!isConnected.value && user) {
        username.value = user;
        isConnected.value = true;
        errorMessage.value = null;

        // Connecte-toi au WebSocket uniquement après la connexion de l'utilisateur
        connectWebSocket(socketUrl, handleIncomingMessage);
      } else {
        username.value = "";
        isConnected.value = false;
        errorMessage.value = null;

        // Déconnecte le WebSocket si l'utilisateur se déconnecte
        disconnectWebSocket();
      }
    };

    const logout = async () => {
      isLogout.value = true;
      try {
        const response = await fetch('https://localhost:3000/api/auth/logout', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
        });
        if (response.ok) {
          errorMessage.value = null;
          toggleConnection(""); // Déconnexion de WebSocket ici aussi
        } else {
          errorMessage.value = 'Erreur lors de la déconnexion : ' + await response.text();
        }
      } catch (err) {
        errorMessage.value = err instanceof TypeError && err.message.includes('Failed to fetch')
          ? 'Impossible de se connecter au serveur.'
          : (err as Error).message || 'Erreur réseau inattendue.';
      } finally {
        isLogout.value = false;
      }
    };

    const addMessage = () => {
      if (!isConnected.value) {
        errorMessage.value = "Vous devez être connecté pour envoyer un message.";
        return;
      }

      const trimmedMessage = newMessage.value.trim();
      if (trimmedMessage === "") {
        errorMessage.value = "Le message ne peut pas être vide.";
        return;
      }

      const message = {
        sender_name: username.value,
        message: trimmedMessage,
        action: 'create',
        send_time: new Date().toISOString()
      };

      // Envoie le message via WebSocket
      sendMessage(message);

      // Réinitialiser le champ du message
      newMessage.value = "";
      errorMessage.value = null;
    };

    // Déconnexion proprement à la destruction du composant
    onBeforeUnmount(() => {
      disconnectWebSocket();
    });

    return {
      isSignInVisible,
      isRegisterVisible,
      isConnected,
      toggleSignIn: () => isSignInVisible.value = !isSignInVisible.value,
      toggleRegister: () => isRegisterVisible.value = !isRegisterVisible.value,
      toggleConnection,
      logout,
      messages,
      newMessage,
      username,
      addMessage,
      errorMessage
    }
  }
});
</script>


<template>
  <div class="main-container">
    <div class="message-input-container">
      <input
        class="input-prompt"
        v-model="newMessage"
        @keyup.enter="addMessage"
        :placeholder="isConnected ? 'What do you want to play ?' : 'Connectez-vous pour envoyer un message'"
      />
      <button class="send-button" @click="addMessage">
        <IconSendButton />
      </button>
    </div>
    <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>

    <div v-if="!isConnected" class="login-container">       
      <button class="sign-in" @click="toggleSignIn">Sign in</button>
      <button class="register" @click="toggleRegister">Register</button>
    </div>

    <div v-if="isConnected" class="logged-container">
      <button class="logout" @click="logout">
        <IconLogout />
      </button>
      <h1 class="username"><b>{{ username }}</b></h1>
    </div>

    <div class="info-container">
      <p class="copyright">©2024 JukeBox</p>
      <p class="help"><u>Help ?</u></p>
      <p class="about_us"><u>About us</u></p>
    </div>
  </div>

  <div class="left-bar">
    <div class="title-container">
      <p class="title">Jukebox</p>
    </div>

    <div class="messages-container">
      <ul class="messages-list">
        <li v-for="(message, index) in messages" :key="index" class="message-item">
          <p><strong>{{ message.username }}</strong> - {{ message.date }}</p>
          <p>{{ message.text }}</p>
        </li>
      </ul>
    </div>
  </div>

  <SignIn v-if="isSignInVisible" @close="toggleSignIn" @success="toggleConnection"/>
  <Register v-if="isRegisterVisible" @close="toggleRegister" @success="toggleConnection"/>
</template>

<style scoped>
.main-container {
  position: fixed;
  top: 0;
  left: 20%;
  width: 80%;
  height: 100%;
  background-color: var(--color-background);
  display: flex;
  align-items: center;
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
  position: absolute;
  right:  2%;
  height: 4vh;
  width: 4vh;
  color: var(--color-background);
  border: none;
  border-radius: 100%;
  cursor: pointer;
  background-color: var(--color-bouton-color);
}

.send-button svg {
  width: 100%;
  height: 100%;
}

.error-message {
  color: red;
  font-size: 14px;
  margin-top: 5px;
  text-align: center;
}

.info-container {
  position: absolute;
  font-family: 'Roboto';
  font-size: 14px;
  width: 100%;
  height: 5%;
  bottom: 0%;
  left: 0%;
}

.copyright {
  position: absolute;
  left: 2%;
}

.help {
  position: absolute;
  right: 2%;
}

.about_us {
  position: absolute;
  right: 9%;
}

.login-container {
  position: absolute;
  width: 20%;
  height: 5%;
  top: 2%;
  right: 2%;
}

.sign-in {
  position: absolute;
  width: 46%;
  height: 100%;
  left: 0%;
  border-radius: 100px;
  background-color: var(--color-login);
  font-family: 'Roboto';
  font-size: 16px;
  color: var(--color-text);
  border: 1px solid;
  cursor: pointer;
  border-color: black;
}

.register {
  position: absolute;
  width: 46%;
  height: 100%;
  right: 0%;
  border-radius: 100px;
  background-color: var(--color-register);
  font-family: 'Roboto';
  font-size: 16px;
  color: var(--color-text);
  border: 1px solid;
  cursor: pointer;
  border-color: black;
}

.logged-container {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 5%;
  top: 2%;
}

.logout {
  position: absolute;
  height: 100%;
  width: 3%;
  right: 1%;
  background-color: var(--color-background);
  border: none;
  cursor: pointer;
}

.username {
  position: absolute;
  left: 2%;
  font-family: 'Roboto';
  font-size: 16px;
  color: var(--color-text);
}

.left-bar {
  position: fixed;
  display: flex;
  align-items: center;
  top: 0%;
  left: 0%;
  width: 20%;
  height: 100%;
  background-color: var(--color-background-soft);
}

.title-container {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  top: 2%;
  height: 5%;
  width: 100%;
}

.title {
  position: absolute;
  padding-bottom: 1.5%;
  font-size: 26px;
  font-family: 'Titan one';
  color: var(--color-heading);
}

.messages-container {
  position: relative;
  height: 85%; 
  width: 100%; 
  margin-top: 5%; 
  overflow-y: auto; 
  background-color: var(--color-background-mute); 
  padding: 10px; 
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1); 
  border-radius: 8px; 
}

.message-item {
  display: flex; 
  flex-direction: column; 
  padding: 8px; 
  margin-bottom: 10px; 
  background-color: var(--color-background-soft); 
  border-radius: 8px; 
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); 
  color: var(--color-text); 
  font-family: 'Roboto', sans-serif;
  font-size: 14px; 
  word-wrap: break-word;
}

.messages-container::-webkit-scrollbar {
  width: 8px;
}

.messages-container::-webkit-scrollbar-track {
  background-color: var(--color-background);
}

.messages-container::-webkit-scrollbar-thumb {
  background-color: var(--color-background-mute);
  border-radius: 4px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background-color: var(--color-background-soft);
}
</style>
