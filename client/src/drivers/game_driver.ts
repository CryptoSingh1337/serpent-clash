import type { Ref } from "vue"
import { Stats } from "@/classes/stats.ts"
import { Constants } from "@/utils/constants.ts"
import type { BackendPlayer, Players } from "@/utils/types"
import { Player } from "@/classes/entity.ts"
import { clamp } from "@/utils/helper.ts"
import { Camera } from "@/classes/camera.ts"
import { SocketDriver } from "@/drivers/socket_driver.ts"
import { HexGrid } from "@/classes/hex_grid.ts"

export class GameDriver {
  ctx: CanvasRenderingContext2D
  socketDriver: SocketDriver | null = null
  clientStatus: Ref
  stats: Stats
  playerId: string = ""
  mouseCoordinate: { x: number; y: number }
  camera: Camera
  hexGrid: HexGrid
  frontendPlayers: Players = {}
  currentPlayer: Player | null = null
  inputs: { seq: number; x: number; y: number }[] = []
  seq: number = 0

  constructor(ctx: CanvasRenderingContext2D, status: Ref) {
    if (!ctx) {
      throw new Error("Can't find canvas element")
    }
    this.ctx = ctx
    this.mouseCoordinate = { x: 0, y: 0 }
    this.stats = new Stats()
    this.stats.updateCameraWidthAndHeight(ctx.canvas.width, ctx.canvas.height)
    this.camera = new Camera(ctx.canvas.width, ctx.canvas.height)
    this.hexGrid = new HexGrid(50)
    this.clientStatus = status
    this.initSocket()
    this.setupMouseTracking()
  }

  setupMouseTracking(): void {
    setInterval(
      () => {
        if (
          this.socketDriver &&
          this.socketDriver.getReadyState() === WebSocket.OPEN
        ) {
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
          this.socketDriver.send(
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
    const onOpen = () => {
      console.log("Socket opened")
      this.clientStatus.value = "Disconnect"
      this.stats.updateStatus("online")
    }
    const onClose = () => {
      console.log("Socket closed")
      this.clientStatus.value = "Connect"
      this.stats.reset()
      this.currentPlayer = null
      for (const id in this.frontendPlayers) {
        delete this.frontendPlayers[id]
      }
    }
    const onError = (err: any) => {
      this.clientStatus.value = "Connect"
      throw err
    }
    const onMessage = (data: any) => {
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
    this.socketDriver = new SocketDriver(onOpen, onClose, onError, onMessage)
  }

  sendPingPayload(): void {
    if (
      this.socketDriver &&
      this.socketDriver.getReadyState() === WebSocket.OPEN
    ) {
      this.socketDriver.send(
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
    if (
      !this.socketDriver ||
      this.socketDriver.getReadyState() === WebSocket.CLOSED
    ) {
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
    if (
      this.socketDriver &&
      this.socketDriver.getReadyState() === WebSocket.OPEN
    ) {
      this.socketDriver.close()
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
