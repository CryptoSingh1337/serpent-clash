<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, useTemplateRef } from "vue"
import { GameDriver } from "@/drivers/game_driver.ts"
import { DebugDriver } from "@/drivers/debug_driver.ts"
import DebugMenu from "@/components/DebugMenu.vue"

const devicePixelRatio = window.devicePixelRatio || 1
const canvasRef = useTemplateRef<HTMLCanvasElement>("canvas-ref")
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
  const c = canvas.getContext("2d")
  if (!c) {
    throw new Error("Can't find canvas element")
  }
  canvas.width = innerWidth * devicePixelRatio
  canvas.height = innerHeight * devicePixelRatio
  window.addEventListener("resize", () => {
    canvas.width = innerWidth
    canvas.height = innerHeight
    if (game) {
      game.updateCameraWidthAndHeight(innerWidth, innerHeight)
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

  game = new GameDriver(c, status)
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
  <div class="h-full w-full">
    <canvas ref="canvas-ref"></canvas>
    <div class="absolute top-5 right-5">
      <div class="relative text-end">
        <button
          class="w-32 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
          @click.prevent="connectOrDisconnect"
        >
          {{ status }}
        </button>
      </div>
      <DebugMenu v-if="debugMode" :debug-menu="debug" />
    </div>
  </div>
</template>
