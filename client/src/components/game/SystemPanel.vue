<script setup lang="ts">
import type { SystemMetrics } from "@/utils/types"
defineProps<{
  system: SystemMetrics
}>()

function avg(arr: number[]) {
  if (!arr || arr.length === 0) return 0
  return Math.round(arr.reduce((a, b) => a + b, 0) / arr.length)
}
</script>

<template>
  <div
    class="bg-gradient-to-br from-gray-900 to-gray-800 rounded-xl shadow-lg border border-green-700 p-6 flex flex-col gap-2 hover:scale-[1.02] transition-transform"
  >
    <div class="flex items-center gap-2 mb-2">
      <span
        class="inline-block w-2 h-2 rounded-full bg-blue-400 animate-pulse"
      ></span>
      <span class="font-semibold text-lg text-blue-300">{{ system.name }}</span>
    </div>
    <div class="flex flex-col gap-1 text-sm text-gray-300">
      <div>
        <span class="font-medium text-gray-400">Last Tick:</span>
        <span class="ml-2">{{ system.systemUpdateTimeInLastTick }} μs</span>
      </div>
      <div>
        <span class="font-medium text-gray-400">Max Tick:</span>
        <span class="ml-2">{{ system.maxSystemUpdateTime }} μs</span>
      </div>
      <div>
        <span class="font-medium text-gray-400">Avg (10 Ticks):</span>
        <span class="ml-2"
          >{{ avg(system.systemUpdateTimeInLastTenTicks) }} μs</span
        >
      </div>
    </div>
  </div>
</template>

<style scoped>
.bg-gradient-to-br {
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
}
</style>
