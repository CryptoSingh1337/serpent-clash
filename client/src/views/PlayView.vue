<script setup lang="ts">
import { onBeforeUnmount, onMounted, useTemplateRef } from 'vue'
import { Player } from '../classes/entity'
import type { BackendPlayer, Players, Position } from '../utils/types'

let socket: WebSocket | null = null
let positions: Position[] = [{
  x: innerWidth / 2, y: innerHeight / 2
}]
let currentPlayer: Player = new Player({ id: '1', color: '#fff', positions: positions, direction: 0 })
const players: Players = {}

const canvasRef = useTemplateRef<HTMLCanvasElement>('canvasRef')
onMounted(() => {
  socket = new WebSocket("ws://localhost:8080/ws")
  socket.onopen = () => {
    console.log("Socket opened")
  }
  socket.onclose = () => {
    console.log("Socket closed")
  }
  socket.onmessage = (data: any) => {
    data = JSON.parse(data.data)
    switch (data.type) {
      case "initialize": {
        currentPlayer.id = data.id
        console.log("Player", currentPlayer)
        break;
      }
      case "game_state": {
        const backendPlayers = data.body.playerStates as {[id: string]: BackendPlayer}
        for (const id in backendPlayers) {
          const backendPlayer = backendPlayers[id]
          if (!players[id]) {
            players[id] = new Player({id: id, color: '#fff', positions: backendPlayer.positions, direction: backendPlayer.direction})
          } else {
            if (currentPlayer.id === backendPlayer.id) {
              currentPlayer.positions = backendPlayer.positions
              currentPlayer.direction = backendPlayer.direction
            }
          }
        }
        for (const id in players) {
          if (!backendPlayers[id]) {
            delete players[id]
          }
        }
        break;
      }
      default: {
        console.log("invalid message type", data.type)
      }
    }
    setInterval(() => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        const payload = {
          id: currentPlayer.id,
          body: JSON.stringify({
            message: "Hello, server"
          })
        }
        socket.send(JSON.stringify(payload))
      }
    }, 5000)
  }
  socket.onerror = (err: any) => {
    console.error(err)
  }
  const _canvas = canvasRef.value
  if (!_canvas) {
    throw new Error("Can't find canvas element")
  }
  const c = _canvas.getContext('2d')
  if (!c) {
    throw new Error("Can't find canvas element")
  }

  _canvas.width = innerWidth
  _canvas.height = innerHeight

  let animationId
  function animate() {
    animationId = requestAnimationFrame(animate)
    if (!c) {
      throw new Error("Can't find canvas element")
    }
    c.fillStyle = 'rgba(0, 0, 0, 0.1)'
    c.fillRect(0, 0, innerWidth, innerHeight)
    for (const id in players) {
      const player = players[id]
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
  <canvas ref="canvasRef"></canvas>
</template>
