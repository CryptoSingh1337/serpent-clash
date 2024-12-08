import { clamp, getServerBaseUrl } from "@/utils/helper.ts"
import { Constants } from "@/utils/constants.ts"
import type { GameDriver } from "@/drivers/game_driver.ts"
import type { Coordinate } from "@/utils/types"

export class DebugDriver {
  game: GameDriver
  serverBaseUrl: string

  constructor(game: GameDriver) {
    this.game = game
    this.serverBaseUrl = getServerBaseUrl(false)
  }

  async teleport(coordinate: Coordinate): Promise<void> {
    if (this.game.playerId === "") {
      console.log("invalid player id...")
      return
    }
    if (
      this.game.socketDriver &&
      this.game.socketDriver.socket &&
      this.game.socketDriver.getReadyState() !== WebSocket.OPEN
    ) {
      console.log("Socket is closed...")
      return
    }
    console.log(`Teleporting to (${coordinate})`)
    await fetch(`${this.serverBaseUrl}/player/${this.game.playerId}/teleport`, {
      method: "POST",
      headers: {
        "Content-type": "application/json"
      },
      body: JSON.stringify(coordinate)
    })
    const { currentPlayer, displayDriver, statsDriver, ctx } = this.game
    if (currentPlayer && statsDriver && ctx && ctx.canvas) {
      const head = currentPlayer.positions[0]
      const cameraX = clamp(
        head.x - ctx.canvas.width / 2,
        Constants.worldBoundary.minX,
        Constants.worldBoundary.maxX - ctx.canvas.width
      )
      const cameraY = clamp(
        head.y - ctx.canvas.height / 2,
        Constants.worldBoundary.minY,
        Constants.worldBoundary.maxY - ctx.canvas.height
      )
      displayDriver.updateCameraCoordinates(cameraX, cameraY)
      statsDriver.updateCameraCoordinate(cameraX, cameraY)
    }
  }
}
