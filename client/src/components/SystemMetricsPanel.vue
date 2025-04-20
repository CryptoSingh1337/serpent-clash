<script setup lang="ts">
defineProps<{
  playerCount: number
}>()

const serverUptime = Math.floor(Math.random() * 24)
const cpuUsage = Math.floor(Math.random() * 100)
const networkTraffic = Math.floor(Math.random() * 1000)
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
          {{ serverUptime }} hours
        </div>
      </div>
      <div
        class="bg-gray-700 rounded-lg p-3 border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
      >
        <div class="font-medium text-gray-300 flex items-center">
          <i class="bi bi-cpu mr-2 text-red-400"></i>CPU Usage
        </div>
        <div class="text-lg font-semibold text-red-500">{{ cpuUsage }}%</div>
      </div>
      <div
        class="bg-gray-700 rounded-lg p-3 border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
      >
        <div class="font-medium text-gray-300 flex items-center">
          <i class="bi bi-reception-4 mr-2 text-blue-400"></i>Network Traffic
        </div>
        <div class="text-lg font-semibold text-blue-500">
          {{ networkTraffic }} KB/s
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
          {{ playerCount }}
        </div>
      </div>
    </div>
  </div>
</template>
