import { Camera } from "@/classes/camera.ts"
import { Constants } from "@/utils/constants.ts"

export class HexGrid {
  hexSize: number
  hexHeight: number
  hexWidth: number
  verticalSpacing: number
  horizontalSpacing: number

  constructor(hexSize: number) {
    this.hexSize = hexSize
    this.hexHeight = hexSize * 2
    this.hexWidth = Math.sqrt(25) * hexSize
    this.verticalSpacing = (this.hexHeight * 8) / 4
    this.horizontalSpacing = this.hexWidth / 1.5
  }

  drawHex(ctx: CanvasRenderingContext2D, x: number, y: number) {
    ctx.beginPath()
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i
      const hx = x + this.hexSize * Math.cos(angle)
      const hy = y + this.hexSize * Math.sin(angle)
      if (i === 0) {
        ctx.moveTo(hx, hy)
      } else {
        ctx.lineTo(hx, hy)
      }
    }
    ctx.closePath()
    ctx.stroke()
  }

  render(ctx: CanvasRenderingContext2D, camera: Camera) {
    const startX = Math.max(
      Constants.worldBoundary.minX,
      Math.floor(camera.x / this.horizontalSpacing) * this.horizontalSpacing -
        this.horizontalSpacing
    )
    const startY = Math.max(
      Constants.worldBoundary.minY,
      Math.floor(camera.y / this.verticalSpacing) * this.verticalSpacing -
        this.verticalSpacing
    )
    const endX = Math.min(
      Constants.worldBoundary.maxX,
      camera.x + camera.width + this.horizontalSpacing
    )
    const endY = Math.min(
      Constants.worldBoundary.maxY,
      camera.y + camera.height + this.verticalSpacing
    )

    ctx.strokeStyle = "rgba(200, 200, 200, 0.5)"
    ctx.lineWidth = 1

    for (let y = startY; y < endY; y += this.verticalSpacing) {
      for (let x = startX; x < endX; x += this.horizontalSpacing) {
        const screenPos = camera.worldToScreen(x, y)
        this.drawHex(ctx, screenPos.x, screenPos.y)
        if (y + this.verticalSpacing / 2 < Constants.worldBoundary.maxY) {
          const offsetScreenPos = camera.worldToScreen(
            x + this.horizontalSpacing / 2,
            y + this.verticalSpacing / 2
          )
          this.drawHex(ctx, offsetScreenPos.x, offsetScreenPos.y)
        }
      }
    }
  }
}
