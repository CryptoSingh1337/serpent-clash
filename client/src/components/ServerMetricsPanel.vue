<script setup lang="ts">
defineProps<{
  playerCount: number
  foodCount: number
  memoryUsage: number
}>()
</script>

<template>
  <div class="w-full">
    <div class="flex flex-wrap justify-center gap-4 mb-4">
      <div
        class="bg-gray-700 rounded-lg p-4 text-center min-w-[160px] flex-1 max-w-[240px] shadow-md border border-gray-600 hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1"
      >
        <div class="text-xl font-bold mb-2 text-gray-200">
          <i class="bi bi-people-fill mr-2 text-green-400"></i>Players
        </div>
        <div class="text-3xl text-green-500 font-bold">{{ playerCount }}</div>
      </div>

      <div
        class="bg-gray-700 rounded-lg p-4 text-center min-w-[160px] flex-1 max-w-[240px] shadow-md border border-gray-600 hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1"
      >
        <div class="text-xl font-bold mb-2 text-gray-200">
          <i class="bi bi-circle-fill mr-2 text-blue-400"></i>Food Items
        </div>
        <div class="text-3xl text-blue-500 font-bold">{{ foodCount }}</div>
      </div>

      <div
        class="bg-gray-700 rounded-lg p-4 text-center min-w-[160px] flex-1 max-w-[240px] shadow-md border border-gray-600 hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1"
      >
        <div class="text-xl font-bold mb-2 text-gray-200">
          <i class="bi bi-memory mr-2 text-yellow-400"></i>Memory Usage
        </div>
        <div class="text-3xl text-yellow-500 font-bold">
          {{ memoryUsage }} MiB
        </div>
      </div>
    </div>
  </div>
</template>
