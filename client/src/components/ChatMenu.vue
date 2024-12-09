<script setup lang="ts">
import { ref } from "vue"

const messages = ref<string[]>([])

messages.value.push(
  ...[
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello",
    "hello"
  ]
)

function handleMessageInputFocus(showChat: boolean) {
  const chat = document.getElementById("chat-menu")
  chat
    ? showChat
      ? chat.classList.remove("opacity-0")
      : chat.classList.add("opacity-0")
    : null
  chat
    ? showChat
      ? chat.classList.add("opacity-100")
      : chat.classList.remove("opacity-100")
    : null
}
</script>

<template>
  <div>
    <div
      id="chat-menu"
      class="flex flex-col-reverse p-4 h-60 w-96 backdrop-blur-sm overflow-hidden opacity-0 transition-opacity ease-in-out delay-100"
    >
      <div id="chat-view" class="overflow-y-scroll">
        <div :key="idx" v-for="(message, idx) in messages" class="flex">
          <span class="mr-1">GG:</span>
          <span class="truncate"> {{ message }} {{ idx }} </span>
        </div>
      </div>
    </div>
    <div class="flex h-7 w-96 items-center px-2 backdrop-blur border-x">
      <input
        class="w-full bg-transparent py-1 outline-0 focus:outline-0"
        type="text"
        aria-label="Type your message here..."
        placeholder="Type your message here..."
        @focusin="handleMessageInputFocus(true)"
        @focusout="handleMessageInputFocus(false)"
      />
      <img class="w-6 h-6 ml-1 cursor-pointer" src="@/assets/message.svg" />
    </div>
  </div>
</template>

<style scoped>
#chat-view {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

#chat-vew::-webkit-scrollbar {
  display: none;
}
</style>
