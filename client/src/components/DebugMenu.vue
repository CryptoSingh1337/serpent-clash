<script setup lang="ts">
import type { DebugDriver } from "@/drivers/debug_driver.ts"

const props = defineProps<{
  debugMenu: DebugDriver | null
}>()
const menuItems = [
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
</script>

<template>
  <div class="mt-2 relative border rounded-s p-2 backdrop-blur text-xs">
    <h1 class="text-center">Debug menu</h1>
    <div class="space-y-1" :key="item.id" v-for="item in menuItems">
      <h1>{{ item.title }}</h1>
      <div :class="item.css">
        <div :key="subField.id" v-for="subField in item.subFields">
          <label v-if="subField.tag === 'input'" :for="subField.id">{{
            subField.label
          }}</label>
          <input
            v-if="subField.tag === 'input'"
            :id="subField.id"
            class="text-black ml-1"
            :type="subField.type"
            :style="subField.style"
          />
          <button
            v-if="item.id === 'teleport' && subField.tag === 'button'"
            class="w-full py-1 px-1 text-center bg-blue-500 hover:bg-blue-700"
            @click.prevent="props.debugMenu && props.debugMenu.teleport(1, 2)"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
