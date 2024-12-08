import { type Ref } from "vue"
import { Constants } from "@/utils/constants.ts"
import type {
  BackendPlayer,
  Coordinate,
  Players,
  ReconcileEvent,
  SpeedBoost
} from "@/utils/types"
import { Player } from "@/classes/entity.ts"
import { clamp } from "@/utils/helper.ts"
import { SocketDriver } from "@/drivers/socket_driver.ts"
import { StatsDriver } from "@/drivers/stats_driver.ts"
import { DisplayDriver } from "@/drivers/display_driver.ts"

export class GameDriver {
  ctx: CanvasRenderingContext2D
  displayDriver: DisplayDriver
  socketDriver: SocketDriver | null = null
  statsDriver: StatsDriver
  clientStatus: Ref
  playerId: string = ""
  mouseCoordinate: Coordinate
  frontendPlayers: Players = {}
  currentPlayer: Player | null = null
  inputs: ReconcileEvent[] = []
  seq: number = 0

  constructor(ctx: CanvasRenderingContext2D, status: Ref) {
    if (!ctx) {
      throw new Error("Can't find canvas element")
    }
    this.ctx = ctx
    this.displayDriver = new DisplayDriver(ctx)
    this.statsDriver = new StatsDriver()
    this.statsDriver.updateCameraWidthAndHeight(
      ctx.canvas.width,
      ctx.canvas.height
    )
    this.mouseCoordinate = { x: 0, y: 0 }
    this.clientStatus = status
    this.initSocket()
    this.initMouseControls()
  }

  initMouseControls(): void {
    this.ctx.canvas.addEventListener("mousedown", (): void => {
      if (
        this.socketDriver &&
        this.socketDriver.getReadyState() === WebSocket.OPEN
      ) {
        if (this.currentPlayer) {
          this.currentPlayer.speedBoost = true
        }
        const event: SpeedBoost = {
          enabled: true
        }
        this.inputs.push({
          seq: ++this.seq,
          event: event
        })
        this.socketDriver.send(
          JSON.stringify({
            type: "boost",
            body: event
          })
        )
      }
    })
    this.ctx.canvas.addEventListener("mouseup", (): void => {
      if (
        this.socketDriver &&
        this.socketDriver.getReadyState() === WebSocket.OPEN
      ) {
        if (this.currentPlayer) {
          this.currentPlayer.speedBoost = false
        }
        const event: SpeedBoost = {
          enabled: false
        }
        this.inputs.push({
          seq: ++this.seq,
          event: event
        })
        this.socketDriver.send(
          JSON.stringify({
            type: "boost",
            body: {
              seq: this.seq,
              event
            }
          })
        )
      }
    })
    setInterval(
      (): void => {
        if (
          this.socketDriver &&
          this.socketDriver.getReadyState() === WebSocket.OPEN
        ) {
          const worldCoordinate =
            this.displayDriver.getCameraScreenToWorldCoordinates(
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
            event: {
              x: this.mouseCoordinate.x,
              y: this.mouseCoordinate.y
            }
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
      this.statsDriver.updateStatus("online")
    }
    const onClose = () => {
      console.log("Socket closed")
      this.clientStatus.value = "Connect"
      this.statsDriver.reset()
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
          this.statsDriver.updatePlayerId(this.playerId)
          break
        }
        case "pong": {
          this.statsDriver.updatePing(performance.now() - body.timestamp)
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
                this.inputs.forEach((input: ReconcileEvent) => {
                  const { event } = input
                  if ("x" in event && "y" in event) {
                    this.updateMouseCoordinate(event.x, event.y)
                    const worldCoordinate =
                      this.displayDriver.getCameraScreenToWorldCoordinates(
                        event.x,
                        event.y
                      )
                    if (this.currentPlayer) {
                      this.currentPlayer.move(
                        worldCoordinate.x,
                        worldCoordinate.y
                      )
                    }
                  } else if ("enabled" in event) {
                    if (this.currentPlayer) {
                      this.currentPlayer.speedBoost = event.enabled
                    }
                  }
                })
                this.statsDriver.updateHeadCoordinate(
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
      this.statsDriver.resetPingCooldown()
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

  updateMouseCoordinate(x: number, y: number): void {
    this.mouseCoordinate.x = x
    this.mouseCoordinate.y = y
    this.statsDriver.updateMouseCoordinate(x, y)
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
      this.displayDriver.updateCameraCoordinates(cameraX, cameraY)
      this.statsDriver.updateCameraCoordinate(cameraX, cameraY)
    }
  }

  render(): void {
    this.ctx.clearRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height)
    this.displayDriver.drawHexGrid()
    this.displayDriver.drawWorldBoundary()
    this.displayDriver.drawPlayers(this.frontendPlayers)
    this.statsDriver.renderStats(this.ctx)
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this.displayDriver.updateCameraHeight(width, height)
    this.statsDriver.updateCameraWidthAndHeight(width, height)
  }

  gameLoop(): void {
    this.update()
    this.render()
    this.statsDriver.reducePingCooldown()
    if (this.statsDriver.getPingCooldown() <= 0) {
      this.sendPingPayload()
    }
    requestAnimationFrame(() => {
      this.gameLoop()
    })
  }

  start(): void {
    this.statsDriver.calculateFps()
    this.gameLoop()
  }
}
