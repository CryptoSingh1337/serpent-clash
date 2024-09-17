<script setup lang="ts">
import { onBeforeUnmount, onMounted, useTemplateRef } from 'vue'
import { Player } from '../classes/entity'

let socket: WebSocket | null = null
let positions: Position[] = [{
  x: innerWidth / 2, y: innerHeight / 2
}]
let player: Player = new Player('1', '2', '#fff', positions)

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
        player.clientId = data.clientId
        player.sessionId = data.sessionId
        console.log("Player", player)
        break;
      }
      default: {
        console.log("invalid message type", data.type)
      }
    }
    setInterval(() => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        const payload = {
          clientId: player.clientId,
          sessionId: player.sessionId,
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
    player.draw(c)
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
