import type { Game } from "@/classes/Game.ts"
import { GameRenderer } from "@/classes/GameRenderer.ts"
import { Camera } from "@/classes/Camera.ts"

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

  update(): void {
    this.camera.update()
    this.renderer.render()
  }

  stop(): void {
    this.renderer.stop()
  }
}
