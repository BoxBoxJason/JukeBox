<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue'
import PlayIcon from '../icons/PlayIcon.vue'
import StopIcon from '../icons/StopIcon.vue'

export default defineComponent({
  name: 'SoundWave',
  components: { PlayIcon, StopIcon },
  setup() {
    const audioCtx = ref<AudioContext | null>(null)
    const analyser = ref<AnalyserNode | null>(null)
    const dataArray = ref<Uint8Array | null>(null)
    const bufferLength = ref<number>(0)
    const animationFrameId = ref<number | null>(null)
    const numberOfBars = 151

    const audioRef = ref<HTMLAudioElement | null>(null)
    const canvasRef = ref<HTMLCanvasElement | null>(null)

    const initAudio = () => {
      const AudioContextConstructor = window.AudioContext || (window as any).webkitAudioContext
      audioCtx.value = new AudioContextConstructor()

      if (!audioRef.value) {
        console.error('Élément audio non trouvé.')
        return
      }
      const audioElement = audioRef.value
      const audioSrc = audioCtx.value.createMediaElementSource(audioElement)

      analyser.value = audioCtx.value.createAnalyser()
      analyser.value.fftSize = 2048
      audioSrc.connect(analyser.value)
      analyser.value.connect(audioCtx.value.destination)

      bufferLength.value = analyser.value.frequencyBinCount
      dataArray.value = new Uint8Array(bufferLength.value)

      audioElement.addEventListener('play', startAnimation)
      audioElement.addEventListener('ended', stopAnimation)
    }

    const initCanvas = () => {
      if (!canvasRef.value) return

      const canvas = canvasRef.value
      const canvasCtx = canvas.getContext('2d')

      if (!canvasCtx) return

      // Couleur du fond
      const rootStyles = getComputedStyle(document.documentElement)
      const background = rootStyles.getPropertyValue('--color-background').trim()

      canvasCtx.fillStyle = background
      canvasCtx.fillRect(0, 0, canvas.width, canvas.height)

      // Largeur des rectangles
      const barWidth = canvas.width / numberOfBars
      const step = Math.floor(bufferLength.value / numberOfBars)

      // Remplissage des rectangles
      for (let i = 0; i < numberOfBars; i++) {
        // Gap entre le heut et le bas
        const gap = 0

        // Hauteur d'une barre sur une plage
        const barHeight = 1

        // Couleurs des rectangles
        const topColor = rootStyles.getPropertyValue('--color-sound-top').trim()
        const bottomColor = rootStyles.getPropertyValue('--color-sound-bottom').trim()

        // Haut
        let x = i * barWidth
        let y = canvas.height / 2 - gap - barHeight
        canvasCtx.fillStyle = topColor
        canvasCtx.fillRect(x, y, barWidth - 2, barHeight)

        // Bas
        x = i * barWidth
        y = canvas.height / 2 + gap
        canvasCtx.fillStyle = bottomColor
        canvasCtx.fillRect(x, y, barWidth - 2, barHeight)
      }
    }

    const playAudio = async () => {
      if (!audioRef.value) return

      if (!audioCtx.value) {
        initAudio()
      }
      if (audioCtx.value && audioCtx.value.state === 'suspended') {
        await audioCtx.value.resume()
      }
      audioRef.value.play()
    }

    const endAudio = () => {
      if (!audioRef.value || !canvasRef.value) return
      audioRef.value.pause()
      audioRef.value.currentTime = 0
      stopAnimation()
      const canvas = canvasRef.value
      const canvasCtx = canvas.getContext('2d')
      if (canvasCtx) {
        canvasCtx.clearRect(0, 0, canvas.width, canvas.height)
      }
      initCanvas()
    }

    const startAnimation = () => {
      if (!animationFrameId.value) {
        draw()
      }
    }

    const stopAnimation = () => {
      if (animationFrameId.value !== null) {
        cancelAnimationFrame(animationFrameId.value)
        animationFrameId.value = null
      }
    }

    const draw = () => {
      if (
        !canvasRef.value ||
        !audioRef.value ||
        !analyser.value ||
        !dataArray.value ||
        !bufferLength.value
      )
        return

      const canvas = canvasRef.value
      const canvasCtx = canvas.getContext('2d')
      if (!canvasCtx) return

      if (audioRef.value.paused) {
        animationFrameId.value = null
        return
      }

      animationFrameId.value = requestAnimationFrame(draw)
      analyser.value.getByteFrequencyData(dataArray.value)

      const rootStyles = getComputedStyle(document.documentElement)
      const background = rootStyles.getPropertyValue('--color-background').trim()

      canvasCtx.fillStyle = background
      canvasCtx.fillRect(0, 0, canvas.width, canvas.height)

      // Largeur des rectangles
      const barWidth = canvas.width / numberOfBars
      const step = Math.floor(bufferLength.value / numberOfBars)

      // Remplissage des rectangles
      for (let i = 0; i < numberOfBars / 2; i++) {
        let sum = 0
        for (let j = 0; j < step; j++) {
          sum += dataArray.value[i * step + j]
        }

        // Moyenne des valeurs sur une plage
        const avg = sum / step

        // Gap entre le heut et le bas
        const gap = 1

        // Hauteur d'une barre sur une plage
        let barHeight = ((avg / 255) * canvas.height) / 2

        // Correction de la hauteur de la barre
        const pourcent = 0.07
        if (i > numberOfBars / 4) {
          barHeight = barHeight * (1 - pourcent) ** (i - numberOfBars / 4)
        }

        // Couleurs des rectangles
        const topColor = rootStyles.getPropertyValue('--color-sound-top').trim()
        const bottomColor = rootStyles.getPropertyValue('--color-sound-bottom').trim()

        // Coin haut droit
        let x = i * barWidth + canvas.width / 2 - barWidth / 2
        let y = canvas.height / 2 - gap - barHeight
        canvasCtx.fillStyle = topColor
        canvasCtx.fillRect(x, y, barWidth - 2, barHeight)

        // Coin bas droit
        x = i * barWidth + canvas.width / 2 - barWidth / 2
        y = canvas.height / 2 + gap
        canvasCtx.fillStyle = bottomColor
        canvasCtx.fillRect(x, y, barWidth - 2, barHeight)

        // Coin haut gauche
        x = canvas.width / 2 - i * barWidth - barWidth / 2
        y = canvas.height / 2 - gap - barHeight
        if (i !== 0) {
          canvasCtx.fillStyle = topColor
          canvasCtx.fillRect(x, y, barWidth - 2, barHeight)
        }

        // Coin bas gauche
        x = canvas.width / 2 - i * barWidth - barWidth / 2
        y = canvas.height / 2 + gap
        if (i !== 0) {
          canvasCtx.fillStyle = bottomColor
          canvasCtx.fillRect(x, y, barWidth - 2, barHeight)
        }
      }
    }

    onMounted(() => {
      initAudio()
      initCanvas()
    })

    return {
      audioRef,
      canvasRef,
      playAudio,
      endAudio,
      startAnimation,
      stopAnimation,
      draw
    }
  }
})
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-4">
    <!-- Canvas pour dessiner les rectangles -->
    <canvas ref="canvasRef" width="1000" height="200"></canvas>

    <!-- Élément audio sans contrôles intégrés -->
    <audio ref="audioRef">
      <source src="./test.mp3" type="audio/mp3" />
      Votre navigateur ne supporte pas l'élément audio.
    </audio>

    <!-- Boutons de contrôle personnalisés -->
    <div class="flex flex-row gap-8 w-full justify-center mr-2">
      <button
        @click="playAudio"
        class="rounded-full bg-[var(--color-button-color)] border border-[var(--color-border)] p-2 hover:bg-[var(--color-background)]"
      >
        <PlayIcon />
      </button>
      <button @click="endAudio" class="rounded-full bg-[var(--color-button-color)] border border-[var(--color-border)] p-2 hover:bg-[var(--color-background)]">
        <StopIcon />
      </button>
    </div>
  </div>
</template>
