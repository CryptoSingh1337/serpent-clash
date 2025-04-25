<script setup lang="ts">
import {onBeforeUnmount, onMounted, ref, useTemplateRef} from "vue"
import {Game} from "@/classes/v2/Game.ts"

const gameCanvas = useTemplateRef<HTMLDivElement>("game-canvas")
const pointer = ref<{x: number, y: number}>({x: 0, y: 0})
let game: Game

onMounted(async () => {
  game = new Game(gameCanvas.value, pointer)
  await game.init()
  game.start()
})

onBeforeUnmount(() => {
  game.stop()
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
