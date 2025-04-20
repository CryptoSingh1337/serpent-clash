<script setup lang="ts">
import { ref, onMounted } from "vue"

const fps = ref(60)
const renderTime = ref(5)
const networkLatency = ref(50)
const memoryTrend = ref([50, 52, 55, 53, 58, 60, 62, 65, 63, 67])

onMounted(() => {
  setInterval(() => {
    fps.value = Math.floor(55 + Math.random() * 10)
    renderTime.value = Math.floor(3 + Math.random() * 5)
    networkLatency.value = Math.floor(40 + Math.random() * 30)
    memoryTrend.value.shift()
    memoryTrend.value.push(Math.floor(50 + Math.random() * 20))
  }, 2000)
})
</script>

<template>
  <div class="w-full">
    <div
      class="bg-gray-800 rounded-lg p-6 border border-gray-700 shadow-lg mb-6"
    >
      <h4 class="text-xl font-bold mb-4 text-blue-400 flex items-center">
        <i class="bi bi-speedometer2 mr-2"></i>Real-time Performance
      </h4>
      <div class="grid grid-cols-3 gap-4 sm:grid-cols-1 md:grid-cols-3">
        <div
          class="bg-gray-700 rounded-lg p-4 text-center border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
        >
          <div
            class="text-sm text-gray-400 mb-1 flex items-center justify-center"
          >
            <i class="bi bi-display mr-1"></i>FPS
          </div>
          <div
            class="text-2xl font-bold"
            :class="fps > 55 ? 'text-green-500' : 'text-yellow-500'"
          >
            {{ fps }}
          </div>
        </div>
        <div
          class="bg-gray-700 rounded-lg p-4 text-center border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
        >
          <div
            class="text-sm text-gray-400 mb-1 flex items-center justify-center"
          >
            <i class="bi bi-stopwatch mr-1"></i>Render Time
          </div>
          <div
            class="text-2xl font-bold"
            :class="renderTime < 5 ? 'text-green-500' : 'text-yellow-500'"
          >
            {{ renderTime }} ms
          </div>
        </div>
        <div
          class="bg-gray-700 rounded-lg p-4 text-center border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
        >
          <div
            class="text-sm text-gray-400 mb-1 flex items-center justify-center"
          >
            <i class="bi bi-wifi mr-1"></i>Network Latency
          </div>
          <div
            class="text-2xl font-bold"
            :class="networkLatency < 60 ? 'text-green-500' : 'text-yellow-500'"
          >
            {{ networkLatency }} ms
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-800 rounded-lg p-6 border border-gray-700 shadow-lg">
      <h4 class="text-xl font-bold mb-4 text-blue-400 flex items-center">
        <i class="bi bi-graph-up mr-2"></i>Memory Usage Trend
      </h4>
      <div class="h-40 flex items-end justify-between gap-1">
        <div
          v-for="(value, index) in memoryTrend"
          :key="index"
          class="bg-blue-500 rounded-t flex-1 min-w-5 max-w-10 transition-all duration-500"
          :style="{ height: `${value}%` }"
        ></div>
      </div>
      <div class="text-xs text-gray-400 mt-2 text-center">Last 10 minutes</div>
    </div>
  </div>
</template>
