import {
  Application,
  Container,
  type Particle,
  ParticleContainer,
  type Sprite
} from "pixi.js"
import type { Game } from "@/classes/Game.ts"
import { Background } from "@/classes/Background.ts"

export class GameRenderer {
  app: Application
  game: Game
  worldContainer: Container
  background: Background
  playerEntityLayer: Container
  foodEntityLayer: ParticleContainer

  constructor(game: Game) {
    this.app = new Application()
    this.game = game
    this.worldContainer = new Container({
      isRenderGroup: true,
      cullable: true,
      cullableChildren: true
    })
    this.background = new Background(this.game)
    this.playerEntityLayer = new Container({
      isRenderGroup: true,
      cullable: true,
      cullableChildren: true
    })
    this.foodEntityLayer = new ParticleContainer({
      dynamicProperties: {
        position: true
      },
      isRenderGroup: true,
      cullable: true,
      cullableChildren: true
    })
    this.foodEntityLayer.cullable = true
    this.foodEntityLayer.cullableChildren = true
  }

  async init(): Promise<void> {
    if (!this.game.div) {
      throw Error("invalid canvas")
    }
    await this.app.init({
      preference: "webgpu",
      resizeTo: window,
      background: 0x191825,
      antialias: true
    })
    this.background.init()
    this.worldContainer.addChild(this.background.container)
    this.worldContainer.addChild(this.foodEntityLayer)
    this.worldContainer.addChild(this.playerEntityLayer)
    this.app.stage.addChild(this.worldContainer)
    this.game.div.appendChild(this.app.canvas)
  }

  render(): void {
    for (const id in this.game.playerEntities) {
      const player = this.game.playerEntities[id]
      player.updateSprite()
    }
  }

  addSpriteEntity(entityType: string, sprites: Sprite[]) {
    if (entityType === "player") {
      this.playerEntityLayer.addChild(...sprites)
    }
  }

  addParticleEntity(particle: Particle) {
    this.foodEntityLayer.addParticle(particle)
  }

  stop(): void {
    this.app.destroy(
      { removeView: true },
      {
        children: true,
        texture: true,
        textureSource: true
      }
    )
  }
}
