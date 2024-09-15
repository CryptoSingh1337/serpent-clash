<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
let socket: WebSocket | null = null
let player: Player = {
  clientId: "",
  sessionId: "",
}
const messages = ref<Object[]>([])


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
    const body = JSON.parse(data.body)
    console.log("Body", body)
    switch (data.eventType) {
      case "initialize": {
        player.clientId = body.clientId
        player.sessionId = body.sessionId
      }
    }
    messages.value.push(data)
  }
  socket.onerror = (err: any) => {
    console.error(err)
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
      messages.value.push(payload)
    }
  }, 5000)
})

onBeforeUnmount(() => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close()
  }
})
</script>

<template>
  <h1>Game View</h1>
  <div>
    <ul>
      <li v-for="message in messages">{{ message }}</li>
    </ul>
  </div>
</template>
