import { Stats } from "@/classes/stats"
import { Constants } from "@/utils/constants"
import type { BackendPlayer, Players } from "@/utils/types"
import { Player } from "@/classes/entity"

export class Game {
  ctx: CanvasRenderingContext2D
  socket: WebSocket | null = null
  stats: Stats
  playerId: string = ""
  mouseCoordinate: { x: number; y: number }
  frontendPlayers: Players = {}
  currentPlayer: Player | null = null

  constructor(ctx: CanvasRenderingContext2D) {
    if (!ctx) {
      throw new Error("Can't find canvas element")
    }
    this.ctx = ctx
    this.mouseCoordinate = { x: 0, y: 0 }
    this.stats = new Stats()
    this.initSocket()
    setInterval(() => {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(
          JSON.stringify({
            type: "movement",
            body: {
              coordinate: this.mouseCoordinate
            }
          })
        )
      }
    }, 1000 / Constants.tickRate)
  }

  initSocket(): void {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    const url = window.location.hostname.includes("localhost")
      ? `${protocol}//${window.location.hostname}:${Constants.serverPort}/ws`
      : `${protocol}//${window.location.hostname}/ws`
    this.socket = new WebSocket(url)
    this.socket.onopen = () => {
      console.log("Socket opened")
      this.stats.updateStatus("online")
    }
    this.socket.onclose = () => {
      console.log("Socket closed")
      this.stats.reset()
      this.currentPlayer = null
      for (const id in this.frontendPlayers) {
        delete this.frontendPlayers[id]
      }
    }
    this.socket.onerror = (err: any) => {
      console.error(err)
    }
    this.socket.onmessage = (data: any) => {
      data = JSON.parse(data.data)
      const body = data.body
      switch (data.type) {
        case "hello": {
          this.playerId = body.id
          this.stats.updatePlayerId(this.playerId)
          break
        }
        case "pong": {
          this.stats.updatePing(performance.now() - body.timestamp)
          break
        }
        case "game_state": {
          const backendPlayers = body.playerStates as {
            [id: string]: BackendPlayer
          }
          for (const id in backendPlayers) {
            const backendPlayer = backendPlayers[id]
            if (!this.frontendPlayers[id]) {
              this.frontendPlayers[id] = new Player({
                id: id,
                color: backendPlayer.color,
                radius: 10,
                positions: backendPlayer.positions
              })
              if (this.playerId === id) {
                this.currentPlayer = this.frontendPlayers[id]
              }
            } else {
              const frontendPlayer = this.frontendPlayers[id]
              frontendPlayer.positions = backendPlayer.positions
              if (this.playerId === id) {
                this.stats.updateHeadCoordinate(
                  frontendPlayer.positions[0].x,
                  frontendPlayer.positions[0].y
                )
              }
            }
          }
          for (const id in this.frontendPlayers) {
            if (!backendPlayers[id]) {
              delete this.frontendPlayers[id]
            }
          }
          break
        }
        default: {
          console.log("invalid message type", data.type)
        }
      }
    }
  }

  sendPingPayload(): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(
        JSON.stringify({
          type: "ping",
          body: {
            timestamp: performance.now()
          }
        })
      )
      this.stats.pingCooldown = Constants.pingCooldown
    }
  }

  connect(): void {
    if (!this.socket || this.socket.readyState === WebSocket.CLOSED) {
      this.initSocket()
    }
  }

  renderPlayers(): void {
    for (const id in this.frontendPlayers) {
      const player = this.frontendPlayers[id]
      player.draw(this.ctx)
    }
  }

  renderStats(): void {
    this.stats.renderStats(this.ctx)
  }

  calculateFps(): void {
    this.stats.calculateFps()
  }

  updateMouseCoordinate(x: number, y: number): void {
    this.mouseCoordinate.x = x
    this.mouseCoordinate.y = y
    this.stats.updateMouseCoordinate(x, y)
  }

  disconnect(): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.close()
    }
  }
}
