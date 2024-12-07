import type { Camera } from "@/classes/camera.ts"

export class DisplayDriver {
  ctx: CanvasRenderingContext2D
  camera: Camera

  constructor(ctx: CanvasRenderingContext2D, camera: Camera) {
    this.ctx = ctx
    this.camera = camera
  }
}
