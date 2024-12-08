import { Stats } from "@/classes/stats.ts"
import { ref, type Ref } from "vue"
import { DisplayDriver } from "@/drivers/display_driver.ts"

export class StatsDriver {
  displayDriver: DisplayDriver
  stats: Ref<Stats> | null
  _stats: Stats
  debugMode: boolean = false

  constructor(displayDriver: DisplayDriver) {
    this.displayDriver = displayDriver
    this.debugMode = import.meta.env.VITE_DEBUG_MODE === "true"
    if (this.debugMode) {
      this.stats = ref<Stats>(new Stats())
      this._stats = this.stats.value
    } else {
      this.stats = null
      this._stats = new Stats()
    }
  }

  renderStats(): void {
    this.displayDriver.renderStats(this._stats)
  }

  calculateFps(): void {
    this._stats.calculateFps()
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

  updateMouseCoordinate(x: number, y: number): void {
    this._stats.updateMouseCoordinate(x, y)
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
