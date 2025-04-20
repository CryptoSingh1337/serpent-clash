<script setup lang="ts">
import { ref } from "vue"

const props = defineProps<{
  tabs: Array<{
    id: string
    label: string
    icon: string
  }>
  initialActiveTab?: string
}>()

const emit = defineEmits<{
  (e: "tab-change", tabId: string): void
}>()

const activeTab = ref(
  props.initialActiveTab || (props.tabs.length > 0 ? props.tabs[0].id : "")
)

function setActiveTab(tabId: string) {
  activeTab.value = tabId
  emit("tab-change", tabId)
}
</script>

<template>
  <div class="tab-navigation">
    <div class="flex flex-col border-r border-gray-700">
      <span
        v-for="tab in tabs"
        :key="tab.id"
        class="flex items-center px-6 py-3 font-medium focus:outline-none text-left"
        :class="[
          activeTab === tab.id
            ? 'border-l-4 border-blue-500 text-blue-500 bg-gray-800'
            : 'text-gray-400 hover:text-gray-300 border-l-4 border-transparent'
        ]"
      >
        <i :class="tab.icon" class="mr-2 text-xl"></i>
        <button class="text-sm" @click="setActiveTab(tab.id)">
          {{ tab.label }}
        </button>
      </span>
    </div>
  </div>
</template>
