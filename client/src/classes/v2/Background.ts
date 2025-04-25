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
    const verticalSpacing = (hexHeight * 8) / 4
    const horizontalSpacing = hexWidth / 1.5
    const padding = Constants.worldBoundary.padding

    // Create hex
    for (let startX = Constants.worldBoundary.minX + padding;
         startX <= Constants.worldBoundary.maxX;
         startX += horizontalSpacing) {
      for (let startY = Constants.worldBoundary.minY + padding;
           startY <= Constants.worldBoundary.maxY;
           startY += verticalSpacing) {
        const points: Point[] = []
        for (let i = 0; i < 6; i++) {
          const angle = (Math.PI / 3) * i
          const hx = startX + hexSize * Math.cos(angle)
          const hy = startY + hexSize * Math.sin(angle)
          points.push(new Point(hx, hy))
        }
        this.container.addChild(new Graphics()
          .poly(points)
          .stroke({color: 0x808080, width: 1, alpha: 0.75}))
      }
    }

    // Create boundary
    this.container.addChild(new Graphics()
      .circle(0, 0, Constants.worldBoundary.radius)
      .stroke({color: 0xff0000, width: 2, alpha: 1}))
  }
}
