import { Graphics, Particle, Texture } from "pixi.js"
import type { Game } from "@/classes/Game.ts"
import {Constants} from "@/utils/constants.ts"

export class Food {
  game: Game
  id: string
  coordinate: { x: number; y: number }
  color: number
  particle: Particle | null

  static sharedTexture: Texture | null = null

  constructor(
    game: Game,
    id: string,
    coordinate: { x: number; y: number },
    color: number = 0x00ff00
  ) {
    this.game = game
    this.id = id
    this.coordinate = coordinate
    this.color = color
    this.particle = null
  }

  create(): void {
    if (!Food.sharedTexture) {
      const graphics = new Graphics()
        .circle(0, 0, Constants.foodRadius)
        .fill({ color: this.color })
      Food.sharedTexture =
        this.game.displayDriver.renderer.app.renderer.generateTexture(graphics)
    }
    this.particle = new Particle({
      texture: Food.sharedTexture,
      x: this.coordinate.x,
      y: this.coordinate.y,
      anchorX: 0.5,
      anchorY: 0.5
    })
  }

  destroy(): void {
    if (this.particle) {
      this.game.displayDriver.renderer.foodEntityLayer.removeParticle(
        this.particle
      )
      this.particle = null
    }
  }
}
