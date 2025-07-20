import {
  Graphics,
  Sprite,
  Texture,
  Text
} from "pixi.js"
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
  static sharedTexture: Texture | null = null

  constructor(game: Game, id: string, snake: Snake) {
    this.id = id
    this.game = game
    this.snake = snake
    this.sprite = []
    this.lastUpdatedTime = 0
    this.animationCounter = 0
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

  updateSegments(positions: Coordinate[], initial = false): {
    oldSnakeLength: number
    newSnakeLength: number
    lengthIncrease: number
  } {
    if (!Player.sharedTexture) {
      Player.sharedTexture =
        this.game.displayDriver.renderer.app.renderer.generateTexture(
          new Graphics()
            .circle(0, 0, Constants.snakeSegmentRadius)
            .fill({ color: this.snake.color })
            .stroke({ color: 0x000000, width: 1 })
        )
    }
    const lengthIncrease = positions.length - this.snake.segments.length
    const oldSnakeLength = this.snake.segments.length
    this.snake.segments = positions
    if (lengthIncrease === 0) {
      this.updateSprite()
    } else if (lengthIncrease > 0) {
      for (let i = 0; i < lengthIncrease; i++) {
        const sprite = new Sprite(Player.sharedTexture)
        sprite.anchor.set(0.5, 0.5)
        sprite.zIndex = 999 - oldSnakeLength - i
        this.sprite.push(sprite)
      }
      this.updateSprite()
    } else {
      for (let i = 0; i < -lengthIncrease; i++) {
        const sprite = this.sprite.pop()
        if (sprite) {
          sprite.destroy(true)
        }
      }
      this.updateSprite()
    }
    if (initial) {
      const head = this.sprite[0]
      const username = new Text({
        text: this.game.username,
        anchor: 0.5,
        style: { fontSize: 20, fill: 0x0000ff }
      })
      username.position.set(0.1 * head.width, 0.1 * head.height)
      head.addChild(username)
    }
    return {
      oldSnakeLength,
      newSnakeLength: this.snake.segments.length,
      lengthIncrease
    }
  }

  updateSprite(): void {
    for (let i = 0; i < this.snake.segments.length; i++) {
      const segment = this.snake.segments[i]
      const sprite = this.sprite[i]
      sprite.position.set(segment.x, segment.y)
    }
  }

  destroy(): void {
    this.sprite.forEach((sprite) => sprite.destroy({ children: true, texture: false, textureSource: false }))
    this.sprite = []
  }
}
