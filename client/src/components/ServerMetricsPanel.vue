<script setup lang="ts">
import ServerPanel from "@/components/server/Panel.vue"
import { bytesToMB } from "@/utils/helper.ts"
import type { ServerMetrics } from "@/utils/types"

defineProps<{
  serverMetrics: ServerMetrics
}>()
</script>

<template>
  <div class="w-full">
    <div class="flex flex-wrap justify-center gap-4 mb-4">
      <ServerPanel label="CPU Usage" :value="serverMetrics.cpuUsage" suffix="%" icon="bi bi-pc-horizontal" color="text-green-500" />
      <ServerPanel label="Memory Usage" :value="bytesToMB(serverMetrics.memoryUsage)" suffix="MB" icon="bi bi-memory" color="text-yellow-400" />
      <ServerPanel label="No. of Goroutines" :value="serverMetrics.numGoroutines" icon="bi bi-cpu" color="text-indigo-400" />
      <ServerPanel label="Heap objects" :value="serverMetrics.heapObjects" icon="bi bi-circle-fill" color="text-fuchsia-400" />
      <ServerPanel label="Last GC run" :value="serverMetrics.lastGCMs" suffix="ms ago" icon="bi bi-circle-fill" color="text-orange-400" />
      <ServerPanel label="Last GC Pause" :value="serverMetrics.gcPauseMicro" suffix="μs ago" icon="bi bi-circle-fill" color="text-rose-400" />

    </div>
  </div>
</template>
