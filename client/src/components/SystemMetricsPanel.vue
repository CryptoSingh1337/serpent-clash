<script setup lang="ts">
defineProps<{
  playerCount: number
}>()

// These would ideally come from the server, but for now we'll generate random values
const serverUptime = Math.floor(Math.random() * 24)
const cpuUsage = Math.floor(Math.random() * 100)
const networkTraffic = Math.floor(Math.random() * 1000)
</script>

<template>
  <div
    class="system-metrics mt-8 bg-gray-800 rounded-lg p-6 border border-gray-700"
  >
    <h5 class="text-xl font-bold mb-4">System Information</h5>
    <div class="grid grid-cols-2 gap-4 sm:grid-cols-1 md:grid-cols-2">
      <div class="metric-item">
        <div class="font-medium text-gray-300">Server Uptime</div>
        <div class="text-lg">{{ serverUptime }} hours</div>
      </div>
      <div class="metric-item">
        <div class="font-medium text-gray-300">CPU Usage</div>
        <div class="text-lg">{{ cpuUsage }}%</div>
      </div>
      <div class="metric-item">
        <div class="font-medium text-gray-300">Network Traffic</div>
        <div class="text-lg">{{ networkTraffic }} KB/s</div>
      </div>
      <div class="metric-item">
        <div class="font-medium text-gray-300">Active Connections</div>
        <div class="text-lg">{{ playerCount * 2 }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.system-metrics {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.metric-item {
  background-color: #222;
  border-radius: 8px;
  padding: 0.75rem;
  border: 1px solid #333;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}
</style>
