<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, useTemplateRef } from "vue"
import { useRoute, useRouter } from "vue-router"
import DebugMenu from "@/components/DebugMenu.vue"
import ChatMenu from "@/components/ChatMenu.vue"
import { Game } from "@/classes/Game.ts"
import { DebugManager } from "@/classes/DebugManager.ts"

const route = useRoute()
const router = useRouter()

let username = (route.query.username as string) || ""
if (!username || username.length === 0) {
  router.push("/menu")
}

const gameCanvas = useTemplateRef<HTMLDivElement>("game-canvas")
const statsContainer = useTemplateRef<HTMLDivElement>("stats-container")
const status = ref("Connect")
const isGameReady = ref(false)
let game: Game
let debugManager: DebugManager
const debugMode: boolean = import.meta.env.VITE_DEBUG_MODE === "true"

if (debugMode) {
  console.log("Debug mode enabled")
}

function connectOrDisconnect(): void {
  if (
    game &&
    game.networkManager &&
    game.networkManager.socketState() === WebSocket.CLOSED
  ) {
    console.debug("Connecting socket connection...")
    game.connect()
  } else if (
    game &&
    game.networkManager &&
    game.networkManager.socketState() === WebSocket.OPEN
  ) {
    console.debug("Disconnecting socket connection...")
    game.disconnect()
  }
}

onMounted(() => {
  game = new Game(gameCanvas.value, statsContainer, status, username)
  debugManager = new DebugManager(game)
  game.init().then(() => {
    game.start()
    isGameReady.value = true
  })
})

onBeforeUnmount(() => {
  game.stop()
})
</script>

<template>
  <div class="w-full h-full">
    <div v-if="debugMode" ref="stats-container"></div>
    <div ref="game-canvas"></div>
    <div class="absolute top-2.5 right-2.5">
      <div class="text-end">
        <button
          class="w-32 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
          @click.prevent="connectOrDisconnect"
        >
          {{ status }}
        </button>
      </div>
      <DebugMenu
        v-if="debugMode && isGameReady"
        :debug-manager="debugManager"
      />
    </div>
    <div class="absolute bottom-1 left-1 space-y-2 text-sm">
      <ChatMenu />
    </div>
  </div>
</template>
