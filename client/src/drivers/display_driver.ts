import { Camera } from "@/classes/camera.ts"
import type { CameraCoordinates, Coordinate, Players } from "@/utils/types"
import { Constants } from "@/utils/constants.ts"
import { HexGrid } from "@/classes/hex_grid.ts"
import type { CustomStats } from "@/classes/custom_stats.ts"

export class DisplayDriver {
  ctx: CanvasRenderingContext2D
  camera: Camera
  hexGrid: HexGrid

  constructor(ctx: CanvasRenderingContext2D) {
    this.ctx = ctx
    this.camera = new Camera(ctx.canvas.width, ctx.canvas.height)
    this.hexGrid = new HexGrid(this, 50)
  }

  renderHex(x: number, y: number, hexSize: number): void {
    this.ctx.beginPath()
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i
      const hx = x + hexSize * Math.cos(angle)
      const hy = y + hexSize * Math.sin(angle)
      if (i === 0) {
        this.ctx.moveTo(hx, hy)
      } else {
        this.ctx.lineTo(hx, hy)
      }
    }
    this.ctx.closePath()
    this.ctx.stroke()
  }

  renderPlayers(players: Players): void {
    for (const id in players) {
      const player = players[id]
      player.draw(this.ctx, this.camera)
    }
  }

  renderWorldBoundary(): void {
    const worldCenterX =
      (Constants.worldBoundary.minX + Constants.worldBoundary.maxX) / 2
    const worldCenterY =
      (Constants.worldBoundary.minY + Constants.worldBoundary.maxY) / 2
    const screenCenter = this.getCameraWorldToScreenCoordinates(
      worldCenterX,
      worldCenterY
    )

    const screenRadius = Constants.worldBoundary.radius
    this.ctx.beginPath()
    this.ctx.strokeStyle = "rgba(255, 0, 0, 0.5)"
    this.ctx.lineWidth = 5
    this.ctx.arc(screenCenter.x, screenCenter.y, screenRadius, 0, Math.PI * 2)
    this.ctx.stroke()
  }

  renderStats(stats: CustomStats): void {
    this.ctx.fillStyle = "White"
    this.ctx.font = "normal 12px Arial"
    this.ctx.fillText(`Status: ${stats.status}`, 5, 60)
    this.ctx.fillText(`Ping: ${Math.trunc(stats.ping * 100) / 100} ms`, 5, 75)
  }

  renderHexGrid(): void {
    this.hexGrid.render()
  }

  updateCameraCoordinates(x: number, y: number): void {
    this.camera.follow(x, y)
  }

  updateCanvasStrokeStyleAndLineWidth(
    strokeStyle: string,
    lineWidth: number
  ): void {
    this.ctx.strokeStyle = strokeStyle
    this.ctx.lineWidth = lineWidth
  }

  updateCameraHeight(width: number, height: number): void {
    this.camera.width = width
    this.camera.height = height
  }

  getCameraWorldToScreenCoordinates(
    worldX: number,
    worldY: number
  ): Coordinate {
    return this.camera.worldToScreen(worldX, worldY)
  }

  getCameraScreenToWorldCoordinates(
    screenX: number,
    screenY: number
  ): Coordinate {
    return this.camera.screenToWorld(screenX, screenY)
  }

  getCameraCoordinates(): CameraCoordinates {
    return {
      x: this.camera.x,
      y: this.camera.y,
      width: this.camera.width,
      height: this.camera.height
    }
  }
}
