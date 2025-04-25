<script setup lang="ts">
import { computed, ref, onMounted } from "vue"
import type { Coordinate } from "@/utils/types"
import type {DebugManager} from "@/classes/v2/DebugManager.ts"

const props = defineProps<{
  debugManager: DebugManager | null
}>()
const stats = computed(
  () => props.debugManager && props.debugManager.game.statsManager.stats
)

// For draggable functionality
const isDragging = ref(false)
// Position below the connect button in the top-right corner
const position = ref({ x: window.innerWidth - 330, y: 60 })
const dragOffset = ref({ x: 0, y: 0 })

const startDrag = (event: MouseEvent) => {
  isDragging.value = true
  dragOffset.value = {
    x: event.clientX - position.value.x,
    y: event.clientY - position.value.y
  }

  // Add event listeners for drag and end
  document.addEventListener("mousemove", onDrag)
  document.addEventListener("mouseup", endDrag)
}

const onDrag = (event: MouseEvent) => {
  if (isDragging.value) {
    position.value = {
      x: event.clientX - dragOffset.value.x,
      y: event.clientY - dragOffset.value.y
    }
  }
}

const endDrag = () => {
  isDragging.value = false
  document.removeEventListener("mousemove", onDrag)
  document.removeEventListener("mouseup", endDrag)
}

// Update position on window resize
const updatePosition = () => {
  if (!isDragging.value) {
    position.value = { x: window.innerWidth - 340, y: 60 }
  }
}

// Set up and clean up event listeners
onMounted(() => {
  window.addEventListener("resize", updatePosition)

  return () => {
    document.removeEventListener("mousemove", onDrag)
    document.removeEventListener("mouseup", endDrag)
    window.removeEventListener("resize", updatePosition)
  }
})

const menuItems = [
  {
    id: "info",
    title: "Info",
    icon: "ℹ️",
    css: "grid gap-1 items-start",
    subFields: [
      {
        tag: "span",
        id: "coordinates",
        label: "Coordinates:"
      },
      {
        tag: "span",
        id: "mouse-coordinates",
        label: "Mouse coordinates:"
      },
      {
        tag: "span",
        id: "camera-coordinates",
        label: "Camera coordinates:"
      },
      {
        tag: "span",
        id: "camera-dimensions",
        label: "Camera dimensions:"
      },
      {
        tag: "span",
        id: "player-id",
        label: "Player id:"
      },
      {
        tag: "span",
        id: "reconcile-events",
        label: "Events:"
      }
    ]
  },
  {
    id: "teleport",
    title: "Teleport",
    icon: "🚀",
    css: "grid grid-cols-[auto_auto_auto] gap-2 items-center",
    subFields: [
      {
        tag: "input",
        id: "teleport-x",
        label: "X:",
        type: "number",
        style: "width: 60px"
      },
      {
        tag: "input",
        id: "teleport-y",
        label: "Y:",
        type: "number",
        style: "width: 60px"
      },
      {
        tag: "button",
        id: "teleport-btn",
        label: "Teleport"
      }
    ]
  }
]

const teleportX = ref<number>(0)
const teleportY = ref<number>(0)

async function teleport(): Promise<void> {
  if (!props.debugManager) {
    console.log("debug menu is not initialized")
    return
  }
  const coordinate: Coordinate = {
    x: teleportX.value,
    y: teleportY.value
  }
  await props.debugManager.teleport(coordinate)
}
</script>

<template>
  <div
    class="fixed w-80 border border-gray-700 rounded-lg p-3 backdrop-blur-md text-xs bg-gray-900/70 text-gray-200 shadow-lg z-50"
    :style="{ left: position.x + 'px', top: position.y + 'px' }"
    @mousedown="startDrag"
  >
    <h1 class="text-center text-sm font-bold mb-3 text-blue-400 cursor-move">
      Debug Menu
    </h1>

    <div class="space-y-3" :key="item.id" v-for="item in menuItems">
      <!-- Section Header -->
      <div
        class="flex items-center justify-between border-b border-gray-700 pb-1 hover:bg-gray-800/50 rounded px-1"
      >
        <h2 class="font-bold flex items-center">
          <span class="mr-2">{{ item.icon }}</span>
          <span>{{ item.title }}</span>
        </h2>
      </div>

      <!-- Section Content -->
      <div :class="[item.css, 'pl-2 transition-all duration-200']">
        <div :key="subField.id" v-for="subField in item.subFields" class="mb-1">
          <!-- Info Fields -->
          <template v-if="subField.tag === 'span'">
            <div class="flex flex-col">
              <span class="text-blue-300 font-medium">{{
                subField.label
              }}</span>
              <span class="pl-2 text-gray-300" :id="subField.id">
                <template v-if="subField.label === 'Coordinates:'">
                  {{ stats && stats.value && stats.value.headCoordinate }}
                </template>
                <template v-else-if="subField.label === 'Mouse coordinates:'">
                  {{ stats && stats.value && stats.value.mouseCoordinate }}
                </template>
                <template v-else-if="subField.label === 'Camera coordinates:'">
                  {{ stats && stats.value && stats.value.cameraCoordinate.x }},
                  {{ stats && stats.value && stats.value.cameraCoordinate.y }}
                </template>
                <template v-else-if="subField.label === 'Camera dimensions:'">
                  {{
                    stats && stats.value && stats.value.cameraCoordinate.width
                  }},
                  {{
                    stats && stats.value && stats.value.cameraCoordinate.height
                  }}
                </template>
                <template v-else-if="subField.label === 'Player id:'">
                  {{ stats && stats.value && stats.value.playerId }}
                </template>
                <template v-else-if="subField.label === 'Events:'">
                  {{ stats && stats.value && stats.value.reconcileEvents }}
                </template>
              </span>
            </div>
          </template>

          <!-- Input Fields -->
          <template v-else-if="subField.tag === 'input'">
            <label class="text-blue-300 font-medium" :for="subField.id">
              {{ subField.label }}
            </label>
            <input
              v-if="subField.id === 'teleport-x'"
              class="ml-1 bg-gray-800 border border-gray-700 rounded px-2 py-1 text-white"
              v-model="teleportX"
              :id="subField.id"
              :type="subField.type"
              :style="subField.style"
            />
            <input
              v-else-if="subField.id === 'teleport-y'"
              class="ml-1 bg-gray-800 border border-gray-700 rounded px-2 py-1 text-white"
              v-model="teleportY"
              :id="subField.id"
              :type="subField.type"
              :style="subField.style"
            />
            <input
              v-else
              class="ml-1 bg-gray-800 border border-gray-700 rounded px-2 py-1 text-white"
              :id="subField.id"
              :type="subField.type"
              :style="subField.style"
            />
          </template>

          <!-- Button -->
          <template v-else-if="subField.tag === 'button'">
            <button
              class="w-full py-1 px-2 text-center bg-blue-600 hover:bg-blue-700 rounded text-white font-medium transition-colors duration-200"
              @click.prevent="teleport"
            >
              {{ subField.label }}
            </button>
          </template>
        </div>
      </div>
    </div>

    <!-- Drag handle indicator -->
    <div class="absolute top-0 right-0 p-1 text-gray-500 text-xs">
      <span title="Drag to move">⋮⋮</span>
    </div>
  </div>
</template>
