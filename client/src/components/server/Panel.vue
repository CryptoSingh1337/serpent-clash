<script setup lang="ts">
defineProps<{
  label: string
  value: number
  icon: string
  color: string
  suffix?: string
}>()
</script>

<template>
  <div
      class="bg-gray-700 rounded-lg p-4 text-center min-w-[160px] flex-1 max-w-[240px] shadow-md border border-gray-600 hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1"
  >
    <div class="text-xl font-bold mb-2 text-gray-200">
      <i :class="[icon, color]" class="mr-2"></i>
      {{ label }}
    </div>
    <div :class="color" class="text-3xl font-bold">
      {{ value.toFixed(0) }} {{ suffix ?? "" }}
    </div>
  </div>
</template>
