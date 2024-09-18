<script setup lang="ts">
import { onBeforeUnmount, onMounted, useTemplateRef } from "vue"
import { Player } from "../classes/entity"
import type { BackendPlayer, Players, Position } from "../utils/types"

const devicePixelRatio = window.devicePixelRatio || 1
const canvasRef = useTemplateRef<HTMLCanvasElement>('canvas-ref')
let socket: WebSocket | null = null
let positions: Position[] = [{
  x: innerWidth / 2, y: innerHeight / 2
}]
let currentPlayer: Player = new Player({
  id: '1',
  color: '#fff',
  radius: 10,
  positions: positions,
  direction: 0
})
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

  socket = new WebSocket("ws://localhost:8080/ws")
  socket.onopen = () => {
    console.log("Socket opened")
  }
  socket.onclose = () => {
    console.log("Socket closed")
  }
  socket.onmessage = (data: any) => {
    data = JSON.parse(data.data)
    const body = data.body
    switch (data.type) {
      case "hello": {
        currentPlayer.id = body.id
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
              positions: backendPlayer.positions,
              direction: backendPlayer.direction})
          } else {
            if (currentPlayer.id === backendPlayer.id) {
              currentPlayer.positions = backendPlayer.positions
              currentPlayer.direction = backendPlayer.direction
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

  let animationId
  function animate() {
    animationId = requestAnimationFrame(animate)
    if (!c) {
      throw new Error("Can't find canvas element")
    }
    c.fillStyle = 'rgba(0, 0, 0, 0.1)'
    c.fillRect(0, 0, innerWidth, innerHeight)
    for (const id in frontendPlayers) {
      const player = frontendPlayers[id]
      player.draw(c)
    }
  }
  animate()
})

onBeforeUnmount(() => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close()
  }
})
</script>

<template>
  <h1>Game View</h1>
  <canvas ref="canvas-ref"></canvas>
</template>

<style>
canvas {
  width: 100%;
  height: 100%;
}
</style>
