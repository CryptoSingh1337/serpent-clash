<script setup lang="ts">
import {onMounted, ref, type Ref, useTemplateRef} from "vue"
import type {QuadTree} from "@/utils/types"
import { Constants as constants } from "@/utils/constants.ts"

const canvasRef = useTemplateRef<HTMLCanvasElement>("canvas-ref")
let ctx: CanvasRenderingContext2D | null = null
const devicePixelRatio = window.devicePixelRatio || 1
let quadTree: QuadTree | null = null
const canvasWidth = 800;
const canvasHeight = 800;
const scale = canvasWidth / (constants.worldBoundary.maxX - constants.worldBoundary.minX);
let fpsInterval = 0
const fps = 60
let then: number

function worldToCanvas(worldX: number, worldY: number): {x: number, y: number} {
  return {
    x: (worldX - constants.worldBoundary.minX) * scale,
    y: (worldY - constants.worldBoundary.minY) * scale
  }
}

function renderQuadTree(ctx: CanvasRenderingContext2D): void {
  const qt = quadTree
  if (!qt) {
    return
  }
  const nodes: QuadTree[] = []
  nodes.push(qt)
  let nodeId = 1
  while (nodes.length > 0) {
    const node = nodes[0]
    const center = worldToCanvas(node.boundary.x, node.boundary.y)
    const w = node.boundary.w * scale
    const h = node.boundary.h * scale
    const boundary = {x: center.x - w, y: center.y - h, w: w*2, h: h*2}
    node.points.forEach(p => {
      const c = worldToCanvas(p.x, p.y)
      ctx.beginPath()
      ctx.arc(c.x, c.y, 1, 0, 2 * Math.PI, true)
      if (p.pointType == "head") {
        ctx.fillStyle = "rgba(0, 255, 0, 1)"
      } else {
        ctx.fillStyle = "rgba(0, 145, 255, 1)"
      }
      ctx.fill()
    })
    ctx.strokeStyle = "rgba(255, 255, 255, 1)"
    ctx.lineWidth = 0.5
    ctx.rect(boundary.x, boundary.y, boundary.w, boundary.h)
    ctx.stroke()
    if (node.divided) {
      ctx.beginPath()
      ctx.moveTo(boundary.x + boundary.w, boundary.y)
      ctx.lineTo(boundary.x + boundary.w, boundary.y + boundary.h * 2)
      ctx.moveTo(boundary.x, boundary.y + boundary.h)
      ctx.lineTo(boundary.x + boundary.w * 2, boundary.y + boundary.h)
      ctx.stroke()
      nodes.push(node.nw)
      nodes.push(node.ne)
      nodes.push(node.sw)
      nodes.push(node.se)
    }
    nodes.shift()
    nodeId++
  }
}

function render(): void {
  requestAnimationFrame(() => {
    render()
  })
  const now = Date.now();
  const elapsed = now - then;
  if (elapsed > fpsInterval) {
    then = now - (elapsed % fpsInterval);
    const canvas = canvasRef.value
    if (!canvas) {
      return
    }
    if (!ctx) {
      return
    }
    ctx.fillStyle = "#000"
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    renderQuadTree(ctx)
  }
}

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas) {
    return
  }
  canvas.width = canvasWidth * devicePixelRatio
  canvas.height = canvasHeight * devicePixelRatio
  ctx = canvas.getContext("2d", { alpha: false })
  if (!ctx) {
    throw new Error("not able to get context")
  }
  setInterval(async () => {
    const response = await fetch("/metrics")
    if (response.ok) {
      const body = await response.json()
      if (body.quadTree) {
        quadTree = body.quadTree as QuadTree
      }
    }
  }, 2000)
  ctx.scale(devicePixelRatio, devicePixelRatio)
  fpsInterval = 1000 / fps;
  then = Date.now();
  render()
})
</script>

<template>
  <div class="w-full h-full flex justify-center flex-col">
    <h3 class="text-center p-5 font-bold text-4xl">Quad Tree - Visualization</h3>
    <canvas ref="canvas-ref" class="m-auto"></canvas>
  </div>
</template>
