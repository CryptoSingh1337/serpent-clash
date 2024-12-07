import { Stats } from "@/classes/stats.ts"
import { ref, type Ref } from "vue"

export class StatsDriver {
  stats: Ref<Stats>

  constructor() {
    this.stats = ref<Stats>(new Stats())
  }

  renderStats(c: CanvasRenderingContext2D): void {
    this.stats.value.renderStats(c)
  }

  calculateFps(): void {
    this.stats.value.calculateFps()
  }

  updatePing(ping: number): void {
    this.stats.value.updatePing(ping)
  }

  updatePlayerId(playerId: string): void {
    this.stats.value.updatePlayerId(playerId)
  }

  updateStatus(status: string): void {
    this.stats.value.updateStatus(status)
  }

  updateMouseCoordinate(x: number, y: number): void {
    this.stats.value.updateMouseCoordinate(x, y)
  }

  updateHeadCoordinate(x: number, y: number): void {
    this.stats.value.updateHeadCoordinate(x, y)
  }

  updateCameraCoordinate(x: number, y: number): void {
    this.stats.value.updateCameraCoordinate(x, y)
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this.stats.value.updateCameraWidthAndHeight(width, height)
  }

  getPingCooldown(): number {
    return this.stats.value.pingCooldown
  }

  reducePingCooldown(): void {
    this.stats.value.pingCooldown -= 1
  }

  resetPingCooldown(): void {
    this.stats.value.resetPingCooldown()
  }

  reset(): void {
    this.stats.value.reset()
  }
}
