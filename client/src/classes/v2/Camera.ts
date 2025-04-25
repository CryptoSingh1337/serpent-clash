import type {Game} from "@/classes/v2/Game.ts"
import type {Player} from "@/classes/v2/Player.ts"
import {Constants} from "@/utils/constants.ts"
import type {Coordinate} from "@/utils/types"

export class Camera {
  game: Game
  target: Player|null
  position: {x: number, y: number}

  constructor(game: Game) {
    this.game = game
    this.target = null
    this.position = {x: 0, y: 0}
  }

  follow(target: Player) {
    this.target = target
  }

  update(): void {
    if (!this.target) {
      return
    }
    const targetX = -this.target.sprite[0].x + this.game.displayDriver.renderer.app.screen.width / 2
    const targetY = -this.target.sprite[0].y + this.game.displayDriver.renderer.app.screen.height / 2
    const app = this.game.displayDriver.renderer.app
    const world = this.game.displayDriver.renderer.worldContainer
    world.x = world.x + (targetX - world.x)
    world.y = world.y + (targetY - world.y)
    const cameraPadding = Constants.worldBoundary.padding
    const leftBoundValue = Constants.worldBoundary.minX - cameraPadding
    const rightBoundValue = Constants.worldBoundary.maxX + cameraPadding - app.screen.width
    const topBoundValue = Constants.worldBoundary.minY - cameraPadding
    const bottomBoundValue = Constants.worldBoundary.maxY + cameraPadding - app.screen.height
    const maxX = -rightBoundValue
    const minX = -leftBoundValue
    const maxY = -bottomBoundValue
    const minY = -topBoundValue
    world.x = Math.max(Math.min(world.x, minX), maxX);
    world.y = Math.max(Math.min(world.y, minY), maxY);
    this.position.x = world.x
    this.position.y = world.y
  }

  worldToScreen(worldX: number, worldY: number): Coordinate {
    return {
      x: Math.floor(worldX - this.position.x),
      y: Math.floor(worldY - this.position.y)
    }
  }

  screenToWorld(screenX: number, screenY: number): Coordinate {
    return {
      x: Math.floor(screenX + this.position.x),
      y: Math.floor(screenY + this.position.y)
    }
  }
}
