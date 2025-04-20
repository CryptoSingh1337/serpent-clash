<script setup lang="ts">
import { ref, onMounted } from "vue"

// Simulated performance metrics
const fps = ref(60)
const renderTime = ref(5)
const networkLatency = ref(50)
const memoryTrend = ref([50, 52, 55, 53, 58, 60, 62, 65, 63, 67])

// Update FPS randomly to simulate changes
onMounted(() => {
  setInterval(() => {
    fps.value = Math.floor(55 + Math.random() * 10)
    renderTime.value = Math.floor(3 + Math.random() * 5)
    networkLatency.value = Math.floor(40 + Math.random() * 30)

    // Update memory trend
    memoryTrend.value.shift()
    memoryTrend.value.push(Math.floor(50 + Math.random() * 20))
  }, 2000)
})
</script>

<template>
  <div class="performance-metrics">
    <div class="metric-panel mb-6">
      <h4 class="text-xl font-bold mb-4">Real-time Performance</h4>
      <div class="grid grid-cols-3 gap-4 sm:grid-cols-1 md:grid-cols-3">
        <div class="metric-item text-center">
          <div class="text-sm text-gray-400 mb-1">FPS</div>
          <div
            class="text-2xl font-bold"
            :class="fps > 55 ? 'text-green-500' : 'text-yellow-500'"
          >
            {{ fps }}
          </div>
        </div>
        <div class="metric-item text-center">
          <div class="text-sm text-gray-400 mb-1">Render Time</div>
          <div
            class="text-2xl font-bold"
            :class="renderTime < 5 ? 'text-green-500' : 'text-yellow-500'"
          >
            {{ renderTime }} ms
          </div>
        </div>
        <div class="metric-item text-center">
          <div class="text-sm text-gray-400 mb-1">Network Latency</div>
          <div
            class="text-2xl font-bold"
            :class="networkLatency < 60 ? 'text-green-500' : 'text-yellow-500'"
          >
            {{ networkLatency }} ms
          </div>
        </div>
      </div>
    </div>

    <div class="metric-panel">
      <h4 class="text-xl font-bold mb-4">Memory Usage Trend</h4>
      <div class="memory-chart flex items-end justify-between">
        <div
          v-for="(value, index) in memoryTrend"
          :key="index"
          class="memory-bar bg-blue-500 rounded-t"
          :style="{ height: `${value}%` }"
        ></div>
      </div>
      <div class="text-xs text-gray-400 mt-2 text-center">Last 10 minutes</div>
    </div>
  </div>
</template>

<style scoped>
.performance-metrics {
  width: 100%;
}

.metric-panel {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.metric-item {
  background-color: #222;
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid #333;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.memory-chart {
  height: 160px;
  gap: 4px;
}

.memory-bar {
  transition: height 0.5s ease;
  flex: 1;
  min-width: 20px;
  max-width: 40px;
}
</style>
