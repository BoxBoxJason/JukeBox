<script lang="ts">
import { defineComponent, handleError, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import IconConnect from './icons/Icon_connect.vue'
import IconQuit from './icons/Icon_quit.vue'
import IconPassword from './icons/Icon_password.vue'

export default defineComponent({
  name: 'loginForm',

  components: {
    IconConnect,
    IconQuit,
    IconPassword
  },

  props: {
  message: {
    type: String,
    default: null,
  },
},

  emits: ['close','success'],

  setup(_, { emit }) {
    // Variables
    const nomIcon = ref('show')
    const typeInput = ref('password')
    const formData = reactive({ username_or_email: '', password: '' })
    const errorMessage = ref<string | null>(null)
    const isSubmitting = ref(false)

    // Fonctions
    const closepopup = () => {
      // Émet l'événement close pour le composant parent
      emit('close')
    }

    const togglePasswordButton = () => {
      if (nomIcon.value == 'show') {
        nomIcon.value = 'hide'
        typeInput.value = 'text'
      } else {
        nomIcon.value = 'show'
        typeInput.value = 'password'
      }
    }

    const handleSubmit = async () => {
      errorMessage.value = null // Réinitialise le message d'erreur
      isSubmitting.value = true // Bloque le bouton pendant la soumission

      try {
        const response = await fetch("https://localhost:3000/api/auth/login", {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData)
        })

        if (response.ok) {
          // Succès : Traitez les données de la réponse ou redirigez
          const data = await response.json()

          if (data.access_token) {
            // Stocke le token dans le localStorage
            localStorage.setItem('authToken', data.access_token)

            // Émet l'événement 'success' avec le username et le token
            emit('success', { username: data.Username, token: data.access_token })
            emit('close')
          } else {
            errorMessage.value = "Token manquant dans la réponse."
          }
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
        isSubmitting.value = false // Débloque le bouton
        formData.password = ''
        formData.username_or_email = ''
      }
    }

    return {
      errorMessage,
      nomIcon,
      typeInput,
      formData,
      closepopup,
      togglePasswordButton,
      handleSubmit,
      isSubmitting
    }
  }
})
</script>

<template>
  <div class="widget-container">
    <form class="widget-sign-in" @submit.prevent="handleSubmit">
      <button class="quit" @click="closepopup">
        <IconQuit />
      </button>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div v-if="message" class="alert alert-success"> {{ message }}</div>
      <h1 class="title">Welcome Back on Jukebox!</h1>
      <h2 class="username-text">Username</h2>
      <input
        id="email"
        v-model="formData.username_or_email"
        class="username-input"
        placeholder="Type your username here !"
        required
      />
      <h3 class="password-text">Password</h3>
      <input
        id="password"
        v-model="formData.password"
        :type="typeInput"
        class="password-input"
        placeholder="Type your password here !"
        required
      />
      <h4 class="forget-pass"><u>Forget password ?</u></h4>
      <button class="password-button" @click="togglePasswordButton" type="button">
        <IconPassword :name="nomIcon" />
      </button>
      <button class="connect-button" v-bind:disabled="isSubmitting" type="submit">
        <IconConnect />
      </button>
    </form>
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
  background-color: rgba(0, 0, 0, 0.8);
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

.title {
  position: relative;
  top: 12%;
  font-family: 'Roboto';
  font-size: 30px;
}

.username-input {
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

.username-text {
  position: absolute;
  top: 31%;
  font-family: 'Roboto';
  font-size: 14px;
  left: 6%;
}

.password-input {
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

.password-text {
  position: absolute;
  top: 50%;
  font-family: 'Roboto';
  font-size: 14px;
  left: 6%;
}

.forget-pass {
  position: absolute;
  top: 69%;
  left: 6%;
  font-family: 'Roboto';
  font-size: 14px;
}

.connect-button {
  position: absolute;
  width: 10%;
  height: 10%;
  top: 82%;
  background-color: var(--color-background);
  border: none;
}

.quit {
  position: absolute;
  width: 6%;
  height: 5%;
  top: 2%;
  right: 3%;
  background-color: var(--color-background);
  border: none;
}

.password-button {
  position: absolute;
  width: 8%;
  height: 5%;
  top: 58.5%;
  right: 8%;
  background-color: transparent;
  border: none;
}
</style>
