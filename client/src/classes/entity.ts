import type { Position } from '../utils/types'

export class Player {
  id: string
  positions: Position[]
  radius: number
  color: string
  direction: number

  constructor({ id, color, radius, positions, direction }: {id: string; color: string; radius:number; positions: Position[]; direction: number }) {
    this.id = id
    this.color = color
    this.radius = radius
    this.positions = positions
    this.direction = direction
  }

  draw(c: CanvasRenderingContext2D) {
    c.beginPath()
    c.arc(this.positions[0].x, this.positions[0].y, this.radius * devicePixelRatio,
      0, Math.PI * 2, false)
    c.fillStyle = this.color
    c.fill()
  }
}
