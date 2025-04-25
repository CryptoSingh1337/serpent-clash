import {Application, Container, type Sprite} from "pixi.js"
import type {Game} from "@/classes/v2/Game.ts"
import {Background} from "@/classes/v2/Background.ts"

export class GameRenderer {
  app: Application
  game: Game
  worldContainer: Container
  entityLayer: Container
  background: Background

  constructor(game: Game) {
    this.app = new Application()
    this.game = game
    this.worldContainer = new Container()
    this.entityLayer = new Container()
    this.background = new Background(this.game)
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
    this.worldContainer.addChild(this.entityLayer)
    this.app.stage.addChild(this.worldContainer)
    this.game.div.appendChild(this.app.canvas)
  }

  addEntity(sprites: Sprite[]) {
    sprites.forEach(sprite => this.entityLayer.addChild(sprite))
  }

  removeEntity() {
    this.entityLayer.removeChildAt(0)
  }

  stop(): void {
    this.app.destroy(true, true)
  }
}
