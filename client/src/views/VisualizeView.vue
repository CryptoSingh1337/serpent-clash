<script setup lang="ts">
import {onMounted, ref, type Ref, useTemplateRef} from "vue"
import type {QuadTree, SpawnRegions} from "@/utils/types"
import { Constants as constants } from "@/utils/constants.ts"

const canvasRef = useTemplateRef<HTMLCanvasElement>("canvas-ref")
let ctx: CanvasRenderingContext2D | null = null
let zoom = window.devicePixelRatio || 1
let panX = 0
let panY = 0
// let isDragging = false
// let dragStartX = 0
// let dragStartY = 0
// let dragStartPanX = panX
// let dragStartPanY = panY
let quadTree: QuadTree | null = null
let spawnRegions: SpawnRegions | null = null
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
  if (!quadTree) {
    return
  }
  const nodes: QuadTree[] = []
  nodes.push(quadTree)
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
      ctx.closePath()
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
      ctx.closePath()
      nodes.push(node.nw)
      nodes.push(node.ne)
      nodes.push(node.sw)
      nodes.push(node.se)
    }
    nodes.shift()
    nodeId++
  }
}

function renderSpawnRegions(ctx: CanvasRenderingContext2D): void {
  ctx.fillStyle = "rgba(255, 255, 255, 0.4)"
  ctx.strokeStyle = "rgba(255, 255, 0, 1)"
  ctx.lineWidth = 0.5
  if (spawnRegions && spawnRegions.regions) {
    const radius = spawnRegions.radius * scale
    spawnRegions.regions.forEach((r, idx) => {
      const c = worldToCanvas(r.x, r.y)
      ctx.beginPath()
      ctx.arc(c.x, c.y, radius, 0, 2 * Math.PI, true)
      ctx.textAlign = "center"
      ctx.textBaseline = "middle"
      ctx.stroke()
      ctx.fillText(`region: ${idx}`, c.x, c.y, 100)
      ctx.closePath()
    })
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
    ctx.save();
    ctx.setTransform(zoom, 0, 0, zoom, panX, panY)
    ctx.scale(zoom, zoom);
    ctx.fillStyle = "#000"
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    renderQuadTree(ctx)
    renderSpawnRegions(ctx)
    ctx.restore()
  }
}

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas) {
    return
  }
  canvas.width = canvasWidth * devicePixelRatio
  canvas.height = canvasHeight * devicePixelRatio
  // canvas.addEventListener("mousedown", (e) => {
  //   isDragging = true
  //   dragStartX = e.clientX
  //   dragStartY = e.clientY
  //   dragStartPanX = panX
  //   dragStartPanY = panY
  // })
  //
  // canvas.addEventListener("mousemove", (e) => {
  //   if (!isDragging) return
  //
  //   const dx = e.clientX - dragStartX
  //   const dy = e.clientY - dragStartY
  //
  //   panX = dragStartPanX + dx
  //   panY = dragStartPanY + dy
  // })
  //
  // canvas.addEventListener("mouseup", () => {
  //   isDragging = false
  // })
  //
  // canvas.addEventListener("mouseleave", () => {
  //   isDragging = false
  // })
  //
  // canvas.addEventListener("wheel", (e) => {
  //   e.preventDefault()
  //   const zoomFactor = 1.1
  //   if (e.deltaY < 0) {
  //     zoom *= zoomFactor
  //   } else {
  //     zoom /= zoomFactor
  //   }
  //   zoom = Math.max(1, zoom)
  //   console.log("Zoom", zoom)
  // })
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
      if (body.spawnRegions) {
        spawnRegions = body.spawnRegions as SpawnRegions
      }
    }
  }, 2000)
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
