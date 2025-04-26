import { clamp, getServerBaseUrl } from "@/utils/helper.ts"
import { Constants } from "@/utils/constants.ts"
import type { Coordinate } from "@/utils/types"
import type { Game } from "@/classes/v2/Game.ts"

export class DebugManager {
  game: Game
  serverBaseUrl: string

  constructor(game: Game) {
    this.game = game
    this.serverBaseUrl = getServerBaseUrl(false)
  }

  async teleport(coordinate: Coordinate): Promise<void> {
    if (this.game.player && this.game.player.id === "") {
      console.log("invalid player id...")
      return
    }
    if (
      this.game.networkManager &&
      this.game.networkManager.socket &&
      this.game.networkManager.socketState() !== WebSocket.OPEN
    ) {
      console.log("Socket is closed...")
      return
    }
    console.log(`Teleporting to (${coordinate})`)
    await fetch(
      `${this.serverBaseUrl}/player/${this.game.player?.id}/teleport`,
      {
        method: "POST",
        headers: {
          "Content-type": "application/json"
        },
        body: JSON.stringify(coordinate)
      }
    )
    const { player, displayDriver, statsManager } = this.game
    if (player && statsManager) {
      const head = player.snake.segments[0]
      const cameraX = clamp(
        head.x - displayDriver.renderer.app.canvas.width / 2,
        Constants.worldBoundary.minX,
        Constants.worldBoundary.maxX - displayDriver.renderer.app.canvas.width
      )
      const cameraY = clamp(
        head.y - displayDriver.renderer.app.canvas.height / 2,
        Constants.worldBoundary.minY,
        Constants.worldBoundary.maxY - displayDriver.renderer.app.canvas.height
      )
      statsManager.updateCameraCoordinate(cameraX, cameraY)
    }
  }
}
