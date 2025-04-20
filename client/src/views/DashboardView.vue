<script setup lang="ts">
import { onMounted, ref } from "vue"
import type { QuadTree, SpawnRegions } from "@/utils/types"
import ServerMetricsPanel from "@/components/ServerMetricsPanel.vue"
import QuadTreeVisualization from "@/components/QuadTreeVisualization.vue"
import TabNavigation from "@/components/TabNavigation.vue"
import SystemMetricsPanel from "@/components/SystemMetricsPanel.vue"
import PerformanceMetricsPanel from "@/components/PerformanceMetricsPanel.vue"

const playerCount = ref(0)
const memoryUsage = ref(0)
const foodCount = ref(0)
let quadTree = ref<QuadTree | null>(null)
let spawnRegions = ref<SpawnRegions | null>(null)

const activeTab = ref("visualization")
const tabs = [
  { id: "visualization", label: "Visualization", icon: "bi bi-diagram-3" },
  { id: "metrics", label: "Server Metrics", icon: "bi bi-graph-up" },
  { id: "performance", label: "Performance", icon: "bi bi-speedometer2" }
]

function handleTabChange(tabId: string) {
  activeTab.value = tabId
}

function countFoodItems(tree: QuadTree): number {
  if (!tree) {
    return 0
  }
  let count = 0
  const nodes: QuadTree[] = []
  nodes.push(tree)
  while (nodes.length > 0) {
    const node = nodes[0]
    node.points.forEach((p) => {
      if (p.pointType === "food") {
        count++
      }
    })
    if (node.divided) {
      nodes.push(node.nw)
      nodes.push(node.ne)
      nodes.push(node.sw)
      nodes.push(node.se)
    }
    nodes.shift()
  }
  return count
}

onMounted(() => {
  setInterval(async () => {
    const response = await fetch("/metrics")
    if (response.ok) {
      const body = await response.json()
      if (body.quadTree) {
        quadTree.value = body.quadTree as QuadTree
        foodCount.value = countFoodItems(quadTree.value)
      }
      if (body.spawnRegions) {
        spawnRegions.value = body.spawnRegions as SpawnRegions
      }
      if (body.playerCount !== undefined) {
        playerCount.value = body.playerCount
      }
      if (body.memoryUsageInMiB !== undefined) {
        memoryUsage.value = body.memoryUsageInMiB
      }
    }
  }, 750)
})
</script>

<template>
  <div class="flex">
    <div class="w-60 h-screen border-r-1">
      <p class="flex justify-center text-xl font-bold p-4 pt-6 text-center">
        <i class="bi bi-clipboard-data mr-1"></i> Server dashboard
      </p>
      <TabNavigation
        :tabs="tabs"
        :initial-active-tab="activeTab"
        @tab-change="handleTabChange"
      />
    </div>
    <div class="flex flex-1 flex-col items-center p-6">
      <div v-if="activeTab === 'visualization'" class="tab-content">
        <QuadTreeVisualization
          :quad-tree="quadTree"
          :spawn-regions="spawnRegions"
        />
      </div>
      <div v-if="activeTab === 'metrics'" class="tab-content">
        <div class="p-4">
          <h4 class="text-2xl font-bold mb-4 text-center">Server Metrics</h4>
          <ServerMetricsPanel
            :player-count="playerCount"
            :food-count="foodCount"
            :memory-usage="memoryUsage"
          />
          <SystemMetricsPanel :player-count="playerCount" />
        </div>
      </div>
      <div v-if="activeTab === 'performance'" class="tab-content">
        <div class="p-4">
          <h4 class="text-2xl font-bold mb-4 text-center">
            Performance Monitoring
          </h4>
          <PerformanceMetricsPanel />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
