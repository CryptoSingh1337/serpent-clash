<script setup lang="ts"></script>

<template>
  <section id="hero" class="min-h-screen bg-cover bg-center flex flex-col justify-center items-center"
           style="background-image: url('https://example.com/hero-image.jpg');">
    <h1 class="text-6xl md:text-8xl font-extrabold neon-text text-teal-300">Serpent Clash</h1>
    <p class="text-xl mt-4 text-center text-gray-300">Outgrow. Outmaneuver. Dominate.</p>
    <button class="">
      Play Now
    </button>
  </section>
  <section id="description" class="py-16 bg-gray-800 text-center">
    <div class="max-w-4xl mx-auto px-4">
      <h2 class="text-4xl font-bold text-teal-300 mb-6">What is Serpent Clash?</h2>
      <p class="text-lg text-gray-300 mb-8">
        Serpent Clash is a fast-paced, real-time multiplayer snake game where players compete to outmaneuver and outgrow
        their opponents.
        With responsive gameplay using WebSocket communication and client-side prediction, you can slither your way to
        victory in an ever-changing battlefield.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">Real-time Multiplayer</h3>
          <p>Battle against players worldwide using real-time WebSocket communication for smooth, interactive gameplay.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">Client-side Prediction</h3>
          <p>Experience fluid movement even in high-latency situations, thanks to advanced prediction algorithms.</p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">Dynamic Snake Rendering</h3>
          <p>Snakes respond to real-time mouse input and update fluidly based on multiple coordinates.</p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">Authoritative Server</h3>
          <p>Our backend ensures fairness and consistent gameplay, providing a competitive environment for all players.
          </p>
        </div>
      </div>
    </div>
  </section>
  <section id="play" class="py-16 bg-gray-900 text-center">
    <div class="max-w-5xl mx-auto px-4">
      <h2 class="text-4xl font-bold text-teal-300 mb-6">Gameplay Preview</h2>
      <div class="mb-8">
        <img src="https://example.com/gameplay-preview.jpg" alt="Gameplay Preview" class="w-full rounded-lg shadow-lg">
      </div>
      <a href="#leaderboard"
         class="px-6 py-3 bg-purple-500 hover:bg-purple-400 text-gray-900 rounded-full glow-btn text-lg font-bold transition-all">
        Watch the Full Trailer
      </a>
    </div>
  </section>
  <section id="features" class="py-16 bg-gray-900 text-center">
    <div class="max-w-4xl mx-auto px-4">
      <h2 class="text-4xl font-bold text-teal-300 mb-6">Future Enhancements</h2>
      <p class="text-lg text-gray-300 mb-8">We’re always working on new features to improve gameplay:</p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">Server Reconciliation</h3>
          <p>Better synchronization to reduce server-side lag and improve the real-time experience.</p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">Player Movement Interpolation</h3>
          <p>Enhance the visibility of player actions with smoother animations and transitions.</p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">Expanded Arenas</h3>
          <p>Introducing larger and more dynamic battlefields for longer, more strategic battles.</p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg neon-border">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">New Game Mechanics</h3>
          <p>Look out for new elements like dynamic food generation and enhanced growth mechanics.</p>
        </div>
      </div>
    </div>
  </section>
  <footer class="py-8 bg-gray-800 text-center text-gray-500">
    <p>&copy; 2024 Serpent Clash. All Rights Reserved.</p>
    <a href="#hero" class="mt-4 inline-block text-teal-300 hover:text-teal-400">Back to Top</a>
  </footer>
</template>

<style scoped>
button {
  border: calc(0.08 * var(--m)) solid transparent;
  position: relative;
  color: #F3F3F3;
  font-size: var(--m);
  border-radius: calc(0.7 * var(--m));
  padding: calc(0.25 * var(--m)) calc(1 * var(--m));
  display: flex;
  justify-content: center;
  cursor: pointer;
  background: linear-gradient(var(--vt-c-black-soft), var(--vt-c-black-soft)), linear-gradient(var(--vt-c-black-soft) 50%, rgba(18, 18, 19, 0.6) 80%, rgba(18, 18, 19, 0)), linear-gradient(90deg, var(--orange), var(--yellow), var(--green), var(--blue), var(--purple), var(--pink), var(--red));
  background-origin: border-box;
  background-clip: padding-box, border-box, border-box;
  background-size: 200%;
  animation: animate 2s infinite linear;
}

button::before {
  content: '';
  background: linear-gradient(90deg, var(--orange), var(--yellow), var(--green), var(--blue), var(--purple), var(--pink), var(--red));
  height: 30%;
  width: 60%;
  position: absolute;
  bottom: -20%;
  z-index: -5;
  background-size: 200%;
  animation: animate 2s infinite linear;
  filter: blur(calc(0.8 * var(--m)));
}

button:hover, button:hover::before {
  animation: animate 0.5s infinite linear;
}

@keyframes animate {
  0% {background-position: 0}
  100% {background-position: 200%}
}
</style>
