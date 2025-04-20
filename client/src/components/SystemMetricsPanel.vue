<script setup lang="ts">
import type { ServerMetrics } from "@/utils/types"

defineProps<{
  serverMetrics: ServerMetrics
}>()

// Format uptime seconds to HH:MM:SS
function formatUptime(seconds: number): string {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)

  return `${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`
}

// Convert bytes to MB
function bytesToMB(bytes: number): number {
  return bytes / (1024 * 1024)
}
</script>

<template>
  <div class="mt-8 bg-gray-800 rounded-lg p-6 border border-gray-700 shadow-lg">
    <h5 class="text-xl font-bold mb-4 text-blue-400 flex items-center">
      <i class="bi bi-pc-display mr-2"></i>System Information
    </h5>
    <div class="grid grid-cols-2 gap-4 sm:grid-cols-1 md:grid-cols-2">
      <div
        class="bg-gray-700 rounded-lg p-3 border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
      >
        <div class="font-medium text-gray-300 flex items-center">
          <i class="bi bi-clock-history mr-2 text-green-400"></i>Server Uptime
        </div>
        <div class="text-lg font-semibold text-green-500">
          {{ formatUptime(serverMetrics.uptimeInSec) }}
        </div>
      </div>
      <div
        class="bg-gray-700 rounded-lg p-3 border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
      >
        <div class="font-medium text-gray-300 flex items-center">
          <i class="bi bi-cpu mr-2 text-red-400"></i>CPU Usage
        </div>
        <div class="text-lg font-semibold text-red-500">
          {{ serverMetrics.cpuUsage }}%
        </div>
      </div>
      <div
        class="bg-gray-700 rounded-lg p-3 border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
      >
        <div class="font-medium text-gray-300 flex items-center">
          <i class="bi bi-reception-4 mr-2 text-blue-400"></i>Total Network
          Traffic in MB
        </div>
        <div class="text-lg font-semibold text-sky-400">
          <i class="bi bi-arrow-down-short text-red-400"></i
          >{{ bytesToMB(serverMetrics.bytesReceived).toFixed(2) }} MB/s
          <i class="bi bi-arrow-up-short text-blue-300"></i
          >{{ bytesToMB(serverMetrics.bytesSent).toFixed(2) }} MB/s
        </div>
      </div>
      <div
        class="bg-gray-700 rounded-lg p-3 border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
      >
        <div class="font-medium text-gray-300 flex items-center">
          <i class="bi bi-hdd-network mr-2 text-purple-400"></i>Active
          Connections
        </div>
        <div class="text-lg font-semibold text-purple-500">
          {{ serverMetrics.playerCount }}
        </div>
      </div>
    </div>
  </div>
</template>
