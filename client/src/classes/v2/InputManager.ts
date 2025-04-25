import type {Game} from "@/classes/v2/Game.ts"

export class InputManager {
  game: Game
  mousePosition: {x: number, y: number}
  boost: boolean

  constructor(game: Game) {
    this.game = game
    this.mousePosition = {x: 0, y: 0}
    this.boost = false
    window.addEventListener("mousemove", this.onMouseMove.bind(this))
    window.addEventListener("mouseup", this.onMouseUp.bind(this))
    window.addEventListener("mousedown", this.onMouseDown.bind(this))
    window.addEventListener("mouseleave", this.onMouseLeave.bind(this))
  }

  onMouseMove(e: MouseEvent) {
    if (!this.game.div) {
      return
    }
    const rect = this.game.div.getBoundingClientRect()
    this.mousePosition.x = e.clientX - rect.left
    this.mousePosition.y = e.clientY - rect.top
  }

  onMouseUp() {
    this.boost = false
  }

  onMouseDown() {
    this.boost = true
  }

  onMouseLeave() {
    this.boost = false
  }
}
