import { clamp } from "@/utils/helper.ts"
import { Constants } from "@/utils/constants.ts"
import type { GameDriver } from "@/drivers/game_driver.ts"

export class DebugDriver {
  game: GameDriver

  constructor(game: GameDriver) {
    this.game = game
  }

  teleport(x: number, y: number): void {
    console.log(`Teleporting to (${x}, ${y})`)
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