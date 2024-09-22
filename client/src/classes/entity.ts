import type { Coordinate } from "@/utils/types"
import { Constants } from "@/utils/constants"
import { lerpAngle } from "@/utils/helper"

export class Player {
  id: string
  positions: Coordinate[]
  radius: number
  color: string
  angle: number = 0
  lastUpdatedTime: number = 0

  constructor({
    id,
    color,
    radius,
    positions
  }: {
    id: string
    color: string
    radius: number
    positions: Coordinate[]
  }) {
    this.id = id
    this.color = color
    this.radius = radius
    this.positions = positions
  }

  move(x: number, y: number): void {
    const currentTime = performance.now()
    if (currentTime - this.lastUpdatedTime <= 1000 / Constants.tickRate) {
      return
    }
    const head = this.positions[0]
    const targetAngle = Math.atan2(y - head.y, x - head.x)
    this.angle = lerpAngle(this.angle, targetAngle, Constants.maxTurnRate)

    head.x += Math.cos(this.angle) * Constants.playerSpeed
    head.y += Math.sin(this.angle) * Constants.playerSpeed

    this.positions[0] = head

    for (let i = 1; i < this.positions.length; i++) {
      const prevSegment = this.positions[i - 1]
      const currentSegment = this.positions[i]

      const angleToPrev = Math.atan2(prevSegment.y - currentSegment.y, prevSegment.x - currentSegment.x)

      currentSegment.x = prevSegment.x - Math.cos(angleToPrev) * Constants.snakeSegmentDistance
      currentSegment.y = prevSegment.y - Math.sin(angleToPrev) * Constants.snakeSegmentDistance
      this.positions[i] = currentSegment
      this.lastUpdatedTime = currentTime
    }
  }

  draw(c: CanvasRenderingContext2D): void {
    this.positions.reverse()
    this.positions.forEach((segment, index) => {
      c.beginPath()
      c.arc(
        segment.x,
        segment.y,
        Constants.snakeSegmentRadius / 2,
        0,
        Math.PI * 2
      )
      c.fillStyle = `hsl(${(index / this.positions.length) * 360}, 100%, 50%)` // Gradient color effect
      c.fill()
      c.stroke()
    })
    this.positions.reverse()
  }
}
