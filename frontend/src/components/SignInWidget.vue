<script lang="ts">
export default {
  emits: ['close'],
  methods: {
    closepopup() {
      // Émet l'événement `close` pour le composant parent
      this.$emit('close');
    },
  },
};
</script>

<template>
  <div class="widget-container">
    <div class="widget-sign-in">
      <button class="quit" @click="closepopup">
        <!--Icon quit-->
        <svg viewBox="7 2 10 24" fill=var(--color-bouton-color) xmlns="http://www.w3.org/2000/svg">
        <path d="M6.99486 7.00636C6.60433 7.39689 6.60433 8.03005 6.99486 8.42058L10.58 12.0057L6.99486 15.5909C6.60433 15.9814 6.60433 16.6146 6.99486 17.0051C7.38538 17.3956 8.01855 17.3956 8.40907 17.0051L11.9942 13.4199L15.5794 17.0051C15.9699 17.3956 16.6031 17.3956 16.9936 17.0051C17.3841 16.6146 17.3841 15.9814 16.9936 15.5909L13.4084 12.0057L16.9936 8.42059C17.3841 8.03007 17.3841 7.3969 16.9936 7.00638C16.603 6.61585 15.9699 6.61585 15.5794 7.00638L11.9942 10.5915L8.40907 7.00636C8.01855 6.61584 7.38538 6.61584 6.99486 7.00636Z"/>
        </svg>
      </button>
      <h1 class="title">Welcome on JukeBox !</h1>
      <h2 class="username-text">Username</h2>
      <input class="username-input" placeholder="Type your username here !"/>
      <h3 class="password-text">Password</h3>
      <input class="password-input" placeholder="Type your password here !"/>
      <h4 class="forget-pass"><u>Forget password ?</u></h4>
      <button class="connect-button">
        <!--Icon connect-->
        <svg viewBox="0 -7 38 38" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
          <g id="icons" stroke="none" stroke-width="1" fill="black" fill-rule="evenodd">
              <g id="ui-gambling-website-lined-icnos-casinoshunter" transform="translate(-1511.000000, -158.000000)" fill=var(--color-bouton-color) fill-rule="nonzero">
                  <g id="1" transform="translate(1350.000000, 120.000000)">
                      <path d="M187.812138,38.5802109 L198.325224,49.0042713 L198.41312,49.0858421 C198.764883,49.4346574 198.96954,49.8946897 199,50.4382227 L198.998248,50.6209428 C198.97273,51.0514917 198.80819,51.4628128 198.48394,51.8313977 L198.36126,51.9580208 L187.812138,62.4197891 C187.031988,63.1934036 185.770571,63.1934036 184.990421,62.4197891 C184.205605,61.6415481 184.205605,60.3762573 184.990358,59.5980789 L192.274264,52.3739093 L162.99947,52.3746291 C161.897068,52.3746291 161,51.4850764 161,50.3835318 C161,49.2819872 161.897068,48.3924345 162.999445,48.3924345 L192.039203,48.3917152 L184.990421,41.4019837 C184.205605,40.6237427 184.205605,39.3584519 184.990421,38.5802109 C185.770571,37.8065964 187.031988,37.8065964 187.812138,38.5802109 Z" id="right-arrow">
                      </path>
                  </g>
              </g>
          </g>
        </svg>
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
</style>

