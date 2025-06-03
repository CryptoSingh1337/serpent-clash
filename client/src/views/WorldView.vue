<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue"
import QuadTreeVisualization from "@/components/QuadTreeVisualization.vue"
import type { QuadTree, SpawnRegions } from "@/utils/types"

const quadTree = ref<QuadTree | null>(null)
const spawnRegions = ref<SpawnRegions | null>(null)
const apiError = ref<boolean>(false)
let interval: number

onMounted(() => {
  interval = setInterval(async () => {
    if (!apiError.value) {
      try {
        const response = await fetch("/metrics/quad-tree")
        if (response.ok) {
          const body = await response.json()
          if (body.quadTree) {
            quadTree.value = body.quadTree as QuadTree
          }
          if (body.spawnRegions) {
            spawnRegions.value = body.spawnRegions as SpawnRegions
          }
        }
      } catch (e) {
        apiError.value = true
        throw e
      }
    }
  }, 750)
})

onBeforeUnmount(() => {
  clearInterval(interval)
})
</script>

<template>
  <div class="flex bg-gray-900 text-white min-h-screen">
    <div
      class="flex flex-1 flex-col items-center bg-gradient-to-br from-gray-900 to-gray-800"
    >
      <div
        class="bg-gray-800 rounded-lg shadow-xl border border-gray-700 p-4 m-4"
      >
        <QuadTreeVisualization
          :quad-tree="quadTree"
          :spawn-regions="spawnRegions"
        />
      </div>
    </div>
  </div>
</template>
