import { Constants } from "@/utils/constants"
import type { CameraCoordinates } from "@/utils/types"
import { roundOff } from "@/utils/helper.ts"

export class CustomStats {
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
  reconcileEvents: number = 0

  // internal
  pingCooldown: number = Constants.pingCooldown

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
    this.headCoordinate.x = roundOff(x)
    this.headCoordinate.y = roundOff(y)
  }

  updateCameraCoordinate(x: number, y: number): void {
    this.cameraCoordinate.x = roundOff(x)
    this.cameraCoordinate.y = roundOff(y)
  }

  updateCameraWidthAndHeight(width: number, height: number): void {
    this.cameraCoordinate.width = width
    this.cameraCoordinate.height = height
  }

  updateReconcileEvent(n: number): void {
    this.reconcileEvents = n
  }

  resetPingCooldown(): void {
    this.pingCooldown = Constants.pingCooldown
  }

  reset(): void {
    this.status = "offline"
    this.ping = 0
    this.pingCooldown = Constants.pingCooldown
    this.reconcileEvents = 0
  }
}
