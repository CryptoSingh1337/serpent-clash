<script setup lang="ts">
import { onMounted, ref } from "vue"
import QuadTreeVisualization from "@/components/QuadTreeVisualization.vue"
import type { QuadTree, SpawnRegions } from "@/utils/types"

const quadTree = ref<QuadTree | null>(null)
const spawnRegions = ref<SpawnRegions | null>(null)

onMounted(() => {
  const source = new EventSource("/metrics/subscribe/quad-tree")
  source.onmessage = function (event: MessageEvent<string>) {
    const info = JSON.parse(event.data) as {
      quadTree: QuadTree
      spawnRegions: SpawnRegions
    }
    quadTree.value = info.quadTree
    spawnRegions.value = info.spawnRegions
  }
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
