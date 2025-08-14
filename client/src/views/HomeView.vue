<script setup lang="ts">
import { useRouter } from "vue-router"
import { ref, onMounted, onUnmounted } from "vue"

const router = useRouter()
const snakeBgCanvas = ref<HTMLCanvasElement | null>(null)
const animationFrameId = ref<number | null>(null)

interface SnakeSegment {
  x: number
  y: number
}

interface Snake {
  segments: SnakeSegment[]
  direction: { x: number; y: number }
  speed: number
  color: string
  size: number
}

interface Food {
  x: number
  y: number
  color: string
  size: number
  pulseState: number
  pulseDirection: number
}

const snakes: Snake[] = []
const foods: Food[] = []

const colors = {
  primary: "#192231",
  secondary: "#fff9c4",
  accent1: "#4caf50",
  accent2: "#8bc34a",
  red: "#f44336",
  orange: "#ffc107",
  purple: "#865dff",
  pink: "#e384ff",
  lightPink: "#ffa3fd"
}

const initAnimation = () => {
  const canvas = snakeBgCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext("2d")
  if (!ctx) return
  const resizeCanvas = () => {
    if (canvas) {
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    }
  }

  resizeCanvas()
  window.addEventListener("resize", resizeCanvas)
  for (let i = 0; i < 5; i++) {
    const snakeColors = [
      colors.accent1,
      colors.accent2,
      colors.purple,
      colors.pink,
      colors.lightPink
    ]
    const snake: Snake = {
      segments: [],
      direction: {
        x: Math.random() * 2 - 1,
        y: Math.random() * 2 - 1
      },
      speed: 1 + Math.random() * 1.5,
      color: snakeColors[i % snakeColors.length],
      size: 8 + Math.floor(Math.random() * 6)
    }
    const magnitude = Math.sqrt(snake.direction.x ** 2 + snake.direction.y ** 2)
    snake.direction.x /= magnitude
    snake.direction.y /= magnitude

    const segmentCount = 15 + Math.floor(Math.random() * 10)
    const startX = Math.random() * canvas.width
    const startY = Math.random() * canvas.height
    for (let j = 0; j < segmentCount; j++) {
      snake.segments.push({
        x: startX - j * snake.size * 1.5 * snake.direction.x,
        y: startY - j * snake.size * 1.5 * snake.direction.y
      })
    }
    snakes.push(snake)
  }

  for (let i = 0; i < 15; i++) {
    const foodColors = [colors.secondary, colors.red, colors.orange]
    foods.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      color: foodColors[Math.floor(Math.random() * foodColors.length)],
      size: 4 + Math.floor(Math.random() * 4),
      pulseState: Math.random(),
      pulseDirection: Math.random() > 0.5 ? 0.02 : -0.02
    })
  }

  const animate = () => {
    if (!canvas || !ctx) return
    ctx.fillStyle = "rgba(25, 34, 49, 1.0)"
    ctx.fillRect(0, 0, canvas.width, canvas.height)
    foods.forEach((food) => {
      food.pulseState += food.pulseDirection
      if (food.pulseState > 1 || food.pulseState < 0) {
        food.pulseDirection *= -1
      }

      const currentSize = food.size * (0.8 + food.pulseState * 0.4)
      ctx.save()
      ctx.beginPath()
      ctx.arc(food.x, food.y, currentSize * 2, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(${
        food.color === colors.red
          ? "244, 67, 54"
          : food.color === colors.orange
            ? "255, 193, 7"
            : "255, 249, 196"
      }, 0.1)`
      ctx.fill()
      ctx.closePath()

      ctx.beginPath()
      ctx.arc(food.x, food.y, currentSize, 0, Math.PI * 2)
      ctx.fillStyle = food.color
      ctx.fill()
      ctx.closePath()
      ctx.restore()

      if (Math.random() < 0.01) {
        food.x += (Math.random() * 2 - 1) * 2
        food.y += (Math.random() * 2 - 1) * 2
      }

      if (food.x < 0) food.x = canvas.width
      if (food.x > canvas.width) food.x = 0
      if (food.y < 0) food.y = canvas.height
      if (food.y > canvas.height) food.y = 0
    })

    snakes.forEach((snake) => {
      const head = { ...snake.segments[0] }
      head.x += snake.direction.x * snake.speed
      head.y += snake.direction.y * snake.speed

      if (head.x < 0) head.x = canvas.width
      if (head.x > canvas.width) head.x = 0
      if (head.y < 0) head.y = canvas.height
      if (head.y > canvas.height) head.y = 0

      snake.segments.unshift(head)
      snake.segments.pop()

      if (Math.random() < 0.005) {
        const angle = (Math.random() * Math.PI) / 4 - Math.PI / 8
        const newX =
          snake.direction.x * Math.cos(angle) -
          snake.direction.y * Math.sin(angle)
        const newY =
          snake.direction.x * Math.sin(angle) +
          snake.direction.y * Math.cos(angle)
        snake.direction.x = newX
        snake.direction.y = newY
        const magnitude = Math.sqrt(
          snake.direction.x ** 2 + snake.direction.y ** 2
        )
        snake.direction.x /= magnitude
        snake.direction.y /= magnitude
      }
      ctx.save()

      for (let i = 0; i < snake.segments.length; i++) {
        const segment = snake.segments[i]
        const segmentSize = snake.size * (1 - (i / snake.segments.length) * 0.5)
        ctx.beginPath()
        ctx.arc(segment.x, segment.y, segmentSize, 0, Math.PI * 2)
        const gradient = ctx.createRadialGradient(
          segment.x,
          segment.y,
          0,
          segment.x,
          segment.y,
          segmentSize
        )
        gradient.addColorStop(0, snake.color)
        gradient.addColorStop(1, `${snake.color}80`) // 50% transparency
        ctx.fillStyle = gradient
        ctx.fill()
        ctx.closePath()
      }
      const head2 = snake.segments[0]
      const eyeSize = snake.size / 4
      const eyeOffset = snake.size / 2
      const eyeOffsetX = snake.direction.y * eyeOffset
      const eyeOffsetY = -snake.direction.x * eyeOffset

      // Left eye
      ctx.beginPath()
      ctx.arc(
        head2.x + eyeOffsetX,
        head2.y + eyeOffsetY,
        eyeSize,
        0,
        Math.PI * 2
      )
      ctx.fillStyle = "white"
      ctx.fill()
      ctx.closePath()

      // Right eye
      ctx.beginPath()
      ctx.arc(
        head2.x - eyeOffsetX,
        head2.y - eyeOffsetY,
        eyeSize,
        0,
        Math.PI * 2
      )
      ctx.fillStyle = "white"
      ctx.fill()
      ctx.closePath()

      // Pupils
      ctx.beginPath()
      ctx.arc(
        head2.x + eyeOffsetX + (snake.direction.x * eyeSize) / 2,
        head2.y + eyeOffsetY + (snake.direction.y * eyeSize) / 2,
        eyeSize / 2,
        0,
        Math.PI * 2
      )
      ctx.fillStyle = "black"
      ctx.fill()
      ctx.closePath()

      ctx.beginPath()
      ctx.arc(
        head2.x - eyeOffsetX + (snake.direction.x * eyeSize) / 2,
        head2.y - eyeOffsetY + (snake.direction.y * eyeSize) / 2,
        eyeSize / 2,
        0,
        Math.PI * 2
      )
      ctx.fillStyle = "black"
      ctx.fill()
      ctx.closePath()
      ctx.restore()
    })
    animationFrameId.value = requestAnimationFrame(animate)
  }

  animate()
}

onMounted(() => {
  initAnimation()
})

onUnmounted(() => {
  if (animationFrameId.value) {
    cancelAnimationFrame(animationFrameId.value)
  }
  window.removeEventListener("resize", () => {})
})
</script>

<template>
  <section
    id="hero"
    class="relative min-h-screen overflow-hidden bg-gradient-to-b from-[#192231] to-[#253649] flex flex-col justify-center items-center"
  >
    <canvas
      ref="snakeBgCanvas"
      class="absolute top-0 left-0 w-full h-full"
    ></canvas>
    <div class="absolute top-0 left-0 w-full h-full">
      <div
        class="absolute top-[10%] left-[5%] w-32 h-32 rounded-full bg-gradient-to-br from-[#4caf5020] to-[#8bc34a20] blur-xl"
      ></div>
      <div
        class="absolute bottom-[15%] right-[10%] w-40 h-40 rounded-full bg-gradient-to-tr from-[#865dff20] to-[#e384ff20] blur-xl"
      ></div>
      <div
        class="absolute top-[40%] right-[15%] w-24 h-24 rounded-full bg-gradient-to-r from-[#ffc10720] to-[#f4433620] blur-xl"
      ></div>
    </div>
    <div
      class="relative z-10 text-center mx-auto px-4 py-8 backdrop-blur-sm bg-[#19223180] rounded-2xl border border-[#ffffff20] shadow-2xl max-w-4xl"
    >
      <div class="animate-pulse-slow">
        <img
          src="@/assets/hero.png"
          alt="Serpent Clash Logo"
          class="w-[25rem] mx-auto mb-8 drop-shadow-[0_0_15px_rgba(76,175,80,0.5)]"
        />
      </div>
      <h1 class="text-4xl mb-4">
        <span class="mr-1">🎮</span>
        <span
          class="font-bold text-transparent bg-clip-text font-cherry bg-gradient-to-r from-[#4caf50] to-[#8bc34a]"
          >Multiplayer Snake Battle Arena
        </span>
      </h1>
      <p class="text-xl mb-8 max-w-2xl mx-auto text-[#fff9c4] leading-relaxed">
        The ultimate multiplayer snake battle arena where strategy meets speed.
        Control your serpent with precision, dodge enemies, and strategically
        cut them off as you slither your way to the top of the leaderboard.
      </p>
      <div class="flex justify-center gap-8 mb-6">
        <div
          class="button-container w-40 h-24 transform hover:scale-105 transition-transform"
        >
          <button
            class="text-white font-extrabold text-4xl button-content flex items-center justify-center"
            @click.prevent="router.push('/menu')"
          >
            <span class="relative z-10 font-cherry">Play</span>
          </button>
        </div>
      </div>
      <div class="flex flex-wrap justify-center gap-4 mt-8">
        <div
          class="px-4 py-2 bg-[#4caf5040] border border-[#4caf50] rounded-full text-[#4caf50] font-semibold"
        >
          ⚔️ Multiplayer
        </div>
        <div
          class="px-4 py-2 bg-[#8bc34a40] border border-[#8bc34a] rounded-full text-[#8bc34a] font-semibold"
        >
          🛡️ Authoritative Server
        </div>
        <div
          class="px-4 py-2 bg-[#ffc10740] border border-[#ffc107] rounded-full text-[#ffc107] font-semibold"
        >
          🔄 Server Reconciliation
        </div>
        <div
          class="px-4 py-2 bg-[#865dff40] border border-[#865dff] rounded-full text-[#865dff] font-semibold"
        >
          ✨ Smooth Rendering
        </div>
      </div>
    </div>
  </section>
  <section
    id="description"
    class="py-20 bg-gradient-to-b from-[#253649] to-[#192231] text-center relative overflow-hidden"
  >
    <div class="absolute top-0 left-0 w-full h-full">
      <div
        class="absolute top-[20%] right-[5%] w-48 h-48 rounded-full bg-gradient-to-br from-[#4caf5010] to-[#8bc34a10] blur-xl"
      ></div>
      <div
        class="absolute bottom-[30%] left-[10%] w-56 h-56 rounded-full bg-gradient-to-tr from-[#865dff10] to-[#e384ff10] blur-xl"
      ></div>
    </div>
    <div class="absolute top-[10%] right-[15%] w-32 h-32 opacity-20">
      <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M30,20 Q45,10 60,20 T90,30 Q80,45 90,60 T80,90 Q65,80 50,90 T20,80 Q30,65 20,50 T30,20"
          fill="none"
          stroke="#4caf50"
          stroke-width="5"
          stroke-linecap="round"
        />
      </svg>
    </div>
    <div
      class="absolute bottom-[15%] left-[10%] w-24 h-24 opacity-20 rotate-45"
    >
      <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M30,20 Q45,10 60,20 T90,30 Q80,45 90,60 T80,90 Q65,80 50,90 T20,80 Q30,65 20,50 T30,20"
          fill="none"
          stroke="#865dff"
          stroke-width="5"
          stroke-linecap="round"
        />
      </svg>
    </div>
    <div class="max-w-5xl mx-auto px-6 relative z-10">
      <div class="inline-block mb-8 relative">
        <h2
          class="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-[#4caf50] to-[#8bc34a] mb-2 p-2 font-cherry"
        >
          What is Serpent Clash?
        </h2>
        <div
          class="h-1 w-32 bg-gradient-to-r from-[#4caf50] to-[#8bc34a] rounded-full mx-auto"
        ></div>
      </div>
      <p class="text-xl text-[#fff9c4] mb-12 max-w-3xl mx-auto leading-relaxed">
        Serpent Clash is a
        <span class="text-[#ffc107] font-semibold">fast-paced</span>, real-time
        multiplayer snake game where players compete to
        <span class="text-[#8bc34a] font-semibold">outmaneuver</span> and
        <span class="text-[#4caf50] font-semibold">outgrow</span> their
        opponents. With responsive gameplay using WebSocket communication and
        server reconciliation, you can slither your way to victory in an
        ever-changing battlefield.
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border border-[#4caf5040] shadow-lg transform hover:scale-[1.02] transition-transform group"
        >
          <div
            class="w-16 h-16 mx-auto mb-4 text-[#4caf50] group-hover:scale-110 transition-transform"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#4caf50] mb-4 font-cherry">
            ⚔️ Real-time Multiplayer
          </h3>
          <p class="text-white leading-relaxed">
            Battle against players worldwide using real-time WebSocket
            communication for smooth, interactive gameplay with low latency.
          </p>
        </div>
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border border-[#865dff40] shadow-lg transform hover:scale-[1.02] transition-transform group"
        >
          <div
            class="w-16 h-16 mx-auto mb-4 text-[#865dff] group-hover:scale-110 transition-transform"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#865dff] mb-4 font-cherry">
            🛡️ Authoritative Server
          </h3>
          <p class="text-white leading-relaxed">
            Our backend ensures fairness and consistent gameplay, providing a
            competitive environment for all players with server-side validation.
          </p>
        </div>
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border border-[#ffc10740] shadow-lg transform hover:scale-[1.02] transition-transform group"
        >
          <div
            class="w-16 h-16 mx-auto mb-4 text-[#ffc107] group-hover:scale-110 transition-transform"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#ffc107] mb-4 font-cherry">
            🔄 Server Reconciliation
          </h3>
          <p class="text-white leading-relaxed">
            Experience accurate game state even under lag conditions, with
            advanced server reconciliation techniques for smooth gameplay.
          </p>
        </div>
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border border-[#e384ff40] shadow-lg transform hover:scale-[1.02] transition-transform group"
        >
          <div
            class="w-16 h-16 mx-auto mb-4 text-[#e384ff] group-hover:scale-110 transition-transform"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#e384ff] mb-4 font-cherry">
            ✨ Dynamic Snake Rendering
          </h3>
          <p class="text-white leading-relaxed">
            Snakes respond to real-time mouse input and update fluidly based on
            multiple coordinates for precise control and smooth visuals.
          </p>
        </div>
      </div>
    </div>
  </section>
  <section
    id="play"
    class="py-20 bg-gradient-to-b from-[#192231] to-[#253649] text-center relative overflow-hidden"
  >
    <div class="absolute top-0 left-0 w-full h-full">
      <div
        class="absolute top-[15%] left-[8%] w-40 h-40 rounded-full bg-gradient-to-br from-[#ffc10710] to-[#f4433610] blur-xl"
      ></div>
      <div
        class="absolute bottom-[20%] right-[5%] w-64 h-64 rounded-full bg-gradient-to-tr from-[#4caf5010] to-[#8bc34a10] blur-xl"
      ></div>
    </div>
    <div
      class="absolute top-[25%] left-[15%] w-6 h-6 rounded-full bg-[#ffc10730] blur-sm"
    ></div>
    <div
      class="absolute bottom-[30%] right-[20%] w-8 h-8 rounded-full bg-[#f4433630] blur-sm"
    ></div>
    <div
      class="absolute top-[60%] right-[30%] w-5 h-5 rounded-full bg-[#fff9c430] blur-sm"
    ></div>
    <div class="max-w-5xl mx-auto px-6 relative z-10">
      <div class="inline-block mb-8 relative">
        <h2
          class="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-[#ffc107] to-[#f44336] mb-2 p-2 font-cherry"
        >
          Gameplay Preview
        </h2>
        <div
          class="h-1 w-32 bg-gradient-to-r from-[#ffc107] to-[#f44336] rounded-full mx-auto"
        ></div>
      </div>
      <div class="mb-10">
        <div
          class="relative rounded-2xl overflow-hidden border-4 border-[#ffffff20] shadow-2xl group"
        >
          <div
            class="absolute inset-0 bg-gradient-to-r from-[#4caf5040] via-[#ffc10740] to-[#865dff40] opacity-0 group-hover:opacity-100 transition-opacity duration-500"
          ></div>
          <iframe
            class="mx-auto aspect-video w-full relative z-10"
            src="https://www.youtube.com/embed/qjxnNrzlXaQ?si=b0WaEwt6iWnkQL0z"
            title="Serpent Clash - Demo"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            referrerpolicy="strict-origin-when-cross-origin"
            :allowfullscreen="false"
          ></iframe>
        </div>
      </div>
      <div class="flex justify-center">
        <a
          href="#features"
          class="relative px-8 py-4 overflow-hidden group rounded-full"
        >
          <span
            class="absolute top-0 left-0 w-full h-full bg-gradient-to-r from-[#865dff] to-[#e384ff] opacity-70 group-hover:opacity-100 transition-opacity"
          ></span>
          <span
            class="relative z-10 flex items-center justify-center text-white font-bold text-lg"
          >
            <svg
              class="w-6 h-6 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
              ></path>
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            Explore Future Features
          </span>
        </a>
      </div>
    </div>
  </section>
  <section
    id="features"
    class="py-20 bg-gradient-to-b from-[#253649] to-[#192231] text-center relative overflow-hidden"
  >
    <div class="absolute top-0 left-0 w-full h-full">
      <div
        class="absolute top-[30%] right-[10%] w-48 h-48 rounded-full bg-gradient-to-br from-[#e384ff10] to-[#ffa3fd10] blur-xl"
      ></div>
      <div
        class="absolute bottom-[10%] left-[5%] w-56 h-56 rounded-full bg-gradient-to-tr from-[#ffc10710] to-[#f4433610] blur-xl"
      ></div>
    </div>
    <div class="absolute top-[5%] left-0 w-full h-20 opacity-10">
      <svg viewBox="0 0 1200 100" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M0,50 Q300,0 600,50 T1200,50"
          fill="none"
          stroke="#4caf50"
          stroke-width="8"
          stroke-linecap="round"
          stroke-dasharray="20,20"
        />
      </svg>
    </div>
    <div class="absolute bottom-[5%] left-0 w-full h-20 opacity-10">
      <svg viewBox="0 0 1200 100" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M1200,50 Q900,100 600,50 T0,50"
          fill="none"
          stroke="#865dff"
          stroke-width="8"
          stroke-linecap="round"
          stroke-dasharray="20,20"
        />
      </svg>
    </div>
    <div class="max-w-5xl mx-auto px-6 relative z-10">
      <div class="inline-block mb-8 relative">
        <h2
          class="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-[#865dff] to-[#e384ff] mb-2 p-2 font-cherry"
        >
          Future Enhancements
        </h2>
        <div
          class="h-1 w-32 bg-gradient-to-r from-[#865dff] to-[#e384ff] rounded-full mx-auto"
        ></div>
      </div>
      <p class="text-xl text-[#fff9c4] mb-12 max-w-3xl mx-auto leading-relaxed">
        We're constantly evolving Serpent Clash with exciting new features to
        enhance your gameplay experience:
      </p>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-10">
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border-l-4 border-[#ffc107] shadow-lg transform hover:translate-y-[-8px] transition-transform"
        >
          <div class="w-16 h-16 mx-auto mb-4 text-[#ffc107]">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#ffc107] mb-4 font-cherry">
            🍎 Food Generation & Growth (Implemented)
          </h3>
          <p class="text-white leading-relaxed">
            Dynamic elements where snakes grow by consuming food, adding a new
            layer of strategy and progression to the gameplay experience.
          </p>
        </div>
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border-l-4 border-[#f44336] shadow-lg transform hover:translate-y-[-8px] transition-transform"
        >
          <div class="w-16 h-16 mx-auto mb-4 text-[#f44336]">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#f44336] mb-4 font-cherry">
            🏆 Leaderboard & UI Enhancements
          </h3>
          <p class="text-white leading-relaxed">
            Improved player experience with a more interactive UI and detailed
            leaderboard statistics to showcase top players and their
            achievements.
          </p>
        </div>
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border-l-4 border-[#4caf50] shadow-lg transform hover:translate-y-[-8px] transition-transform"
        >
          <div class="w-16 h-16 mx-auto mb-4 text-[#4caf50]">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#4caf50] mb-4 font-cherry">
            💬 Chat System
          </h3>
          <p class="text-white leading-relaxed">
            Real-time chat functionality using Server-Sent Events and channels,
            allowing players to communicate, strategize, and build community
            during gameplay.
          </p>
        </div>
        <div
          class="bg-[#19223190] backdrop-blur-sm p-8 rounded-2xl border-l-4 border-[#865dff] shadow-lg transform hover:translate-y-[-8px] transition-transform"
        >
          <div class="w-16 h-16 mx-auto mb-4 text-[#865dff]">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-semibold text-[#865dff] mb-4 font-cherry">
            ⚡ Performance Optimizations
          </h3>
          <p class="text-white leading-relaxed">
            Binary data formats for faster network communication, improved
            rendering efficiency, and enhanced server-side processing for a
            smoother gameplay experience.
          </p>
        </div>
      </div>
    </div>
  </section>
  <footer class="py-12 bg-[#192231] text-center relative overflow-hidden">
    <div
      class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-[#4caf50] via-[#ffc107] to-[#865dff] opacity-70"
    ></div>
    <div class="absolute top-0 left-0 w-full h-full">
      <div
        class="absolute bottom-0 left-0 w-full h-32 bg-gradient-to-t from-[#13192380] to-transparent"
      ></div>
    </div>
    <div class="max-w-5xl mx-auto px-6 relative z-10">
      <div class="w-24 h-24 mx-auto mb-6">
        <img
          src="@/assets/hero.png"
          alt="Serpent Clash Logo"
          class="w-full drop-shadow-[0_0_10px_rgba(76,175,80,0.5)]"
        />
      </div>
      <div class="flex flex-wrap justify-center gap-8 mb-8">
        <a
          href="#hero"
          class="text-[#fff9c4] hover:text-white transition-colors"
          >Home</a
        >
        <a
          href="#description"
          class="text-[#fff9c4] hover:text-white transition-colors"
          >About</a
        >
        <a
          href="#play"
          class="text-[#fff9c4] hover:text-white transition-colors"
          >Gameplay</a
        >
        <a
          href="#features"
          class="text-[#fff9c4] hover:text-white transition-colors"
          >Features</a
        >
        <a
          href="https://github.com/CryptoSingh1337/serpent-clash"
          target="_blank"
          class="text-[#fff9c4] hover:text-white transition-colors"
          >GitHub</a
        >
      </div>
      <div class="flex justify-center gap-4 mt-6">
        <a
          href="#hero"
          class="w-10 h-10 flex items-center justify-center rounded-full bg-[#ffffff10] hover:bg-[#ffffff20] transition-colors group"
        >
          <svg
            class="w-5 h-5 text-[#4caf50] group-hover:text-[#8bc34a] transition-colors"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 10l7-7m0 0l7 7m-7-7v18"
            ></path>
          </svg>
        </a>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.font-cherry {
  font-family: "Cherry Bomb One", system-ui;
  font-weight: 400;
  font-style: normal;
}
</style>
