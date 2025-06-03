<script setup lang="ts">
defineProps<{
  label: string
  value: number
  threshold: number
  isThresholdReverse: boolean
  suffix?: string
}>()
</script>

<template>
  <div
    class="bg-gray-700 rounded-lg p-4 text-center border border-gray-600 shadow-md hover:shadow-lg transition-all duration-300"
  >
    <div class="text-sm text-gray-400 mb-1 flex items-center justify-center">
      <i class="bi bi-display mr-1"></i>
      {{ label }}
    </div>
    <div
      class="text-2xl font-bold"
      :class="
        isThresholdReverse
          ? value > threshold
            ? 'text-green-500'
            : 'text-yellow-500'
          : value < threshold
            ? 'text-green-500'
            : 'text-yellow-500'
      "
    >
      {{ value }} {{ suffix ?? "" }}
    </div>
  </div>
</template>
