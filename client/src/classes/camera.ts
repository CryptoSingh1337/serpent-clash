import type { Coordinate } from "@/utils/types"

export class Camera {
  x: number = 0
  y: number = 0
  width: number
  height: number

  constructor(width: number, height: number) {
    this.width = width
    this.height = height
  }

  follow(x: number, y: number) {
    this.x = x
    this.y = y
  }

  worldToScreen(worldX: number, worldY: number): Coordinate {
    return {
      x: Math.floor(worldX - this.x),
      y: Math.floor(worldY - this.y)
    }
  }

  screenToWorld(screenX: number, screenY: number): Coordinate {
    return {
      x: Math.floor(screenX + this.x),
      y: Math.floor(screenY + this.y)
    }
  }
}
