<script setup lang="ts">
import { onBeforeUnmount, onMounted, useTemplateRef } from "vue"
import { Player } from "../classes/entity"
import type { BackendPlayer, Players } from "../utils/types"
import { Stats } from "../classes/stats"
import { Constants } from '../utils/constants'

const devicePixelRatio = window.devicePixelRatio || 1
const canvasRef = useTemplateRef<HTMLCanvasElement>('canvas-ref')
let socket: WebSocket | null = null
let stats = new Stats()
let currentPlayerId: string = ""
const mouseCoordinate = { x: innerWidth / 2, y: innerHeight / 2 }
const frontendPlayers: Players = {}

onMounted(() => {
  // canvas setup
  const canvas = canvasRef.value
  if (!canvas) {
    throw new Error("Can't find canvas element")
  }
  const c = canvas.getContext('2d')
  if (!c) {
    throw new Error("Can't find canvas element")
  }
  canvas.width = innerWidth * devicePixelRatio
  canvas.height = innerHeight * devicePixelRatio
  window.addEventListener('resize', () => {
    canvas.width = innerWidth
    canvas.height = innerHeight
  })

  canvas.addEventListener("mousemove", (event) => {
    mouseCoordinate.x = event.clientX - canvas.offsetLeft
    mouseCoordinate.y = event.clientY - canvas.offsetTop
    stats.updateMouseCoordinate(mouseCoordinate.x, mouseCoordinate.y)
  })

  // socket setup
  socket = new WebSocket("ws://localhost:8080/ws")
  socket.onopen = () => {
    console.log("Socket opened")
    stats.updateStatus("online")
  }
  socket.onclose = () => {
    console.log("Socket closed")
    stats.updateStatus("offline")
    stats.updateHeadCoordinate(0, 0)
    for (const id in frontendPlayers) {
      delete frontendPlayers[id]
    }
  }
  socket.onmessage = (data: any) => {
    data = JSON.parse(data.data)
    const body = data.body
    switch (data.type) {
      case "hello": {
        currentPlayerId = body.id
        stats.updatePlayerId(currentPlayerId)
        break;
      }
      case "game_state": {
        const backendPlayers = body.playerStates as {[id: string]: BackendPlayer}
        for (const id in backendPlayers) {
          const backendPlayer = backendPlayers[id]
          if (!frontendPlayers[id]) {
            frontendPlayers[id] = new Player({
              id: id,
              color: backendPlayer.color,
              radius: 10,
              positions: backendPlayer.positions
            })
          } else {
            const frontendPlayer = frontendPlayers[id]
            frontendPlayer.positions = backendPlayer.positions
            if (currentPlayerId === id) {
              stats.updateHeadCoordinate(frontendPlayer.positions[0].x,
                frontendPlayer.positions[0].y)
            }
          }
        }
        for (const id in frontendPlayers) {
          if (!backendPlayers[id]) {
            delete frontendPlayers[id]
          }
        }
        break;
      }
      default: {
        console.log("invalid message type", data.type)
      }
    }
  }
  socket.onerror = (err: any) => {
    console.error(err)
  }

  setInterval(() => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({
        type: "movement",
        body: {
          coordinate: mouseCoordinate
        }
      }))
    }
  }, 1000 / Constants.tickRate)

  function renderPlayers(): void {
    for (const id in frontendPlayers) {
      const player = frontendPlayers[id]
      player.draw(c)
    }
  }

  function animate() {
    if (!c) {
      throw new Error("Can't find canvas element")
    }
    if (!canvas) {
      throw new Error("Can't find canvas element")
    }
    c.clearRect(0, 0, canvas.width, canvas.height)
    for (const id in frontendPlayers) {
      const player = frontendPlayers[id]
      player.draw(c)
    }
    stats.renderStats(c)
    renderPlayers()
    requestAnimationFrame(animate)
  }
  stats.calculateFps()
  animate()
})

onBeforeUnmount(() => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close()
  }
})
</script>

<template>
  <div class="h-full w-full">
    <canvas ref="canvas-ref"></canvas>
  </div>
</template>
