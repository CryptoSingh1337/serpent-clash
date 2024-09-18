import type { Position } from '../utils/types'

export class Player {
  id: string
  positions: Position[]
  color: string
  direction: number

  constructor({ id, color, positions, direction }: {id: string; color: string; positions: Position[]; direction: number }) {
    this.id = id
    this.color = color
    this.positions = positions
    this.direction = direction
  }

  draw(c: CanvasRenderingContext2D) {
    c.beginPath()
    c.arc(this.positions[0].x, this.positions[0].y, 5, 0, Math.PI * 2, false)
    c.fillStyle = this.color
    c.fill()
  }
}
