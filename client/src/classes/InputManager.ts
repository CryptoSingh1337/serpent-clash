import type { Game } from "@/classes/Game.ts"
import type { ReconcileEvent } from "@/utils/types"
import { clamp } from "@/utils/helper.ts"
import { Constants, WsMessageType } from "@/utils/constants.ts"

export class InputManager {
  game: Game
  mousePosition: { x: number; y: number }
  boost: boolean
  inputQueue: ReconcileEvent[] = []
  seq: number = 0

  constructor(game: Game) {
    this.game = game
    this.mousePosition = { x: 0, y: 0 }
    this.boost = false
    window.addEventListener("mousemove", this.onMouseMove.bind(this))
    window.addEventListener("mouseup", this.onMouseUp.bind(this))
    window.addEventListener("mousedown", this.onMouseDown.bind(this))
    window.addEventListener("mouseleave", this.onMouseLeave.bind(this))
    setInterval(
      (): void => {
        if (
          this.game.networkManager &&
          this.game.networkManager.socketState() === WebSocket.OPEN
        ) {
          const boost = this.game.inputManager.boost || false
          const worldCoordinate = this.game.displayDriver.camera.screenToWorld(
            this.game.inputManager.mousePosition.x,
            this.game.inputManager.mousePosition.y
          )
          worldCoordinate.x = clamp(
            worldCoordinate.x,
            Constants.worldBoundary.minX,
            Constants.worldBoundary.maxX
          )
          worldCoordinate.y = clamp(
            worldCoordinate.y,
            Constants.worldBoundary.minY,
            Constants.worldBoundary.maxY
          )
          const event: ReconcileEvent = {
            seq: ++this.seq,
            event: {
              coordinate: {
                x: this.game.inputManager.mousePosition.x,
                y: this.game.inputManager.mousePosition.y
              },
              boost
            }
          }
          this.inputQueue.push(event)
          this.game.networkManager.send(
            JSON.stringify({
              type: WsMessageType.Movement,
              body: {
                seq: this.seq,
                coordinate: worldCoordinate,
                boost
              }
            })
          )
        }
      },
      Math.floor(1000 / Constants.tickRate)
    )
  }

  onMouseMove(e: MouseEvent) {
    if (!this.game.div) {
      return
    }
    const rect = this.game.div.getBoundingClientRect()
    this.mousePosition.x = e.clientX - rect.left
    this.mousePosition.y = e.clientY - rect.top
  }

  onMouseUp(e: MouseEvent) {
    if (e.button === 0) {
      this.boost = false
    }
  }

  onMouseDown(e: MouseEvent) {
    if (e.button === 0) {
      this.boost = true
    }
  }

  onMouseLeave() {
    this.boost = false
  }
}
