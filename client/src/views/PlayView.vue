<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, useTemplateRef } from "vue"
import { Game } from "@/classes/game"

const devicePixelRatio = window.devicePixelRatio || 1
const canvasRef = useTemplateRef<HTMLCanvasElement>("canvas-ref")
let game: Game | null = null
const status = ref<string>("connect")

function connectOrDisconnect(): void {
  if (game && status.value === "connect") {
    game.connect()
    status.value = "disconnect"
  } else if (game && status.value === "disconnect") {
    game.disconnect()
    status.value = "connect"
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
      game.stats.updateCameraWidthAndHeight(innerWidth, innerHeight)
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

  game = new Game(c)
  if (!game) {
    throw new Error("Cannot initialize game object")
  }
  status.value = "disconnect"

  function animate() {
    if (!canvas) {
      throw new Error("Can't find canvas element")
    }
    if (!c) {
      throw new Error("Can't find canvas element")
    }
    if (!game) {
      throw new Error("game object is not initialized")
    }
    game.stats.pingCooldown -= 1
    if (game.stats.pingCooldown <= 0) {
      game.sendPingPayload()
    }
    c.clearRect(0, 0, canvas.width, canvas.height)
    game.update()
    game.render()
    requestAnimationFrame(animate)
  }
  game.calculateFps()
  animate()
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
    <button
      class="absolute top-5 right-5 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
      @click="connectOrDisconnect"
    >
      {{ status }}
    </button>
  </div>
</template>
