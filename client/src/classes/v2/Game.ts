import type {Ref} from "vue"
import {DisplayDriver} from "@/classes/v2/DisplayDriver.ts"
import {Player} from "@/classes/v2/Player.ts"
import {Snake} from "@/classes/v2/Snake.ts"
import type {Coordinate} from "@/utils/types"
import {Constants} from "@/utils/constants.ts"
import {InputManager} from "@/classes/v2/InputManager.ts"

export class Game {
  div: HTMLDivElement|null
  displayDriver: DisplayDriver
  inputManager: InputManager
  player: Player|null
  pointer: Ref<{x: number, y: number}>

  constructor(div: HTMLDivElement|null, pointer: Ref<{x: number, y: number}>) {
    this.div = div
    this.displayDriver = new DisplayDriver(this)
    this.inputManager = new InputManager(this)
    this.player = null
    this.pointer = pointer
  }

  async init(): Promise<void> {
    await this.displayDriver.init()
    const theta = Math.atan2(0, 0)
    const x = 0
    const y = 0
    const segments: Coordinate[] = []
    segments.push({x: x, y: y})
    for (let i = 1; i < 10; i++) {
      segments.push({
        x: x - i*Constants.snakeSegmentDistance*Math.cos(theta),
        y: y - i*Constants.snakeSegmentDistance*Math.sin(theta),
      })
    }
    segments.reverse()
    this.player = new Player(this, "1", new Snake(segments, 0xffffff))
    console.log(this.player.snake.segments)
    this.displayDriver.renderer.addEntity(this.player.sprite)
  }

  removeEntity(entityId: string, entityType: string): void {
    //TODO: remove player entity
  }

  update(): void {
    if (this.player && this.player.snake) {
      const coordinate = this.player.snake.segments[0]
      this.pointer.value.x = coordinate.x
      this.pointer.value.y = coordinate.y
    }
    this.displayDriver.camera.update()
  }

  start(): void {
    if (!this.player) {
      return
    }
    this.displayDriver.camera.follow(this.player)
    this.displayDriver.renderer.app.ticker.add(() => {
      this.displayDriver.render()
      this.update()
    })
  }

  stop(): void {
    this.player = null
    this.displayDriver.renderer.removeEntity()
    this.displayDriver.stop()
  }
}
