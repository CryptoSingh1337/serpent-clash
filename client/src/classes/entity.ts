import type { Coordinate } from "@/utils/types"
import { Constants } from "@/utils/constants"
import {lerp, lerpAngle} from "@/utils/helper"

export class Player {
  id: string
  positions: Coordinate[]
  targetPositions: Coordinate[]
  radius: number
  color: string
  angle: number = 0
  targetAngle: number = 0
  lastUpdatedTime: number = 0
  lastServerUpdateTime: number = 0

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
    this.targetPositions = JSON.parse(JSON.stringify(positions)) // Deep copy
    this.lastServerUpdateTime = performance.now()
  }

  move(x: number, y: number): void {
    const currentTime = performance.now()
    const deltaTime = (currentTime - this.lastUpdatedTime) / 1000 // Convert to seconds

    if (deltaTime < 1 / Constants.tickRate) {
      return
    }
    // Update target angle and position
    const head = this.targetPositions[0]
    this.targetAngle = Math.atan2(y - head.y, x - head.x)

    // Move target head
    head.x += Math.cos(this.targetAngle) * Constants.playerSpeed * deltaTime
    head.y += Math.sin(this.targetAngle) * Constants.playerSpeed * deltaTime

    // Update target positions for the rest of the body
    for (let i = 1; i < this.targetPositions.length; i++) {
      const prevSegment = this.targetPositions[i - 1]
      const currentSegment = this.targetPositions[i]

      const angleToPrev = Math.atan2(
          prevSegment.y - currentSegment.y,
          prevSegment.x - currentSegment.x
      )

      currentSegment.x = prevSegment.x - Math.cos(angleToPrev) * Constants.snakeSegmentDistance
      currentSegment.y = prevSegment.y - Math.sin(angleToPrev) * Constants.snakeSegmentDistance
    }

    // Interpolate actual positions
    const interpolationFactor = Math.min((currentTime - this.lastServerUpdateTime) / Constants.tickRate, 1)

    this.angle = lerpAngle(this.angle, this.targetAngle, Constants.maxTurnRate * deltaTime)

    for (let i = 0; i < this.positions.length; i++) {
      this.positions[i].x = lerp(this.positions[i].x, this.targetPositions[i].x, interpolationFactor)
      this.positions[i].y = lerp(this.positions[i].y, this.targetPositions[i].y, interpolationFactor)
    }
    this.lastUpdatedTime = currentTime
  }

  moveWithInterpolation(positions: Coordinate[]): void {
    const currentTime = performance.now()
    const interpolationFactor = Math.min((currentTime - this.lastUpdatedTime) / Constants.tickRate, 1)
    for (let i = 0; i < this.positions.length; i++) {
      this.positions[i].x = lerp(this.positions[i].x, positions[i].x, interpolationFactor)
      this.positions[i].y = lerp(this.positions[i].y, positions[i].y, interpolationFactor)
    }
    this.lastUpdatedTime = currentTime
  }

  updateFromServer(serverPositions: Coordinate[]): void {
    this.positions = serverPositions
    this.targetPositions = JSON.parse(JSON.stringify(serverPositions)) // Deep copy
    this.lastServerUpdateTime = performance.now()
  }

  draw(c: CanvasRenderingContext2D): void {
    this.positions.reverse()
    this.positions.forEach((segment, index) => {
      c.beginPath()
      c.arc(
        segment.x,
        segment.y,
        Constants.snakeSegmentDiameter / 2,
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
