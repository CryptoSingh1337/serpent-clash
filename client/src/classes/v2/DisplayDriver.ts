import type { Game } from "@/classes/v2/Game.ts"
import { GameRenderer } from "@/classes/v2/GameRenderer.ts"
import { Camera } from "@/classes/v2/Camera.ts"

export class DisplayDriver {
  game: Game
  renderer: GameRenderer
  camera: Camera

  constructor(game: Game) {
    this.game = game
    this.renderer = new GameRenderer(this.game)
    this.camera = new Camera(this.game)
  }

  async init(): Promise<void> {
    await this.renderer.init()
  }

  render(): void {
    if (this.game.player) {
      const coordinate = this.camera.worldToScreen(
        this.game.inputManager.mousePosition.x,
        this.game.inputManager.mousePosition.y
      )
      this.game.player.move(coordinate.x, coordinate.y)
      this.game.player.updateSprite()
    }
  }

  update(): void {
    this.camera.update()
  }

  stop(): void {
    this.renderer.stop()
  }
}
