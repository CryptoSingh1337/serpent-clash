import type {Game} from "@/classes/v2/Game.ts"
import {Player} from "@/classes/entity.ts"
import type {BackendPlayer, ReconcileEvent} from "@/utils/types"
import {getServerBaseUrl} from "@/utils/helper.ts"
import {WsMessageType} from "@/utils/constants.ts"
import type {Ref} from "vue";

export class NetworkManager {
  clientStatus: Ref
  game: Game
  socket: WebSocket

  constructor(game: Game, clientStatus: Ref, username: string) {
    this.clientStatus = clientStatus
    this.game = game
    this.socket = this.init(username)
  }

  init(
    username: string
  ): WebSocket {
    const baseUrl = getServerBaseUrl(true)
    const socket = new WebSocket(`${baseUrl}/ws?username=${username}`)
    socket.onopen = () => {
      console.log("Socket opened")
      // this.clientStatus.value = "Disconnect"
      // this.statsDriver.updateStatus("online")
    }
    socket.onclose = () => {
      console.log("Socket closed")
      // this.clientStatus.value = "Connect"
      // this.statsDriver.reset()
      // this.currentPlayer = null
      // for (const id in this.frontendPlayers) {
      //   delete this.frontendPlayers[id]
      // }
    }
    socket.onerror = (err: any) => {
      // this.clientStatus.value = "Connect"
      throw err
    }
    socket.onmessage = (data: any) => {
      data = JSON.parse(data.data)
      const body = data.body
      // switch (data.type) {
      //   case WsMessageType.hello: {
      //     this.game.player = body.id
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
