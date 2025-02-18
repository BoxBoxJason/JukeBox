<!--
Website about page
Explanation of our project, our goals, the team, and the technology stack
Displays navigation / legal footer
-->

<script lang="ts">
import { defineComponent, ref } from 'vue'
import HomeIcon from '@/components/icons/HomeIcon.vue'

export default defineComponent({
  components: {
    HomeIcon
  },

  setup() {
    const isAuthVisible = ref(false)
    const currentForm = ref('signin')

    const updateVisibility = ({ visible, form }: { visible: boolean; form?: string }) => {
      if (visible !== undefined) {
        isAuthVisible.value = visible
      }
      if (form) {
        currentForm.value = form
      }
    }

    return {
      isAuthVisible,
      currentForm,
      updateVisibility
    }
  }
})
</script>

<template>
  <div class="flex">
    <!-- Main Content -->
    <main class="flex flex-col gap-4 mx-8 my-8 text-[var(--color-text-2)] min-h-screen">
      <!-- Header -->
      <div class="flex flex-row items-center gap-auto pr-4">
        <h1 class="flex-none text-3xl">
          <span class="text-[var(--color-heading)]">Juke</span
          ><span class="text-[var(--color-heading-2)]">Box</span>
        </h1>
        <div class="flex-grow"></div>
        <a
          href="/"
          class="flex flex-row p-2 rounded-full border bg-[var(--color-button-color)] flex-none text-[var(--color-text-2)] h-full items-center gap-2 hover:text-[var(--color-heading)]"
          ><HomeIcon />Go Home</a
        >
      </div>
      <hr />

      <!-- Header Description -->
      <p>
        JukeBox is an interactive music creation platform that allows users to create music together
        by interacting with an AI model in real time.
      </p>
      <p>This project is powered by three distinct AI models:</p>
      <br />

      <!-- Models Description -->
      <h2 class="text-xl">Music generation</h2>
      <hr />
      <!-- Music generation AI model -->
      <p>
        The music generation AI model is
        <a
          href="https://forum.ircam.fr/projects/detail/rave-vst/"
          class="text-[var(--color-heading)]"
          >RAVE</a
        >.
      </p>
      <p>
        This model is an auto encoder, meaning that is takes sound(s) as input and generates a new
        sound as output. It can generate music in real time based on the user's commands. It fits
        our use case perfectly, as it can generate music in real time, and the music's parameters
        can be controlled by the user's commands.
      </p>
      <p>
        You can check
        <a href="https://github.com/acids-ircam/RAVE" class="text-[var(--color-heading)]"
          >their GitHub repository</a
        >
        for more information.
      </p>
      <br />

      <!-- User requests processing -->
      <h2 class="text-xl">Filtering</h2>
      <hr />
      <p>
        The users chat messages are analyzed by an AI model that can filter out messages that are
        not related to music creation, and can also understand the user's requests and translate
        them into music commands. This model's mission is to protect the music generation AI model
        from being spammed with irrelevant messages. It also recognizes malicious intents and can
        ban users from the platform.
      </p>
      <p>
        The filter model is a custom trained version of
        <a
          href="https://huggingface.co/docs/transformers/model_doc/distilbert"
          class="text-[var(--color-heading)]"
          >DistilBERT</a
        >
        trained on a dataset of music-related messages.
      </p>
      <br />

      <!-- User requests processing -->
      <h2 class="text-xl">Reformulation</h2>
      <hr />
      <p>
        The reformulation model converts a user request into a music command that can be understood
        by the music generation AI model.
      </p>

      <!-- Team -->
      <h2 class="text-xl">The team</h2>
      <hr />
      <div class="flex flex-row gap-20">
        <figure class="flex flex-col items-center gap-2">
          <img src="@/components/img/Mathis.jpg" class="rounded-full h-40" />
          <figcaption>Mathis Front</figcaption>
        </figure>
        <figure class="flex flex-col items-center gap-2">
          <img src="@/components/img/Thomas.jpg" class="rounded-full h-40" />
          <figcaption>Thomas Front</figcaption>
        </figure>
        <figure class="flex flex-col items-center gap-2">
          <img src="@/components/img/Lockin.jpg" class="rounded-full h-40" />
          <figcaption>Lockin PDG</figcaption>
        </figure>
        <figure class="flex flex-col items-center gap-2">
          <img src="@/components/img/Sofiane.jpg" class="rounded-full h-40" />
          <figcaption>Sofiane Back</figcaption>
        </figure>
      </div>
    </main>
  </div>
</template>
