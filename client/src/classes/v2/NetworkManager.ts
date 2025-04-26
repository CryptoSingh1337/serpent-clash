import type { Game } from "@/classes/v2/Game.ts"
import { Player } from "@/classes/v2/Player.ts"
import type { BackendPlayer, ReconcileEvent } from "@/utils/types"
import { getServerBaseUrl } from "@/utils/helper.ts"
import { WsMessageType } from "@/utils/constants.ts"
import { Snake } from "@/classes/v2/Snake.ts"

export class NetworkManager {
  game: Game
  socket: WebSocket

  constructor(game: Game, username: string) {
    this.game = game
    this.socket = this.init(username)
  }

  init(username: string): WebSocket {
    const baseUrl = getServerBaseUrl(true)
    const socket = new WebSocket(`${baseUrl}/ws?username=${username}`)
    socket.onopen = () => {
      console.log("Socket opened")
    }
    socket.onclose = () => {
      console.log("Socket closed")
      this.game.player = null
      this.game.displayDriver.camera.target = null
      for (const id in this.game.playerEntities) {
        this.game.playerEntities[id].destroy()
        delete this.game.playerEntities[id]
      }
    }
    socket.onerror = (err: any) => {
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
        case WsMessageType.GameState: {
          const backendPlayerEntities = body.playerStates as {
            [id: string]: BackendPlayer
          }
          for (const id in backendPlayerEntities) {
            const backendPlayer = backendPlayerEntities[id]
            if (!this.game.playerEntities[id]) {
              if (this.game.player && this.game.player.id === id) {
                this.game.player.snake.segments = backendPlayer.positions
                this.game.player.createSprite()
                this.game.displayDriver.renderer.addEntity(
                  this.game.player.sprite
                )
                this.game.playerEntities[id] = this.game.player
              } else {
                this.game.playerEntities[id] = new Player(
                  this.game,
                  id,
                  new Snake(backendPlayer.positions, 0xffffff)
                )
                this.game.playerEntities[id].createSprite()
                this.game.displayDriver.renderer.addEntity(
                  this.game.playerEntities[id].sprite
                )
              }
            } else {
              const playerEntity = this.game.playerEntities[id]
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
              this.game.playerEntities[id].destroy()
              delete this.game.playerEntities[id]
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

  close(): void {
    this.socket.close()
  }
}
