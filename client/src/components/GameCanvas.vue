<script setup lang="ts">
import {onBeforeUnmount, onMounted, ref, useTemplateRef} from "vue"
import {Application, Container, Graphics, FederatedPointerEvent, Point } from "pixi.js"

const gameCanvas = useTemplateRef<HTMLDivElement>("game-canvas")
const pointer = ref<{x: number, y: number}>({x: 0, y: 0})
let app: Application
let world: Container

function drawHex(x: number, y: number, hexSize: number): Graphics {
  const points: Point[] = []
  for (let i = 0; i < 6; i++) {
    const angle = (Math.PI / 3) * i
    const hx = x + hexSize * Math.cos(angle)
    const hy = y + hexSize * Math.sin(angle)
    points.push(new Point(hx, hy))
    // if (i === 0) {
    //   hex.moveTo(hx, hy)
    // } else {
    //   hex.lineTo(hx, hy)
    // }
  }
  return new Graphics()
      .poly(points)
      .stroke({color: 0x808080, width: 1, alpha: 0.75})
}

onMounted(async () => {
  app = new Application()
  await app.init({
    preference: "webgpu",
    resizeTo: window,
    background: "#191825",
    antialias: true
  })
  world = new Container()
  world.position.set(app.screen.width / 2, app.screen.height / 2)
  console.log("Renderer", app.renderer)
  app.stage.addChild(world)

  console.log(world.getBounds())
  const hexSize = 50
  const hexHeight = hexSize * 2
  const hexWidth = Math.sqrt(25) * hexSize
  const verticalSpacing = (hexHeight * 8) / 4
  const horizontalSpacing = hexWidth / 1.5
  for (let x = -2750; x < 2750; x += verticalSpacing) {
    for (let y = -2750; y < 2750; y += horizontalSpacing) {
      world.addChild(drawHex(x, y, 50))
    }
  }
  // const circle = new Graphics().circle(0, 0, 50).fill({color: 0xffffff})
  // world.addChild(circle)
  app.stage.eventMode = 'static'
  app.stage.hitArea = app.screen
  app.stage.addEventListener('pointermove', (e: FederatedPointerEvent) => {
    pointer.value.x = e.clientX
    pointer.value.y = e.clientY
  })
  if (gameCanvas.value) {
    gameCanvas.value.appendChild(app.canvas)
  }
})

onBeforeUnmount(() => {
  if (app) {
    app.destroy(true, true)
  }
})
</script>

<template>
  <div class="w-full h-full">
    <div ref="game-canvas"></div>
    <div class="absolute top-2.5 right-2.5">
      <span>X: {{ pointer.x }}, </span>
      <span>Y: {{ pointer.y }}</span>
    </div>
  </div>
</template>
