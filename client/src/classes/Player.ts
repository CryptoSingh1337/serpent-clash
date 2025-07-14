import { Graphics, Sprite } from "pixi.js"
import type { Snake } from "@/classes/Snake.ts"
import type { Game } from "@/classes/Game.ts"
import { Constants } from "@/utils/constants.ts"
import { lerp, lerpAngle } from "@/utils/helper.ts"
import type { Coordinate } from "@/utils/types"

export class Player {
  game: Game
  id: string
  snake: Snake
  sprite: Sprite[]
  lastUpdatedTime: number
  animationCounter: number

  constructor(game: Game, id: string, snake: Snake) {
    this.id = id
    this.game = game
    this.snake = snake
    this.sprite = []
    this.lastUpdatedTime = 0
    this.animationCounter = 0
  }

  createSprite(): void {
    const texture =
      this.game.displayDriver.renderer.app.renderer.generateTexture(
        new Graphics()
          .circle(0, 0, Constants.snakeSegmentRadius)
          .fill({ color: this.snake.color })
          .stroke({ color: 0x000000, width: 1 })
      )

    for (let i = 0; i < this.snake.segments.length; i++) {
      const segment = this.snake.segments[i]
      const sprite = new Sprite(texture)
      sprite.anchor.set(0.5, 0.5)
      sprite.zIndex = i === 0 ? 999 : 998 - i
      sprite.position.set(segment.x, segment.y)
      this.sprite.push(sprite)
    }
  }

  move(x: number, y: number): void {
    const currentTime = performance.now()
    const deltaTime = currentTime - this.lastUpdatedTime

    if (deltaTime < Math.floor(1000 / Constants.tickRate)) {
      return
    }

    const head = this.snake.segments[0]
    let angle = this.snake.angle
    const targetAngle = Math.atan2(y - head.y, x - head.x)
    angle = lerpAngle(angle, targetAngle, Constants.maxTurnRate)
    this.snake.angle = angle

    // Move target head
    let speed = Constants.playerSpeed
    if (this.game.inputManager.boost) {
      speed += Constants.playerSpeedBoost
    }
    head.x += Math.cos(angle) * speed
    head.y += Math.sin(angle) * speed
    this.snake.segments[0] = head

    // Update target positions for the rest of the body
    for (let i = 1; i < this.snake.segments.length; i++) {
      const prevSegment = this.snake.segments[i - 1]
      const currentSegment = this.snake.segments[i]

      const angleToPrev = Math.atan2(
        prevSegment.y - currentSegment.y,
        prevSegment.x - currentSegment.x
      )

      currentSegment.x =
        prevSegment.x - Math.cos(angleToPrev) * Constants.snakeSegmentDistance
      currentSegment.y =
        prevSegment.y - Math.sin(angleToPrev) * Constants.snakeSegmentDistance
    }
    this.lastUpdatedTime = currentTime
  }

  moveWithInterpolation(positions: Coordinate[]): void {
    const currentTime = performance.now()
    const interpolationFactor = Math.min(
      (currentTime - this.lastUpdatedTime) / Constants.tickRate,
      1
    )
    for (let i = 0; i < this.snake.segments.length; i++) {
      this.snake.segments[i].x = lerp(
        this.snake.segments[i].x,
        positions[i].x,
        interpolationFactor
      )
      this.snake.segments[i].y = lerp(
        this.snake.segments[i].y,
        positions[i].y,
        interpolationFactor
      )
    }
    this.lastUpdatedTime = currentTime
  }

  updateSprite(): void {
    for (let i = 0; i < this.snake.segments.length; i++) {
      const segment = this.snake.segments[i]
      const sprite = this.sprite[i]
      sprite.position.set(segment.x, segment.y)
    }
  }

  destroy(): void {
    this.sprite.forEach((sprite) => sprite.destroy(true))
  }
}
