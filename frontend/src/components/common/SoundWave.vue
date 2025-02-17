<script>
export default {
  name: 'BarVisualizer',
  data() {
    return {
      audioCtx: null,
      analyser: null,
      dataArray: null,
      bufferLength: null,
      animationFrameId: null,
      numberOfBars: 50 // Nombre de rectangles à afficher
    }
  },
  mounted() {
    this.initAudio()
  },
  methods: {
    initAudio() {
      // Création du contexte audio
      const AudioContext = window.AudioContext || window.webkitAudioContext
      this.audioCtx = new AudioContext()

      // Récupération de l'élément audio et création de la source audio
      const audioElement = this.$refs.audio
      const audioSrc = this.audioCtx.createMediaElementSource(audioElement)

      // Création d'un analyser pour extraire les données audio
      // On choisit un fftSize plus petit pour obtenir moins de données et un affichage plus "pixélisé"
      this.analyser = this.audioCtx.createAnalyser()
      this.analyser.fftSize = 2048
      audioSrc.connect(this.analyser)
      this.analyser.connect(this.audioCtx.destination)

      // La propriété frequencyBinCount correspond à la moitié du fftSize
      this.bufferLength = this.analyser.frequencyBinCount
      this.dataArray = new Uint8Array(this.bufferLength)

      // Ajout d'écouteurs d'événements sur l'élément audio
      audioElement.addEventListener('play', this.startAnimation)
      audioElement.addEventListener('pause', this.stopAnimation)
      audioElement.addEventListener('ended', this.stopAnimation)
    },
    playAudio() {
      const audioElement = this.$refs.audio
      // Reprendre le contexte audio s'il est suspendu
      if (this.audioCtx.state === 'suspended') {
        this.audioCtx.resume()
      }
      audioElement.play()
    },
    pauseAudio() {
      this.$refs.audio.pause()
    },
    endAudio() {
      const audioElement = this.$refs.audio
      audioElement.pause()
      audioElement.currentTime = 0
      this.stopAnimation()
      // Effacer le canvas
      const canvas = this.$refs.canvas
      const canvasCtx = canvas.getContext('2d')
      canvasCtx.clearRect(0, 0, canvas.width, canvas.height)
    },
    startAnimation() {
      if (!this.animationFrameId) {
        this.draw()
      }
    },
    stopAnimation() {
      if (this.animationFrameId) {
        cancelAnimationFrame(this.animationFrameId)
        this.animationFrameId = null
      }
    },
    draw() {
      const canvas = this.$refs.canvas
      const canvasCtx = canvas.getContext('2d')

      // Si la lecture est en pause, on arrête la boucle d'animation
      if (this.$refs.audio.paused) {
        this.animationFrameId = null
        return
      }

      // Planifier le prochain dessin
      this.animationFrameId = requestAnimationFrame(this.draw)

      // Récupérer les données de fréquence
      this.analyser.getByteFrequencyData(this.dataArray)

      // Couleur
      const rootStyles = getComputedStyle(document.documentElement)
      const background = rootStyles.getPropertyValue('--color-background').trim()

      // Effacer le canvas
      canvasCtx.fillStyle = background
      canvasCtx.fillRect(0, 0, canvas.width, canvas.height)

      // Calcul du nombre de rectangles et de leur largeur
      const barWidth = canvas.width / this.numberOfBars
      const step = Math.floor(this.bufferLength / this.numberOfBars)

      // Parcourir la moitié du nombre de barres à afficher
      for (let i = 0; i < this.numberOfBars / 2; i++) {
        // Calculer la valeur moyenne sur le segment de données correspondant
        let sum = 0

        for (let j = 0; j < step; j++) {
          sum += this.dataArray[i * step + j]
        }

        const avg = sum / step

        // Gap entre le haut et le bas
        const gap = 1

        // Calculer la hauteur du rectangle (les valeurs varient de 0 à 255)
        const barHeight = ((avg / 255) * canvas.height) / 2 - gap

        //Couleur
        const TopColor = rootStyles.getPropertyValue('--color-sound-top').trim()
        const BottomColor = rootStyles.getPropertyValue('--color-sound-bottom').trim()

        // Coin haut droit
        let x = i * barWidth + canvas.width / 2
        let y = canvas.height / 2 - gap - barHeight
        canvasCtx.fillStyle = TopColor
        canvasCtx.fillRect(x, y, barWidth - 2, barHeight)

        // Coin bas droit
        x = i * barWidth + canvas.width / 2
        y = canvas.height / 2 + gap
        canvasCtx.fillStyle = BottomColor
        canvasCtx.fillRect(x, y, barWidth - 2, barHeight)

        // Coin haut gauche
        x = canvas.width / 2 - i * barWidth
        y = canvas.height / 2 - gap - barHeight
        // Si i != 0
        if (i != 0) {
          canvasCtx.fillStyle = TopColor
          canvasCtx.fillRect(x, y, barWidth - 2, barHeight)
        }

        // Coin bas gauche
        x = canvas.width / 2 - i * barWidth
        y = canvas.height / 2 + gap
        if (i != 0) {
          canvasCtx.fillStyle = BottomColor
          canvasCtx.fillRect(x, y, barWidth - 2, barHeight)
        }
      }
    }
  }
}
</script>

<template>
  <div>
    <!-- Canvas pour dessiner les rectangles -->
    <canvas ref="canvas" width="800" height="200"></canvas>

    <!-- Élément audio sans contrôles intégrés -->
    <audio ref="audio">
      <source src="./test.mp3" type="audio/mp3" />
      Votre navigateur ne supporte pas l'élément audio.
    </audio>

    <!-- Boutons de contrôle personnalisés -->
    <div class="controls">
      <button @click="playAudio">Play</button>
      <button @click="pauseAudio">Pause</button>
      <button @click="endAudio">End</button>
    </div>
  </div>
</template>

<style scoped>

.controls {
  margin-top: 10px;
}
.controls button {
  margin-right: 10px;
  padding: 5px 10px;
  cursor: pointer;
}
</style>
