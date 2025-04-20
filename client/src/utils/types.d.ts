import type { Player } from "@/classes/entity"

export type ReconcileEvent = {
  seq: number
  event: { coordinate: Coordinate; boost: boolean }
}

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

export type CameraCoordinates = {
  x: number
  y: number
  width: number
  height: number
}

export type QuadTree = {
  boundary: { x: number; y: number; w: number; h: number }
  points: { x: number; y: number; entityId: number; pointType: string }[]
  divided: boolean
  nw: QuadTree
  ne: QuadTree
  sw: QuadTree
  se: QuadTree
}

export type SpawnRegions = {
  radius: number
  regions: Coordinate[]
}

export type ServerMetrics = {
  cpuUsage: number
  memoryUsageInMB: number
  uptimeInSec: number
  bytesSent: number
  bytesReceived: number
  playerCount: number
  foodCount: number
}
