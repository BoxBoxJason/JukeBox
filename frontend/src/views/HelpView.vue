<!--
Website help page
Explanation of how to use the website
Displays navigation / legal footer
-->

<script lang="ts">
import { defineComponent, ref } from 'vue';
import SiteFooter from '@/components/common/SiteFooter.vue';
import AuthBar from '@/components/auth/AuthBar.vue';
import AuthSwapper from '@/components/auth/AuthSwapper.vue';

export default defineComponent({
  components: {
    AuthBar,
    AuthSwapper,
    SiteFooter,
  },

  setup() {
    const isAuthVisible = ref(false);
    const currentForm = ref('signin');

    const updateVisibility = ({ visible, form }: { visible: boolean; form?: string }) => {
      if (visible !== undefined) {
        isAuthVisible.value = visible;
      }
      if (form) {
        currentForm.value = form;
      }
    };

    return {
      isAuthVisible,
      currentForm,
      updateVisibility,
    };
  },
});
</script>

<template>
  <div class="flex max-w-7xl mx-auto min-h-screen items-center justify-center">
    <!-- Main Content -->
    <main class="flex flex-col gap-4 mx-8 my-8 text-[var(--color-text-2)] min-h-screen">
      <!-- Header -->
      <div class="flex flex-row items-center gap-auto pr-4">
        <h1 class="flex-none text-3xl">
          <span class="text-[var(--color-heading)]">Juke</span><span class="text-[var(--color-heading-2)]">Box</span>
        </h1>
        <div class="grow"></div>
        <a href="/"
          class="flex flex-row p-2 rounded-full border bg-[var(--color-button-color)] flex-none text-[var(--color-text-2)] h-full items-center gap-2 hover:text-[var(--color-heading)]">
          <HomeIcon />Go Home
        </a>
      </div>
      <hr />

    <!-- Header Description -->
    <p>
        JukeBox is an interactive music creation platform that allows users to create music together
        by interacting with an AI model in real time.
    </p>

    <!-- Creating an Account -->
    <h2 class="text-xl">Creating an Account</h2>
      <hr />
      <p>
          <p>Before you can start creating music, you need to create an account.</p>
          <li>Click on the Sign Up button on the homepage.</li>
          <li>Enter your email address, username, and password.</li>
          <li>Log in to start using JukeBox!</li>
      </p>
      <br />
    <!-- Exploring the Interface -->
    <h2 class="text-xl">Exploring the Interface</h2>
      <hr />
      <p>
        <p>Once logged in, you'll find:</p>
          <li><strong>Live Music Stream</strong>: The main area where AI-generated music is played in real time.</li>
          <li><strong>Chat Box</strong>: Send messages to influence the AI's music generation.</li>
          <li><strong>Music Controls</strong>: Adjust parameters like tempo, mood, and instruments.</li>
      </p>
      <br />
    <!-- Interacting with the AI -->
    <h2 class="text-xl">Interacting with the AI</h2>
      <hr />
      <p>
        <li><strong>Use the Chat Box</strong>: Type a message describing the musical change you want (e.g., "Make the beat faster" or "Add a jazzy saxophone").</li>
        <li><strong>AI Interpretation</strong>: The system will analyze your request and translate it into a command for the music generation model.</li>
        <li><strong>Real-Time Adaptation</strong>: The AI will update the music based on your input.</li>
      </p>

      <h3>Tips for Effective Commands:</h3>
        <p><li>Be specific (e.g., "Slow down the tempo to a chill vibe").</li>
        <li>Use descriptive language related to music (e.g., "Add an upbeat bassline").</li>
        <li>Avoid spam messages—irrelevant inputs will be filtered out.</li></p>
      <br />
    <!-- Community Interaction & Moderation -->
    <h2 class="text-xl">Community Interaction & Moderation</h2>
      <hr />
      <p>
        <p>JukeBox is a collaborative space. To ensure a great experience for all users:</p>
          <li><strong>Respect Others</strong>: Avoid inappropriate or disruptive messages.</li>
          <li><strong>Stay on Topic</strong>: Messages unrelated to music may be ignored.</li>
          <li><strong>Automatic Moderation</strong>: Our AI filters and bans inappropriate behavior.</li>
      </p>
      <br />
    <!-- Troubleshooting & Support -->
    <h2 class="text-xl">Troubleshooting & Support</h2>
      <hr />
      <p>
        <p><strong>Having trouble?</strong></p>
          <li>Check the <strong>Help section</strong>.</li>
          <li>Visit our FAQ page for common issues.</li>
          <li>Contact support@jukebox.com for further assistance.</li>
        <p><strong>Now you're ready to dive into the world of AI-powered music creation! Enjoy JukeBox! </strong></p>
      </p>
      <br />
    </main>
  </div>
</template>
