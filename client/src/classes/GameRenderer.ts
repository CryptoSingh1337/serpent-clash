import { Application, Container, type Sprite } from "pixi.js"
import type { Game } from "@/classes/Game.ts"
import { Background } from "@/classes/Background.ts"

export class GameRenderer {
  app: Application
  game: Game
  worldContainer: Container
  background: Background
  entityLayer: Container

  constructor(game: Game) {
    this.app = new Application()
    this.game = game
    this.worldContainer = new Container()
    this.background = new Background(this.game)
    this.entityLayer = new Container()
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
    this.app.stage.addChild(this.worldContainer)
    this.worldContainer.addChild(this.background.container)
    this.worldContainer.addChild(this.entityLayer)
    this.game.div.appendChild(this.app.canvas)
  }

  render(): void {
    for (const id in this.game.playerEntities) {
      const player = this.game.playerEntities[id]
      player.updateSprite()
    }
  }

  addEntity(sprites: Sprite[]) {
    sprites.forEach((sprite) => this.entityLayer.addChild(sprite))
  }

  removeEntity() {
    this.entityLayer.children.forEach((object) => object.destroy())
  }

  stop(): void {
    this.app.destroy(true, true)
  }
}
