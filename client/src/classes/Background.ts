import { Container, Graphics, Point, Sprite } from "pixi.js"
import type { Game } from "@/classes/Game.ts"
import { Constants } from "@/utils/constants.ts"

export class Background {
  game: Game
  container: Container

  constructor(game: Game) {
    this.game = game
    this.container = new Container()
  }

  init() {
    const hexSize = 50
    const hexWidth = 5 * hexSize
    const hexHeight = 2 * hexSize
    const verticalSpacing = (hexHeight * 6) / 4.5
    const horizontalSpacing = hexWidth / 2

    // Create hex in honeycomb pattern
    const hexTexture =
      this.game.displayDriver.renderer.app.renderer.generateTexture(
        this.drawHex(0, 0, hexSize)
      )
    let rowCount = 0
    for (
      let y = Constants.worldBoundary.minY;
      y <= Constants.worldBoundary.maxY;
      y += horizontalSpacing, rowCount++
    ) {
      const rowOffset = rowCount % 2 === 0 ? 0 : verticalSpacing / 2
      for (
        let x = Constants.worldBoundary.minX;
        x <= Constants.worldBoundary.maxX;
        x += verticalSpacing
      ) {
        const hexSprite = new Sprite(hexTexture)
        hexSprite.position.set(x + rowOffset, y)
        hexSprite.anchor.set(0.5)
        this.container.addChild(hexSprite)
      }
    }

    // Create outer mask and boundary
    const outerMask = new Graphics()
      .rect(
        Constants.worldBoundary.minX - Constants.worldBoundary.padding,
        Constants.worldBoundary.minY - Constants.worldBoundary.padding,
        Constants.worldBoundary.maxX -
          Constants.worldBoundary.minX +
          Constants.worldBoundary.padding,
        Constants.worldBoundary.maxY -
          Constants.worldBoundary.minY +
          Constants.worldBoundary.padding
      )
      .fill({ color: 0xff0000, alpha: 0.1 })
      .circle(0, 0, Constants.worldBoundary.radius)
      .cut()
      .circle(0, 0, Constants.worldBoundary.radius)
      .stroke({ color: 0xff0000, width: 2, alpha: 1 })
    this.container.addChild(outerMask)
    this.container.cacheAsTexture(true)
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
      .stroke({ color: 0x808080, width: 1, alpha: 0.85 })
  }
}
