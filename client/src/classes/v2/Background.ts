import {Container, Graphics, Point} from "pixi.js"
import type {Game} from "@/classes/v2/Game.ts"
import {Constants} from "@/utils/constants.ts"

export class Background {
  game: Game
  container: Container

  constructor(game: Game) {
    this.game = game
    this.container = new Container()
  }

  init() {
    const hexSize = 50
    const hexHeight = hexSize * 2
    const hexWidth = Math.sqrt(25) * hexSize
    const verticalSpacing = (hexHeight * 6) / 4.5
    const horizontalSpacing = hexWidth / 2

    // Create hex in honeycomb pattern
    let rowCount = 0
    for (let startY = Constants.worldBoundary.minY;
         startY <= Constants.worldBoundary.maxY;
         startY += horizontalSpacing, rowCount++) {
      const rowOffset = rowCount % 2 === 0 ? 0 : verticalSpacing / 2
      for (let startX = Constants.worldBoundary.minX;
           startX <= Constants.worldBoundary.maxX;
           startX += verticalSpacing) {
        this.container.addChild(this.drawHex(startX + rowOffset, startY, hexSize))
      }
    }

    // Create boundary
    this.container.addChild(new Graphics()
      .circle(0, 0, Constants.worldBoundary.radius)
      .stroke({color: 0xff0000, width: 2, alpha: 1}))
  }

  drawHex(x: number, y: number, hexSize: number): Graphics {
    const points = []
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i
      const hx = x + hexSize * Math.cos(angle)
      const hy = y + hexSize * Math.sin(angle)
      points.push(new Point(hx, hy))
    }
    return new Graphics()
      .poly(points)
      .stroke({color: 0x808080, width: 1, alpha: 0.85})
  }
}
