import type { Player } from "@/classes/entity"
import type { Camera } from "@/classes/camera.ts"
import type { Stats } from "@/classes/stats.ts"

export type Coordinate = {
  x: number
  y: number
}

export type BackendPlayer = {
  id: string
  color: string
  positions: Coordinate[]
  seq: number
}

export type Players = {
  [id: string]: Player
}

export type DebugMenuEventPayload = {
  player: Player | null
  camera: Camera
  stats: Stats
  ctx: CanvasRenderingContext2D
}
