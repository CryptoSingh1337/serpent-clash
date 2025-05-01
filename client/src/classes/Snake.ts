import type { Coordinate } from "@/utils/types"

export class Snake {
  angle: number
  segments: Coordinate[]
  color: number

  constructor(segments: Coordinate[], color: number) {
    this.angle = 0
    this.segments = segments
    this.color = color
  }
}
