import type { Game } from "@/classes/Game.ts"
import { Player } from "@/classes/Player.ts"
import type { BackendFood, BackendPlayer, ReconcileEvent } from "@/utils/types"
import { getServerBaseUrl } from "@/utils/helper.ts"
import { WsMessageType } from "@/utils/constants.ts"
import { Snake } from "@/classes/Snake.ts"
import { Food } from "@/classes/Food.ts"

export class NetworkManager {
  game: Game
  socket: WebSocket

  constructor(game: Game, username: string) {
    this.game = game
    this.socket = this.init(username)
  }

  init(username: string): WebSocket {
    const baseUrl = getServerBaseUrl(true)
    const socket = new WebSocket(`${baseUrl}/game/ws?username=${username}`)
    socket.onopen = () => {
      console.log("Socket opened")
      this.game.clientStatusRef.value = "Disconnect"
    }
    socket.onclose = () => {
      console.log("Socket closed")
      this.game.clientStatusRef.value = "Connect"
      this.game.displayDriver.camera.target = null
      for (const id in this.game.playerEntities) {
        this.game.playerEntities[id].destroy()
        delete this.game.playerEntities[id]
      }
      this.game.player = null
    }
    socket.onerror = (err: any) => {
      this.game.clientStatusRef.value = "Connect"
      throw err
    }
    socket.onmessage = (data: any) => {
      data = JSON.parse(data.data)
      const body = data.body
      switch (data.type) {
        case WsMessageType.hello: {
          this.game.player = new Player(
            this.game,
            body.id,
            new Snake([], 0xffffff)
          )
          this.game.statsManager.updatePlayerId(this.game.player.id)
          break
        }
        case WsMessageType.Pong: {
          const ping = Math.max(
            body.reqAck - body.reqInit + Date.now() - body.resInit,
            0
          )
          this.game.statsManager.updatePing(ping)
          break
        }
        case WsMessageType.PlayerState: {
          const backendPlayerEntities = body.players as {
            [id: string]: BackendPlayer
          }
          for (const id in backendPlayerEntities) {
            const backendPlayer = backendPlayerEntities[id]
            if (!this.game.playerEntities[id]) {
              // entity does not exists
              if (this.game.player && this.game.player.id === id) {
                // current player
                this.game.player.username = backendPlayer.username
                this.game.player.updateSegments(backendPlayer.positions, true)
                this.game.displayDriver.renderer.addSpriteEntity(
                  "player",
                  this.game.player.sprite
                )
                this.game.playerEntities[id] = this.game.player
              } else {
                // different player
                this.game.playerEntities[id] = new Player(
                  this.game,
                  id,
                  new Snake([], 0xffffff)
                )
                this.game.playerEntities[id].username = backendPlayer.username
                this.game.playerEntities[id].updateSegments(
                  backendPlayer.positions,
                  true
                )
                this.game.displayDriver.renderer.addSpriteEntity(
                  "player",
                  this.game.playerEntities[id].sprite
                )
              }
            } else {
              // already existing entity
              const playerEntity = this.game.playerEntities[id]
              const delta = playerEntity.updateSegments(
                backendPlayer.positions,
                false
              )
              if (delta.oldSnakeLength !== delta.newSnakeLength) {
                this.game.displayDriver.renderer.addSpriteEntity(
                  "player",
                  playerEntity.sprite.slice(
                    delta.oldSnakeLength,
                    delta.newSnakeLength
                  )
                )
              }
              playerEntity.moveWithInterpolation(backendPlayer.positions)
              if (this.game.player && this.game.player.id === id) {
                const lastProcessedInput =
                  this.game.inputManager.inputQueue.findIndex((input) => {
                    return backendPlayer.seq === input.seq
                  })
                if (lastProcessedInput > -1) {
                  this.game.inputManager.inputQueue.splice(
                    0,
                    lastProcessedInput + 1
                  )
                }
                this.game.inputManager.inputQueue.forEach(
                  (input: ReconcileEvent) => {
                    const { coordinate } = input.event
                    if (coordinate && this.game.player) {
                      this.game.player.move(coordinate.x, coordinate.y)
                    }
                  }
                )
              }
            }
          }
          for (const id in this.game.playerEntities) {
            if (!backendPlayerEntities[id]) {
              if (this.game.playerEntities[id]) {
                this.game.playerEntities[id].destroy()
                delete this.game.playerEntities[id]
              }
            }
          }
          break
        }
        case WsMessageType.FoodState: {
          const backendFoodEntities = body.foods as {
            [id: number]: BackendFood
          }
          for (const id in backendFoodEntities) {
            const backendFood = backendFoodEntities[id]
            if (!this.game.foodEntities[id]) {
              // food does not exist
              this.game.foodEntities[id] = new Food(
                this.game,
                id,
                backendFood.coordinate
              )
              this.game.foodEntities[id].create()
              if (this.game.foodEntities[id].particle) {
                this.game.displayDriver.renderer.addParticleEntity(
                  this.game.foodEntities[id].particle
                )
              }
            } else {
              // already existing food
              // const foodEntity = this.game.foodEntities[id]
              // foodEntity.coordinate = backendFood.coordinate
            }
            for (const id in this.game.foodEntities) {
              // remove food entities that are not in the backend state
              if (!backendFoodEntities[id]) {
                this.game.foodEntities[id].destroy()
                delete this.game.foodEntities[id]
              }
            }
          }
          break
        }
        default: {
          console.log("invalid message type", data.type)
        }
      }
    }
    return socket
  }

  socketState(): number {
    return this.socket.readyState
  }

  send(data: string): void {
    this.socket.send(data)
  }

  sendPingPayload(): void {
    if (this.socket && this.socketState() === WebSocket.OPEN) {
      this.socket.send(
        JSON.stringify({
          type: WsMessageType.Ping,
          body: {
            reqInit: Date.now()
          }
        })
      )
      this.game.statsManager.resetPingCooldown()
    }
  }

  close(): void {
    this.socket.close()
  }
}
