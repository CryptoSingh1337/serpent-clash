import type { Player } from '@/classes/entity'

export type Payload = {
  eventType: string,
  body: string
}

export type Position = {
  x: number,
  y: number
}

export type BackendPlayer = {
  id: string,
  positions: Position[]
  direction: number
}

export type Players = {
  [id: string]: Player
}
