<script setup lang="ts">
defineProps<{
  playerCount: number
  foodCount: number
  memoryUsage: number
}>()
</script>

<template>
  <div class="metrics-panel">
    <div class="flex flex-wrap justify-center gap-4 mb-4">
      <div class="metric-card">
        <div class="text-xl font-bold mb-2">Players</div>
        <div class="text-3xl text-green-500">{{ playerCount }}</div>
      </div>

      <div class="metric-card">
        <div class="text-xl font-bold mb-2">Food Items</div>
        <div class="text-3xl text-blue-500">{{ foodCount }}</div>
      </div>

      <div class="metric-card">
        <div class="text-xl font-bold mb-2">Memory Usage</div>
        <div class="text-3xl text-yellow-500">{{ memoryUsage }} MiB</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.metrics-panel {
  width: 100%;
}

.metric-card {
  background-color: #222;
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
  min-width: 160px;
  flex: 1;
  max-width: 240px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #333;
}
</style>
