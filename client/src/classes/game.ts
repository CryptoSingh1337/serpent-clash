import { Stats } from "@/classes/stats"
import { Constants } from "@/utils/constants"
import type { BackendPlayer, Players } from "@/utils/types"
import { Player } from "@/classes/entity"
import { clamp } from "@/utils/helper"
import { Camera } from "@/classes/camera"

class HexGrid {
  hexSize: number
  hexHeight: number
  hexWidth: number
  verticalSpacing: number
  horizontalSpacing: number

  constructor(hexSize: number) {
    this.hexSize = hexSize
    this.hexHeight = hexSize * 2
    this.hexWidth = Math.sqrt(25) * hexSize
    this.verticalSpacing = (this.hexHeight * 8) / 4
    this.horizontalSpacing = this.hexWidth / 1.5
  }

  drawHex(ctx: CanvasRenderingContext2D, x: number, y: number) {
    ctx.beginPath()
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i
      const hx = x + this.hexSize * Math.cos(angle)
      const hy = y + this.hexSize * Math.sin(angle)
      if (i === 0) {
        ctx.moveTo(hx, hy)
      } else {
        ctx.lineTo(hx, hy)
      }
    }
    ctx.closePath()
    ctx.stroke()
  }

  render(ctx: CanvasRenderingContext2D, camera: Camera) {
    const startX = Math.max(
      Constants.worldBoundary.minX,
      Math.floor(camera.x / this.horizontalSpacing) * this.horizontalSpacing -
        this.horizontalSpacing
    )
    const startY = Math.max(
      Constants.worldBoundary.minY,
      Math.floor(camera.y / this.verticalSpacing) * this.verticalSpacing -
        this.verticalSpacing
    )
    const endX = Math.min(
      Constants.worldBoundary.maxX,
      camera.x + camera.width + this.horizontalSpacing
    )
    const endY = Math.min(
      Constants.worldBoundary.maxY,
      camera.y + camera.height + this.verticalSpacing
    )

    ctx.strokeStyle = "rgba(200, 200, 200, 0.5)"
    ctx.lineWidth = 1

    for (let y = startY; y < endY; y += this.verticalSpacing) {
      for (let x = startX; x < endX; x += this.horizontalSpacing) {
        const screenPos = camera.worldToScreen(x, y)
        this.drawHex(ctx, screenPos.x, screenPos.y)
        if (y + this.verticalSpacing / 2 < Constants.worldBoundary.maxY) {
          const offsetScreenPos = camera.worldToScreen(
            x + this.horizontalSpacing / 2,
            y + this.verticalSpacing / 2
          )
          this.drawHex(ctx, offsetScreenPos.x, offsetScreenPos.y)
        }
      }
    }
  }
}

export class Game {
  ctx: CanvasRenderingContext2D
  socket: WebSocket | null = null
  stats: Stats
  playerId: string = ""
  mouseCoordinate: { x: number; y: number }
  camera: Camera
  hexGrid: HexGrid
  frontendPlayers: Players = {}
  currentPlayer: Player | null = null
  inputs: { seq: number; x: number; y: number }[] = []
  seq: number = 0

  constructor(ctx: CanvasRenderingContext2D) {
    if (!ctx) {
      throw new Error("Can't find canvas element")
    }
    this.ctx = ctx
    this.mouseCoordinate = { x: 0, y: 0 }
    this.stats = new Stats()
    this.stats.updateCameraWidthAndHeight(ctx.canvas.width, ctx.canvas.height)
    this.camera = new Camera(ctx.canvas.width, ctx.canvas.height)
    this.hexGrid = new HexGrid(50)
    this.initSocket()
    this.setupMouseTracking()
  }

  setupMouseTracking(): void {
    setInterval(
      () => {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
          const worldCoordinate = this.camera.screenToWorld(
            this.mouseCoordinate.x,
            this.mouseCoordinate.y
          )
          worldCoordinate.x = clamp(
            worldCoordinate.x,
            Constants.worldBoundary.minX,
            Constants.worldBoundary.maxX
          )
          worldCoordinate.y = clamp(
            worldCoordinate.y,
            Constants.worldBoundary.minY,
            Constants.worldBoundary.maxY
          )
          this.inputs.push({
            seq: ++this.seq,
            x: this.mouseCoordinate.x,
            y: this.mouseCoordinate.y
          })
          this.socket.send(
            JSON.stringify({
              type: "movement",
              body: {
                seq: this.seq,
                coordinate: worldCoordinate
              }
            })
          )
        }
      },
      Math.floor(1000 / Constants.tickRate)
    )
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
              frontendPlayer.moveWithInterpolation(backendPlayer.positions)
              if (this.playerId === id) {
                const lastProcessedInput = this.inputs.findIndex((input) => {
                  return backendPlayer.seq === input.seq
                })
                if (lastProcessedInput > -1) {
                  this.inputs.splice(0, lastProcessedInput + 1)
                }
                this.inputs.forEach((input) => {
                  this.updateMouseCoordinate(input.x, input.y)
                  const worldCoordinate = this.camera.screenToWorld(
                    input.x,
                    input.y
                  )
                  if (this.currentPlayer) {
                    this.currentPlayer.move(
                      worldCoordinate.x,
                      worldCoordinate.y
                    )
                  }
                })
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
      player.draw(this.ctx, this.camera)
    }
  }

  renderStats(): void {
    this.stats.renderStats(this.ctx)
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

  update(): void {
    if (this.currentPlayer) {
      const head = this.currentPlayer.positions[0]
      const cameraX = clamp(
        head.x - this.ctx.canvas.width / 2,
        Constants.worldBoundary.minX,
        Constants.worldBoundary.maxX - this.ctx.canvas.width
      )
      const cameraY = clamp(
        head.y - this.ctx.canvas.height / 2,
        Constants.worldBoundary.minY,
        Constants.worldBoundary.maxY - this.ctx.canvas.height
      )
      this.camera.follow(cameraX, cameraY)
      this.stats.updateCameraCoordinate(cameraX, cameraY)
    }
  }

  renderBackground(): void {
    this.hexGrid.render(this.ctx, this.camera)
  }

  renderWorldBoundary(): void {
    const worldCenterX =
      (Constants.worldBoundary.minX + Constants.worldBoundary.maxX) / 2
    const worldCenterY =
      (Constants.worldBoundary.minY + Constants.worldBoundary.maxY) / 2
    const screenCenter = this.camera.worldToScreen(worldCenterX, worldCenterY)

    const screenRadius = this.camera.worldToScreenDistance(
      Constants.worldBoundary.radius
    )
    this.ctx.beginPath()
    this.ctx.strokeStyle = "rgba(255, 0, 0, 0.5)"
    this.ctx.lineWidth = 5
    this.ctx.arc(screenCenter.x, screenCenter.y, screenRadius, 0, Math.PI * 2)
    this.ctx.stroke()
  }

  render(): void {
    this.ctx.clearRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height)
    this.renderBackground()
    this.renderWorldBoundary()
    this.renderPlayers()
    this.renderStats()
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this.camera.width = width
    this.camera.height = height
    this.stats.updateCameraWidthAndHeight(width, height)
  }

  gameLoop(): void {
    this.update()
    this.render()
    this.stats.pingCooldown -= 1
    if (this.stats.pingCooldown <= 0) {
      this.sendPingPayload()
    }
    requestAnimationFrame(() => {
      this.gameLoop()
    })
  }

  start(): void {
    this.stats.calculateFps()
    this.gameLoop()
  }
}
