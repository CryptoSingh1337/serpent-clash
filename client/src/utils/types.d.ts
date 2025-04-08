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
