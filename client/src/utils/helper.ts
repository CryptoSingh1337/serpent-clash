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
