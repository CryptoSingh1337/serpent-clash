import type { Player } from "@/classes/entity"

export type Coordinate = {
  x: number
  y: number
}

export type BackendPlayer = {
  id: string
  color: string
  positions: Coordinate[]
}

export type Players = {
  [id: string]: Player
}
