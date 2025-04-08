import type { Coordinate } from "@/utils/types"
import { Constants } from "@/utils/constants"
import { lerp, lerpAngle } from "@/utils/helper"
import type { Camera } from "@/classes/camera"

export class Player {
  id: string
  positions: Coordinate[]
  color: string
  angle: number = 0
  lastUpdatedTime: number = 0
  lastServerUpdateTime: number = 0
  speedBoost: boolean = false

  constructor({
    id,
    color,
    positions
  }: {
    id: string
    color: string
    positions: Coordinate[]
  }) {
    this.id = id
    this.color = color
    this.positions = positions
    this.lastServerUpdateTime = performance.now()
  }

  move(x: number, y: number): void {
    const currentTime = performance.now()
    const deltaTime = currentTime - this.lastUpdatedTime

    if (deltaTime < Math.floor(1000 / Constants.tickRate)) {
      return
    }

    const head = this.positions[0]
    let angle = this.angle
    const targetAngle = Math.atan2(y - head.y, x - head.x)
    angle = lerpAngle(angle, targetAngle, Constants.maxTurnRate)

    // Move target head
    let speed = Constants.playerSpeed
    if (this.speedBoost) {
      speed += Constants.playerSpeedBoost
    }
    head.x += Math.cos(angle) * speed
    head.y += Math.sin(angle) * speed
    this.positions[0] = head

    // Update target positions for the rest of the body
    for (let i = 1; i < this.positions.length; i++) {
      const prevSegment = this.positions[i - 1]
      const currentSegment = this.positions[i]

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
    for (let i = 0; i < this.positions.length; i++) {
      this.positions[i].x = lerp(
        this.positions[i].x,
        positions[i].x,
        interpolationFactor
      )
      this.positions[i].y = lerp(
        this.positions[i].y,
        positions[i].y,
        interpolationFactor
      )
    }
    this.lastUpdatedTime = currentTime
  }

  draw(ctx: CanvasRenderingContext2D, camera: Camera): void {
    ctx.lineWidth = 1
    ctx.strokeStyle = "black"
    for (let i = this.positions.length - 1, j = 0; i >= 0; i--, j++) {
      const segment = this.positions[i]
      const screenPos = camera.worldToScreen(segment.x, segment.y)
      ctx.beginPath()
      ctx.arc(
        screenPos.x,
        screenPos.y,
        Constants.snakeSegmentDiameter / 2,
        0,
        Math.PI * 2
      )
      ctx.fillStyle = `hsl(${(j / this.positions.length) * 360}, 100%, 50%)`
      ctx.fill()
      ctx.stroke()
    }
  }
}
