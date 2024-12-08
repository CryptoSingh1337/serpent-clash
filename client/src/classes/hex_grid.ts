import { Constants } from "@/utils/constants.ts"
import type { DisplayDriver } from "@/drivers/display_driver.ts"

export class HexGrid {
  hexSize: number
  hexHeight: number
  hexWidth: number
  verticalSpacing: number
  horizontalSpacing: number
  displayDriver: DisplayDriver

  constructor(displayDriver: DisplayDriver, hexSize: number) {
    this.hexSize = hexSize
    this.hexHeight = hexSize * 2
    this.hexWidth = Math.sqrt(25) * hexSize
    this.verticalSpacing = (this.hexHeight * 8) / 4
    this.horizontalSpacing = this.hexWidth / 1.5
    this.displayDriver = displayDriver
  }

  drawHex(x: number, y: number) {
    this.displayDriver.renderHex(x, y, this.hexSize)
  }

  render() {
    const { x, y, width, height } = this.displayDriver.getCameraCoordinates()
    const startX = Math.max(
      Constants.worldBoundary.minX,
      Math.floor(x / this.horizontalSpacing) * this.horizontalSpacing -
        this.horizontalSpacing
    )
    const startY = Math.max(
      Constants.worldBoundary.minY,
      Math.floor(y / this.verticalSpacing) * this.verticalSpacing -
        this.verticalSpacing
    )
    const endX = Math.min(
      Constants.worldBoundary.maxX,
      x + width + this.horizontalSpacing
    )
    const endY = Math.min(
      Constants.worldBoundary.maxY,
      y + height + this.verticalSpacing
    )

    this.displayDriver.updateCanvasStrokeStyleAndLineWidth(
      "rgba(200, 200, 200, 0.5)",
      1
    )

    for (let y = startY; y < endY; y += this.verticalSpacing) {
      for (let x = startX; x < endX; x += this.horizontalSpacing) {
        const screenPos = this.displayDriver.getCameraWorldToScreenCoordinates(
          x,
          y
        )
        this.drawHex(screenPos.x, screenPos.y)
        if (y + this.verticalSpacing / 2 < Constants.worldBoundary.maxY) {
          const offsetScreenPos =
            this.displayDriver.getCameraWorldToScreenCoordinates(
              x + this.horizontalSpacing / 2,
              y + this.verticalSpacing / 2
            )
          this.drawHex(offsetScreenPos.x, offsetScreenPos.y)
        }
      }
    }
  }
}
