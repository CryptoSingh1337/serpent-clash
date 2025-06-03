<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import GamePanel from "@/components/game/Panel.vue"
import type { GameMetrics } from "@/utils/types"

const props = defineProps<{
  gameMetrics: GameMetrics
}>()

const memoryTrend = ref([50, 52, 55, 53, 58, 60, 62, 65, 63, 67])

onMounted(() => {
  setInterval(() => {
    memoryTrend.value.shift()
    memoryTrend.value.push(Math.floor(50 + Math.random() * 20))
  }, 2000)
})

const avgProcessingTime = computed(() => {
  if (
    !props.gameMetrics.systemUpdateTimeInLastTenTicks ||
    props.gameMetrics.systemUpdateTimeInLastTenTicks.length === 0
  ) {
    return 0
  }
  return (
    props.gameMetrics.systemUpdateTimeInLastTenTicks.reduce(
      (accumulator, currentValue) => accumulator + currentValue,
      0
    ) / props.gameMetrics.systemUpdateTimeInLastTenTicks.length
  )
})
</script>

<template>
  <div class="w-full">
    <div
      class="bg-gray-800 rounded-lg p-6 border border-gray-700 shadow-lg mb-6"
    >
      <h4 class="text-xl font-bold mb-4 text-blue-400 flex items-center">
        <i class="bi bi-speedometer2 mr-2"></i>
        Real-time metrics
      </h4>
      <div
        class="grid grid-cols-5 gap-4 sm:grid-cols-2 md:grid-cols-5 justify-center"
      >
        <GamePanel
          :label="'Average Tick processing time'"
          :value="avgProcessingTime"
          :threshold="16660"
          suffix="μs"
          :is-threshold-reverse="false"
        />
        <GamePanel
          label="Last Tick processing time"
          :value="gameMetrics.systemUpdateTimeInLastTick"
          :threshold="16660"
          suffix="μs"
          :is-threshold-reverse="false"
        />
        <GamePanel
          label="Max Tick processing time"
          :value="gameMetrics.maxSystemUpdateTime"
          :threshold="16660"
          suffix="μs"
          :is-threshold-reverse="false"
        />
      </div>
    </div>

    <div class="bg-gray-800 rounded-lg p-6 border border-gray-700 shadow-lg">
      <h4 class="text-xl font-bold mb-4 text-blue-400 flex items-center">
        <i class="bi bi-graph-up mr-2"></i>
        Memory Usage Trend
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
