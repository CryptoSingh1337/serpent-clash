<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue"
import type { GameMetrics, ServerMetrics } from "@/utils/types"
import ServerMetricsPanel from "@/components/ServerMetricsPanel.vue"
import TabNavigation from "@/components/TabNavigation.vue"
import SystemMetricsPanel from "@/components/SystemMetricsPanel.vue"
import GameMetricsPanel from "@/components/GameMetricsPanel.vue"
import { useRouter } from "vue-router"

const router = useRouter()

const serverMetrics = ref<ServerMetrics>({
  cpuUsage: 0,
  memoryUsage: 0,
  heapAllocated: 0,
  heapReserved: 0,
  totalHeapAllocated: 0,
  heapObjects: 0,
  lastGCMs: 0,
  gcPauseMicro: 0,
  numGoroutines: 0,
  uptimeInSec: 0,
  bytesSent: 0,
  bytesReceived: 0,
  packetsSent: 0,
  packetsReceived: 0,
  errorIn: 0,
  errorOut: 0,
  dropIn: 0,
  dropOut: 0,
  activeConnections: 0
})
const gameMetrics = ref<GameMetrics>({
  playerCount: 0,
  systemUpdateTimeInLastTick: 0,
  maxSystemUpdateTime: 0,
  systemUpdateTimeInLastTenTicks: [],
  noOfCollisionsInLastTenTicks: []
})

let sse: EventSource | null = null
const activeTab = ref("server-metrics")
const tabs = [
  { id: "server-metrics", label: "Server Metrics", icon: "bi bi-graph-up" },
  { id: "game-metrics", label: "Game metrics", icon: "bi bi-speedometer2" }
]

function handleTabChange(tabId: string) {
  activeTab.value = tabId
}

onMounted(() => {
  sse = new EventSource("/metrics/subscribe/info")
  sse.onmessage = function (event: MessageEvent<string>) {
    const info = JSON.parse(event.data) as {
      serverMetrics: ServerMetrics
      gameMetrics: GameMetrics
    }
    serverMetrics.value = info.serverMetrics
    gameMetrics.value = info.gameMetrics
  }
})

onBeforeUnmount(() => {
  if (sse) {
    sse.close()
  }
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
        <i class="bi bi-clipboard-data mr-2 text-blue-300"></i>
        Server Dashboard
      </p>
      <TabNavigation
        :tabs="tabs"
        :initial-active-tab="activeTab"
        @tab-change="handleTabChange"
      />
      <hr class="my-1" />
      <span
        class="flex items-center px-6 py-3 font-medium focus:outline-none text-left text-gray-400 hover:text-gray-300 border-l-4 border-transparent"
      >
        <i class="bi bi-cookie mr-2 text-xl"></i>
        <button class="text-sm" @click="router.push('/world')">
          World View
        </button>
      </span>
    </div>
    <div
      class="flex flex-1 flex-col items-center bg-gradient-to-br from-gray-900 to-gray-800"
    >
      <div
        v-if="activeTab === 'server-metrics'"
        class="tab-content w-full animate-fadeIn"
      >
        <div
          class="bg-gray-800 rounded-lg shadow-xl border border-gray-700 p-6"
        >
          <h4 class="text-2xl font-bold mb-6 text-center text-blue-400">
            <i class="bi bi-graph-up mr-2"></i>
            Server Metrics
          </h4>
          <ServerMetricsPanel :server-metrics="serverMetrics" />
          <SystemMetricsPanel :server-metrics="serverMetrics" />
        </div>
      </div>
      <div
        v-if="activeTab === 'game-metrics'"
        class="tab-content w-full animate-fadeIn"
      >
        <div
          class="bg-gray-800 rounded-lg shadow-xl border border-gray-700 p-6"
        >
          <h4 class="text-2xl font-bold mb-6 text-center text-blue-400">
            <i class="bi bi-speedometer2 mr-2"></i>
            Game Metrics
          </h4>
          <GameMetricsPanel :game-metrics="gameMetrics" />
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
