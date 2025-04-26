import { type Ref, type ShallowRef } from "vue"
import { Constants, WsMessageType } from "@/utils/constants.ts"
import type {
  BackendPlayer,
  Coordinate,
  Players,
  ReconcileEvent
} from "@/utils/types"
import { Player } from "@/classes/entity.ts"
import { clamp } from "@/utils/helper.ts"
import { SocketDriver } from "@/drivers/socket_driver.ts"
import { StatsDriver } from "@/drivers/stats_driver.ts"
import { DisplayDriver } from "@/drivers/display_driver.ts"
import Stats from "stats.js"

const debugMode: boolean = import.meta.env.VITE_DEBUG_MODE === "true"

export class GameDriver {
  ctx: CanvasRenderingContext2D
  stats: Stats | null = null
  displayDriver: DisplayDriver
  socketDriver: SocketDriver | null = null
  statsDriver: StatsDriver
  clientStatus: Ref
  playerId: string = ""
  username: string = ""
  mouseCoordinate: Coordinate
  frontendPlayers: Players = {}
  currentPlayer: Player | null = null
  inputs: ReconcileEvent[] = []
  seq: number = 0

  constructor({
    username,
    ctx,
    statsContainer,
    status
  }: {
    username: string
    ctx: CanvasRenderingContext2D
    statsContainer: Readonly<ShallowRef<HTMLDivElement | null>> | null
    status: Ref
  }) {
    if (!ctx) {
      throw new Error("Can't find canvas element")
    }
    this.ctx = ctx
    this.displayDriver = new DisplayDriver(ctx)
    this.statsDriver = new StatsDriver(this.displayDriver)
    this.statsDriver.updateCameraWidthAndHeight(
      ctx.canvas.width,
      ctx.canvas.height
    )
    this.mouseCoordinate = { x: 0, y: 0 }
    this.clientStatus = status
    this.username = username
    this.initSocket()
    this.initMouseControls()

    if (debugMode && statsContainer && statsContainer.value) {
      this.stats = new Stats()
      this.stats.showPanel(2)
      statsContainer.value.appendChild(this.stats.dom)
    }
  }

  mouseDownHandler(e: Event): void {
    console.log("Mouse down")
    if (this.currentPlayer) {
      this.currentPlayer.boost = true
    }
  }

  resetDefault(e: Event): void {
    console.log("Reset default")
    if (this.currentPlayer) {
      this.currentPlayer.boost = false
    }
  }

  initMouseControls(): void {
    this.mouseDownHandler = this.mouseDownHandler.bind(this)
    this.resetDefault = this.resetDefault.bind(this)
    this.ctx.canvas.addEventListener("mouseleave", this.resetDefault, true)
    this.ctx.canvas.addEventListener("mouseup", this.resetDefault, true)
    this.ctx.canvas.addEventListener("mousedown", this.mouseDownHandler, true)
    setInterval(
      (): void => {
        if (
          this.socketDriver &&
          this.socketDriver.getReadyState() === WebSocket.OPEN
        ) {
          const boost = this.currentPlayer?.boost || false
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
          const event: ReconcileEvent = {
            seq: ++this.seq,
            event: {
              coordinate: {
                x: this.mouseCoordinate.x,
                y: this.mouseCoordinate.y
              },
              boost
            }
          }
          this.inputs.push(event)
          this.statsDriver.updateReconcileEvent(this.inputs.length)
          this.socketDriver.send(
            JSON.stringify({
              type: WsMessageType.Movement,
              body: {
                seq: this.seq,
                coordinate: worldCoordinate,
                boost
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
      // switch (data.type) {
      //   case WsMessageType.hello: {
      //     this.playerId = body.id
      //     this.statsDriver.updatePlayerId(this.playerId)
      //     break
      //   }
      //   case WsMessageType.Pong: {
      //     const ping = Math.max(
      //       body.reqAck - body.reqInit + Date.now() - body.resInit,
      //       0
      //     )
      //     this.statsDriver.updatePing(ping)
      //     break
      //   }
      //   case WsMessageType.GameState: {
      //     const backendPlayers = body.playerStates as {
      //       [id: string]: BackendPlayer
      //     }
      //     for (const id in backendPlayers) {
      //       const backendPlayer = backendPlayers[id]
      //       if (!this.frontendPlayers[id]) {
      //         this.frontendPlayers[id] = new Player({
      //           id: id,
      //           color: backendPlayer.color,
      //           positions: backendPlayer.positions
      //         })
      //         if (!this.currentPlayer && this.playerId === id) {
      //           console.log("Current player changed")
      //           this.currentPlayer = this.frontendPlayers[id]
      //         }
      //       } else {
      //         const frontendPlayer = this.frontendPlayers[id]
      //         frontendPlayer.moveWithInterpolation(backendPlayer.positions)
      //         if (this.playerId === id) {
      //           const lastProcessedInput = this.inputs.findIndex((input) => {
      //             return backendPlayer.seq === input.seq
      //           })
      //           if (lastProcessedInput > -1) {
      //             this.inputs.splice(0, lastProcessedInput + 1)
      //           }
      //           this.inputs.forEach((input: ReconcileEvent) => {
      //             const { coordinate } = input.event
      //             if (coordinate) {
      //               this.mouseCoordinate.x = coordinate.x
      //               this.mouseCoordinate.y = coordinate.y
      //               this.statsDriver.updateMouseCoordinate(this.mouseCoordinate)
      //               const worldCoordinate =
      //                 this.displayDriver.getCameraScreenToWorldCoordinates(
      //                   coordinate.x,
      //                   coordinate.y
      //                 )
      //               if (this.currentPlayer) {
      //                 this.currentPlayer.move(
      //                   worldCoordinate.x,
      //                   worldCoordinate.y
      //                 )
      //               }
      //             }
      //           })
      //           this.statsDriver.updateHeadCoordinate(
      //             frontendPlayer.positions[0].x,
      //             frontendPlayer.positions[0].y
      //           )
      //         }
      //       }
      //     }
      //     for (const id in this.frontendPlayers) {
      //       if (!backendPlayers[id]) {
      //         delete this.frontendPlayers[id]
      //       }
      //     }
      //     break
      //   }
      //   default: {
      //     console.log("invalid message type", data.type)
      //   }
      // }
    }
    this.socketDriver = new SocketDriver(
      this.username,
      onOpen,
      onClose,
      onError,
      onMessage
    )
  }

  sendPingPayload(): void {
    if (
      this.socketDriver &&
      this.socketDriver.getReadyState() === WebSocket.OPEN
    ) {
      this.socketDriver.send(
        JSON.stringify({
          type: WsMessageType.Ping,
          body: {
            reqInit: Date.now()
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
    this.statsDriver.updateMouseCoordinate(this.mouseCoordinate)
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
    this.ctx.fillStyle = "#191825"
    this.ctx.fillRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height)
    this.displayDriver.renderHexGrid()
    this.displayDriver.renderWorldBoundary()
    this.displayDriver.renderPlayers(this.frontendPlayers)
    this.statsDriver.renderStats()
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this.displayDriver.updateCameraHeight(width, height)
    this.statsDriver.updateCameraWidthAndHeight(width, height)
  }

  gameLoop(): void {
    if (this.stats != null) {
      this.stats.begin()
    }
    this.update()
    this.render()
    this.statsDriver.reducePingCooldown()
    if (this.statsDriver.getPingCooldown() <= 0) {
      this.sendPingPayload()
    }
    if (this.stats != null) {
      this.stats.end()
    }
    requestAnimationFrame(() => {
      this.gameLoop()
    })
  }

  start(): void {
    this.gameLoop()
  }
}
