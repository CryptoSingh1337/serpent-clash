import { CustomStats } from "@/classes/custom_stats.ts"
import { ref, type Ref } from "vue"
import { DisplayDriver } from "@/drivers/display_driver.ts"
import type { Coordinate } from "@/utils/types"

export class StatsDriver {
  displayDriver: DisplayDriver
  stats: Ref<CustomStats> | null
  _stats: CustomStats
  debugMode: boolean = false

  constructor(displayDriver: DisplayDriver) {
    this.displayDriver = displayDriver
    this.debugMode = import.meta.env.VITE_DEBUG_MODE === "true"
    if (this.debugMode) {
      this.stats = ref<CustomStats>(new CustomStats())
      this._stats = this.stats.value
    } else {
      this.stats = null
      this._stats = new CustomStats()
    }
  }

  renderStats(): void {
    this.displayDriver.renderStats(this._stats)
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
