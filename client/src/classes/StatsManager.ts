import { ref, type Ref } from "vue"
import type { Game } from "@/classes/Game.ts"
import { CustomStats } from "@/classes/CustomStats.ts"
import type { Coordinate } from "@/utils/types"

export class StatsManager {
  game: Game
  stats: Ref<CustomStats> | null
  _stats: CustomStats
  debugMode: boolean

  constructor(game: Game) {
    this.game = game
    this.debugMode = import.meta.env.VITE_DEBUG_MODE === "true"
    console.log("Debug mode", this.debugMode)
    if (this.debugMode) {
      this.stats = ref<CustomStats>(new CustomStats())
      this._stats = this.stats.value
    } else {
      this.stats = null
      this._stats = new CustomStats()
    }
  }

  update(): void {
    this.updateCameraCoordinate(
      this.game.displayDriver.camera.position.x,
      this.game.displayDriver.camera.position.y
    )
    if (this.game.player) {
      if (this.game.player.snake && this.game.player.snake.segments.length > 0) {
        const head = this.game.player.snake.segments[0]
        this.updateHeadCoordinate(head.x, head.y)
      }
      this.updateReconcileEvent(this.game.inputManager.inputQueue.length)
      this.updateMouseCoordinate({
        x: this.game.inputManager.mousePosition.x,
        y: this.game.inputManager.mousePosition.y
      })
      const worldMouseCoordinate = this.game.displayDriver.camera.screenToWorld(
        this.game.inputManager.mousePosition.x,
        this.game.inputManager.mousePosition.y
      )
      this.updateWorldMouseCoordinate(
        worldMouseCoordinate.x,
        worldMouseCoordinate.y
      )
    }
  }

  updatePing(ping: number): void {
    this._stats.updatePing(ping)
  }

  updatePlayerId(playerId: string): void {
    this._stats.updatePlayerId(playerId)
  }

  updateStatus(status: string): void {
    this._stats.updateStatus(status)
  }

  updateMouseCoordinate(coordinate: Coordinate): void {
    this._stats.updateMouseCoordinate(coordinate.x, coordinate.y)
  }

  updateWorldMouseCoordinate(x: number, y: number): void {
    this._stats.updateWorldMouseCoordinate(x, y)
  }

  updateHeadCoordinate(x: number, y: number): void {
    this._stats.updateHeadCoordinate(x, y)
  }

  updateCameraCoordinate(x: number, y: number): void {
    this._stats.updateCameraCoordinate(x, y)
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this._stats.updateCameraWidthAndHeight(width, height)
  }

  updateReconcileEvent(n: number): void {
    this._stats.updateReconcileEvent(n)
  }

  getPingCooldown(): number {
    return this._stats.pingCooldown
  }

  reducePingCooldown(): void {
    this._stats.pingCooldown -= 1
  }

  resetPingCooldown(): void {
    this._stats.resetPingCooldown()
  }

  reset(): void {
    this._stats.reset()
  }
}
