<script setup lang="ts">
import { useRouter } from "vue-router"
import { ref } from "vue"

const router = useRouter()
const username = ref<string>("")
</script>

<template>
  <div class="w-screen h-screen flex items-center justify-center">
    <div class="flex flex-col space-y-4 items-center">
      <div>Username</div>
      <input
        type="text"
        class=""
        v-model="username"
        placeholder="Enter your name..."
        @keydown.enter.prevent="router.push(`/play?username=${username}`)"
      />
      <button
        class="w-16 h-8 rounded-full bg-indigo-400"
        @click.prevent="router.push(`/play?username=${username}`)"
      >
        Play
      </button>
    </div>
  </div>
</template>
