import type { Player } from "@/classes/Player.ts"
import type { Food } from "@/classes/Food.ts"

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

export type BackendFood = {
  coordinate: Coordinate
}

export type Players = {
  [id: string]: Player
}

export type Foods = {
  [id: number]: Food
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
  memoryUsage: number
  heapAllocated: number
  heapReserved: number
  totalHeapAllocated: number
  heapObjects: number
  lastGCMs: number
  gcPauseMicro: number
  numGoroutines: number
  uptimeInSec: number
  bytesSent: number
  bytesReceived: number
  packetsSent: number
  packetsReceived: number
  errorIn: number
  errorOut: number
  dropIn: number
  dropOut: number
  activeConnections: number
}

export type GameMetrics = {
  playerCount: number
  systemUpdateTimeInLastTick: number
  maxSystemUpdateTime: number
  systemUpdateTimeInLastTenTicks: number[]
  noOfCollisionsInLastTenTicks: number[]
}
