<script setup lang="ts">
import { onMounted, ref } from "vue"
import type {
  GameMetrics,
  QuadTree,
  ServerMetrics,
  SpawnRegions
} from "@/utils/types"
import ServerMetricsPanel from "@/components/ServerMetricsPanel.vue"
import QuadTreeVisualization from "@/components/QuadTreeVisualization.vue"
import TabNavigation from "@/components/TabNavigation.vue"
import SystemMetricsPanel from "@/components/SystemMetricsPanel.vue"
import PerformanceMetricsPanel from "@/components/PerformanceMetricsPanel.vue"

const serverMetrics = ref<ServerMetrics>({
  cpuUsage: 0,
  memoryUsageInMB: 0,
  uptimeInSec: 0,
  bytesSent: 0,
  bytesReceived: 0,
  playerCount: 0,
  foodCount: 0
})
const gameMetrics = ref<GameMetrics>({
  systemUpdateTimeInLastTick: 0,
  maxSystemUpdateTime: 0,
  systemUpdateTimeInLastTenTicks: [],
  noOfCollisionsInLastTenTicks: []
})
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
      }
      if (body.spawnRegions) {
        spawnRegions.value = body.spawnRegions as SpawnRegions
      }
      if (body.serverMetrics) {
        serverMetrics.value = body.serverMetrics
        if (serverMetrics.value && quadTree.value) {
          serverMetrics.value.foodCount = countFoodItems(quadTree.value)
        }
      }
      if (body.gameMetrics) {
        gameMetrics.value = body.gameMetrics
      }
    }
  }, 750)
})
</script>

<template>
  <div class="flex bg-gray-900 text-white min-h-screen">
    <div
      class="w-60 min-h-screen bg-gray-800 border-r border-gray-700 shadow-lg"
    >
      <p
        class="flex justify-center text-xl font-bold p-4 pt-6 text-center text-blue-400"
      >
        <i class="bi bi-clipboard-data mr-2 text-blue-300"></i> Server Dashboard
      </p>
      <TabNavigation
        :tabs="tabs"
        :initial-active-tab="activeTab"
        @tab-change="handleTabChange"
      />
    </div>
    <div
      class="flex flex-1 flex-col items-center bg-gradient-to-br from-gray-900 to-gray-800"
    >
      <div
        v-if="activeTab === 'visualization'"
        class="tab-content animate-fadeIn"
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
      <div
        v-if="activeTab === 'metrics'"
        class="tab-content w-full animate-fadeIn"
      >
        <div
          class="bg-gray-800 rounded-lg shadow-xl border border-gray-700 p-6"
        >
          <h4 class="text-2xl font-bold mb-6 text-center text-blue-400">
            <i class="bi bi-graph-up mr-2"></i>Server Metrics
          </h4>
          <ServerMetricsPanel :server-metrics="serverMetrics" />
          <SystemMetricsPanel :server-metrics="serverMetrics" />
        </div>
      </div>
      <div
        v-if="activeTab === 'performance'"
        class="tab-content w-full animate-fadeIn"
      >
        <div
          class="bg-gray-800 rounded-lg shadow-xl border border-gray-700 p-6"
        >
          <h4 class="text-2xl font-bold mb-6 text-center text-blue-400">
            <i class="bi bi-speedometer2 mr-2"></i>Performance Monitoring
          </h4>
          <PerformanceMetricsPanel :game-metrics="gameMetrics" />
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

.animate-fadeIn {
  animation: fadeIn 0.5s ease-in-out;
}
</style>
