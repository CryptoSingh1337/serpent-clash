<script setup lang="ts">
import type { DebugDriver } from "@/drivers/debug_driver.ts"
import { computed, ref } from "vue"
import type { Coordinate } from "@/utils/types"

const props = defineProps<{
  debugMenu: DebugDriver | null
}>()
const stats = computed(
  () => props.debugMenu && props.debugMenu.game.statsDriver.stats
)
const menuItems = [
  {
    id: "info",
    title: "Info",
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
      }
    ]
  },
  {
    id: "teleport",
    title: "Teleport",
    css: "grid grid-cols-[auto_auto_50px] gap-1 items-center",
    subFields: [
      {
        tag: "input",
        id: "teleport-x",
        label: "X:",
        type: "number",
        style: "width: 50px"
      },
      {
        tag: "input",
        id: "teleport-y",
        label: "Y:",
        type: "number",
        style: "width: 50px"
      },
      {
        tag: "button",
        id: "teleport-btn",
        label: "Update"
      }
    ]
  }
]
const teleportX = ref<number>(0)
const teleportY = ref<number>(0)
async function teleport(): Promise<void> {
  if (!props.debugMenu) {
    console.log("debug menu is initialized")
    return
  }
  const coordinate: Coordinate = {
    x: teleportX.value,
    y: teleportY.value
  }
  await props.debugMenu.teleport(coordinate)
}
</script>

<template>
  <div class="mt-2 w-72 border rounded-s p-2 backdrop-blur-sm text-xs">
    <h1 class="text-center">Debug menu</h1>
    <div class="space-y-1 my-1" :key="item.id" v-for="item in menuItems">
      <h1 class="border-b font-bold italic">{{ item.title }}</h1>
      <div :class="item.css">
        <div :key="subField.id" v-for="subField in item.subFields">
          <span
            v-if="subField.tag === 'span' && subField.label === 'Coordinates:'"
            :id="subField.id"
          >
            {{ subField.label }}
            {{ stats && stats.value && stats.value.headCoordinate }}
          </span>
          <span
            v-if="
              subField.tag === 'span' && subField.label === 'Mouse coordinates:'
            "
            :id="subField.id"
          >
            {{ subField.label }}
            {{ stats && stats.value && stats.value.mouseCoordinate }}
          </span>
          <span
            v-if="
              subField.tag === 'span' &&
              subField.label === 'Camera coordinates:'
            "
            :id="subField.id"
          >
            {{ subField.label }}
            {{ stats && stats.value && stats.value.cameraCoordinate.x }},
            {{ stats && stats.value && stats.value.cameraCoordinate.y }}
          </span>
          <span
            v-if="
              subField.tag === 'span' && subField.label === 'Camera dimensions:'
            "
            :id="subField.id"
          >
            {{ subField.label }}
            {{ stats && stats.value && stats.value.cameraCoordinate.width }},
            {{ stats && stats.value && stats.value.cameraCoordinate.height }}
          </span>
          <span
            v-if="subField.tag === 'span' && subField.label === 'Player id:'"
            :id="subField.id"
          >
            {{ subField.label }}
            {{ stats && stats.value && stats.value.playerId }}
          </span>
          <label v-if="subField.tag === 'input'" :for="subField.id">
            {{ subField.label }}
          </label>
          <input
            class="text-black ml-1"
            v-if="subField.tag === 'input' && subField.id === 'teleport-x'"
            v-model="teleportX"
            :id="subField.id"
            :type="subField.type"
            :style="subField.style"
          />
          <input
            class="text-black ml-1"
            v-else-if="subField.tag === 'input' && subField.id === 'teleport-y'"
            v-model="teleportY"
            :id="subField.id"
            :type="subField.type"
            :style="subField.style"
          />
          <input
            class="text-black ml-1"
            v-else-if="subField.tag === 'input'"
            :id="subField.id"
            :type="subField.type"
            :style="subField.style"
          />
          <button
            v-if="item.id === 'teleport' && subField.tag === 'button'"
            class="w-full py-1 px-1 text-center bg-blue-500 hover:bg-blue-700"
            @click.prevent="teleport"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
