<script setup lang="ts">
import { useRouter } from "vue-router"
import { ref } from "vue"

const router = useRouter()
const username = ref<string>("")
</script>

<template>
  <section
    class="relative min-h-screen overflow-hidden bg-gradient-to-b from-[#192231] to-[#253649] flex flex-col justify-center items-center"
  >
    <!-- Decorative elements -->
    <div class="absolute top-0 left-0 w-full h-full">
      <div
        class="absolute top-[10%] left-[5%] w-32 h-32 rounded-full bg-gradient-to-br from-[#4caf5020] to-[#8bc34a20] blur-xl"
      ></div>
      <div
        class="absolute bottom-[15%] right-[10%] w-40 h-40 rounded-full bg-gradient-to-tr from-[#865dff20] to-[#e384ff20] blur-xl"
      ></div>
      <div
        class="absolute top-[40%] right-[15%] w-24 h-24 rounded-full bg-gradient-to-r from-[#ffc10720] to-[#f4433620] blur-xl"
      ></div>
    </div>

    <!-- Main content -->
    <div
      class="relative z-10 text-center mx-auto px-8 py-10 backdrop-blur-sm bg-[#19223180] rounded-2xl border border-[#ffffff20] shadow-2xl max-w-md"
    >
      <div class="animate-pulse-slow mb-6">
        <a href="/">
          <img
            src="@/assets/hero.png"
            alt="Serpent Clash Logo"
            class="w-48 mx-auto drop-shadow-[0_0_15px_rgba(76,175,80,0.5)]"
          />
        </a>
      </div>

      <h1
        class="text-3xl mb-6 font-cherry text-transparent bg-clip-text bg-gradient-to-r from-[#4caf50] to-[#8bc34a]"
      >
        Enter Your Username
      </h1>

      <div class="flex flex-col space-y-6 items-center">
        <input
          type="text"
          class="w-full px-4 py-3 bg-[#ffffff10] border border-[#ffffff30] rounded-lg text-[#fff9c4] placeholder-[#ffffff80] focus:outline-none focus:ring-2 focus:ring-[#4caf5080] transition-all"
          v-model="username"
          placeholder="Enter your name..."
          @keydown.enter.prevent="router.push(`/play?username=${username}`)"
        />

        <div class="static-button-container w-40 h-14">
          <button
            class="text-white font-extrabold tracking-wide text-xl flex items-center justify-center h-full w-full"
            @click.prevent="router.push(`/play?username=${username}`)"
          >
            <span class="relative z-10 font-cherry">Play</span>
          </button>
        </div>
      </div>

      <div class="mt-8 text-[#fff9c4] text-sm">
        Join the ultimate multiplayer snake battle arena!
      </div>
    </div>

    <!-- Food decorative elements -->
    <div
      class="absolute top-[25%] left-[15%] w-6 h-6 rounded-full bg-[#ffc10730] blur-sm"
    ></div>
    <div
      class="absolute bottom-[30%] right-[20%] w-8 h-8 rounded-full bg-[#f4433630] blur-sm"
    ></div>
    <div
      class="absolute top-[60%] right-[30%] w-5 h-5 rounded-full bg-[#fff9c430] blur-sm"
    ></div>
  </section>
</template>

<style scoped>
.font-cherry {
  font-family: "Cherry Bomb One", system-ui;
  font-weight: 400;
  font-style: normal;
}

.static-button-container {
  position: relative;
  border-radius: 9999px;
  overflow: hidden;
  background: linear-gradient(to right, #4caf50, #8bc34a);
  box-shadow: 0 0 15px rgba(76, 175, 80, 0.5);
  transition: all 0.3s ease;
}

.static-button-container:hover {
  box-shadow: 0 0 20px rgba(76, 175, 80, 0.7);
  transform: scale(1.05) translateY(-2px);
}
</style>
