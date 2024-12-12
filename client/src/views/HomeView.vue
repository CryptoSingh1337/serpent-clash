<script setup lang="ts">
import { useRouter } from "vue-router"

const router = useRouter()
</script>

<template>
  <section
    id="hero"
    class="min-h-screen bg-cover bg-center flex flex-col justify-center items-center grid-bg ba-grid anim"
  >
    <div class="inner"></div>
    <div class="z-50 text-center space-y-8">
      <h1 class="text-6xl md:text-8xl font-extrabold text-accent-1">
        Serpent Clash
      </h1>
      <p class="text-xl mt-4 text-gray-300">Outgrow. Outmaneuver. Dominate.</p>
      <div class="flex justify-around">
        <div class="button-container w-32 h-20">
          <button
            class="text-white font-bold text-2xl button-content"
            @click.prevent="router.push('/menu')"
          >
            Play
          </button>
        </div>
      </div>
    </div>
  </section>
  <section id="description" class="py-16 bg-gray-800 text-center">
    <div class="max-w-4xl mx-auto px-4">
      <h2 class="text-4xl font-bold text-teal-300 mb-6">
        What is Serpent Clash?
      </h2>
      <p class="text-lg text-gray-300 mb-8">
        Serpent Clash is a fast-paced, real-time multiplayer snake game where
        players compete to outmaneuver and outgrow their opponents. With
        responsive gameplay using WebSocket communication and client-side
        prediction, you can slither your way to victory in an ever-changing
        battlefield.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">
            Real-time Multiplayer
          </h3>
          <p>
            Battle against players worldwide using real-time WebSocket
            communication for smooth, interactive gameplay.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">
            Client-side Prediction
          </h3>
          <p>
            Experience fluid movement even in high-latency situations, thanks to
            advanced prediction algorithms.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">
            Dynamic Snake Rendering
          </h3>
          <p>
            Snakes respond to real-time mouse input and update fluidly based on
            multiple coordinates.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-green-400 mb-4">
            Authoritative Server
          </h3>
          <p>
            Our backend ensures fairness and consistent gameplay, providing a
            competitive environment for all players.
          </p>
        </div>
      </div>
    </div>
  </section>
  <section id="play" class="py-16 bg-gray-900 text-center">
    <div class="max-w-5xl mx-auto px-4">
      <h2 class="text-4xl font-bold text-teal-300 mb-6">Gameplay Preview</h2>
      <div class="mb-8">
        <iframe
          class="mx-auto my-10 aspect-video w-full"
          src="https://www.youtube-nocookie.com/embed/LiovUkBWOos"
          title="YouTube video player"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
          :allowfullscreen="false"
        ></iframe>
      </div>
      <a
        href="#leaderboard"
        class="px-6 py-3 bg-purple-500 hover:bg-purple-400 text-gray-900 rounded-full glow-btn text-lg font-bold transition-all"
      >
        Watch the Full Trailer
      </a>
    </div>
  </section>
  <section id="features" class="py-16 bg-gray-900 text-center">
    <div class="max-w-4xl mx-auto px-4">
      <h2 class="text-4xl font-bold text-teal-300 mb-6">Future Enhancements</h2>
      <p class="text-lg text-gray-300 mb-8">
        We’re always working on new features to improve gameplay:
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">
            Server Reconciliation
          </h3>
          <p>
            Better synchronization to reduce server-side lag and improve the
            real-time experience.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">
            Player Movement Interpolation
          </h3>
          <p>
            Enhance the visibility of player actions with smoother animations
            and transitions.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">
            Expanded Arenas
          </h3>
          <p>
            Introducing larger and more dynamic battlefields for longer, more
            strategic battles.
          </p>
        </div>
        <div class="bg-gray-900 p-6 rounded-lg">
          <h3 class="text-2xl font-semibold text-purple-400 mb-4">
            New Game Mechanics
          </h3>
          <p>
            Look out for new elements like dynamic food generation and enhanced
            growth mechanics.
          </p>
        </div>
      </div>
    </div>
  </section>
  <footer class="py-8 bg-gray-800 text-center text-gray-500">
    <p>&copy; 2024 Serpent Clash. All Rights Reserved.</p>
    <a href="#hero" class="mt-4 inline-block text-teal-300 hover:text-teal-400"
      >Back to Top</a
    >
  </footer>
</template>
