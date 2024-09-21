import type { Coordinate } from "@/utils/types"
import { Constants } from "@/utils/constants"

export class Player {
  id: string
  positions: Coordinate[]
  radius: number
  color: string

  constructor({ id, color, radius, positions }: {id: string; color: string; radius:number; positions: Coordinate[] }) {
    this.id = id
    this.color = color
    this.radius = radius
    this.positions = positions
  }

  draw(c: CanvasRenderingContext2D) {
    this.positions.reverse()
    this.positions.forEach((segment, index) => {
      c.beginPath();
      c.arc(segment.x, segment.y, Constants.snakeSegmentRadius / 2, 0, Math.PI * 2);
      c.fillStyle = `hsl(${(index / this.positions.length) * 360}, 100%, 50%)`; // Gradient color effect
      c.fill();
      c.stroke();
    })
    this.positions.reverse()
  }
}
