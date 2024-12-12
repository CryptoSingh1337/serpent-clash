import { Constants } from "@/utils/constants.ts"

export function lerpAngle(a: number, b: number, t: number): number {
  let diff = b - a
  // Handle wrapping from -π to π
  while (diff < -Math.PI) {
    diff += 2 * Math.PI
  }
  while (diff > Math.PI) {
    diff -= 2 * Math.PI
  }
  return a + diff * Math.min(t, 1.0)
}

export function lerp(start: number, end: number, t: number): number {
  return start * (1 - t) + end * t
}

export function clamp(value: number, min: number, max: number): number {
  if (value < min) {
    return min
  } else if (value > max) {
    return max
  }
  return value
}

export function roundOff(value: number): number {
  return Math.floor(value * 100) / 100
}

export function getServerBaseUrl(ws: boolean): string {
  let protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
  if (!ws) {
    protocol = window.location.protocol
  }
  return window.location.hostname.includes("localhost")
    ? `${protocol}//${window.location.hostname}:${Constants.serverPort}`
    : `${protocol}//${window.location.hostname}`
}
