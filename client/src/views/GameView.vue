<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, useTemplateRef } from "vue"
import { useRoute, useRouter } from "vue-router"
import { GameDriver } from "@/drivers/game_driver.ts"
import { DebugDriver } from "@/drivers/debug_driver.ts"
import DebugMenu from "@/components/DebugMenu.vue"
import ChatMenu from "@/components/ChatMenu.vue"

const route = useRoute()
const router = useRouter()

let username = (route.query.username as string) || ""
if (!username || username.length === 0) {
  router.push("/menu")
}

const devicePixelRatio = window.devicePixelRatio || 1
const canvasRef = useTemplateRef<HTMLCanvasElement>("canvas-ref")
const statsContainer = useTemplateRef<HTMLDivElement>("stats-container")
let game: GameDriver | null = null
let debug: DebugDriver | null = null
const status = ref<string>("Connect")
const debugMode: boolean = import.meta.env.VITE_DEBUG_MODE === "true"

if (debugMode) {
  console.log("Debug mode enabled")
}

function connectOrDisconnect(): void {
  if (
    game &&
    game.socketDriver &&
    game.socketDriver.getReadyState() === WebSocket.CLOSED
  ) {
    console.debug("Connecting socket connection...")
    game.connect()
  } else if (
    game &&
    game.socketDriver &&
    game.socketDriver.getReadyState() === WebSocket.OPEN
  ) {
    console.debug("Disconnecting socket connection...")
    game.disconnect()
  }
}

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas) {
    throw new Error("Can't find canvas element")
  }
  const ctx = canvas.getContext("2d")
  if (!ctx) {
    throw new Error("Can't find canvas element")
  }
  canvas.width = window.innerWidth * devicePixelRatio
  canvas.height = window.innerHeight * devicePixelRatio
  window.addEventListener("resize", () => {
    canvas.width = innerWidth * devicePixelRatio
    canvas.height = innerHeight * devicePixelRatio
    if (game) {
      game.updateCameraWidthAndHeight(canvas.width, canvas.height)
    }
  })

  canvas.addEventListener("mousemove", (event) => {
    const rect = canvas.getBoundingClientRect()
    const scaleX = canvas.width / rect.width
    const scaleY = canvas.height / rect.height
    const screenX = (event.clientX - rect.left) * scaleX
    const screenY = (event.clientY - rect.top) * scaleY
    if (game) {
      game.updateMouseCoordinate(screenX, screenY)
    }
  })

  game = new GameDriver({ username, ctx, statsContainer, status })
  if (!game) {
    throw new Error("Cannot initialize game object")
  }
  debug = new DebugDriver(game)
  game.start()
})

onBeforeUnmount(() => {
  if (game) {
    game.disconnect()
  }
})
</script>

<template>
  <div class="w-full h-full">
    <div v-if="debugMode" ref="stats-container"></div>
    <canvas ref="canvas-ref"></canvas>
    <div class="absolute top-2.5 right-2.5">
      <div class="text-end">
        <button
          class="w-32 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
          @click.prevent="connectOrDisconnect"
        >
          {{ status }}
        </button>
      </div>
      <DebugMenu v-if="debugMode" :debug-menu="debug" />
    </div>
    <div class="absolute bottom-1 left-1 space-y-2 text-sm">
      <ChatMenu />
    </div>
  </div>
</template>
