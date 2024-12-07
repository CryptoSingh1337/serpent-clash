import { Constants } from "@/utils/constants"
import type { CameraCoordinates } from "@/utils/types"

export class Stats {
  fps: number = 0
  ping: number = 0
  mouseCoordinate: { x: number; y: number } = { x: 0, y: 0 }
  headCoordinate: { x: number; y: number } = {
    x: innerWidth / 2,
    y: innerHeight / 2
  }
  cameraCoordinate: CameraCoordinates = {
    x: 0,
    y: 0,
    width: 0,
    height: 0
  }
  playerId: string = ""
  status: string = "offline"

  // internal
  times: number[] = []
  pingCooldown: number = Constants.pingCooldown

  renderStats(c: CanvasRenderingContext2D): void {
    if (!c) {
      throw new Error("Can't find canvas element")
    }
    c.fillStyle = "White"
    c.font = "normal 12px Arial"
    c.fillText(Math.floor(this.fps) + " fps", 5, 15)
    c.fillText(`Status: ${this.status}`, 5, 30)
    c.fillText(`Ping: ${Math.trunc(this.ping * 100) / 100} ms`, 5, 45)
  }

  calculateFps(): void {
    window.requestAnimationFrame((): void => {
      const now = performance.now()
      while (this.times.length > 0 && this.times[0] <= now - 1000) {
        this.times.shift()
      }
      this.times.push(now)
      this.fps = this.times.length
      this.calculateFps()
    })
  }

  updatePing(ping: number): void {
    this.ping = ping
  }

  updatePlayerId(playerId: string): void {
    this.playerId = playerId
  }

  updateStatus(status: string): void {
    this.status = status
  }

  updateMouseCoordinate(x: number, y: number): void {
    this.mouseCoordinate.x = x
    this.mouseCoordinate.y = y
  }

  updateHeadCoordinate(x: number, y: number): void {
    this.headCoordinate.x = x
    this.headCoordinate.y = y
  }

  updateCameraCoordinate(x: number, y: number): void {
    this.cameraCoordinate.x = x
    this.cameraCoordinate.y = y
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this.cameraCoordinate.width = width
    this.cameraCoordinate.height = height
  }

  resetPingCooldown(): void {
    this.pingCooldown = Constants.pingCooldown
  }

  reset(): void {
    this.status = "offline"
    this.ping = 0
    this.pingCooldown = Constants.pingCooldown
  }
}
