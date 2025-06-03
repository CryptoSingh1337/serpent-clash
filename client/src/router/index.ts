import { createRouter, createWebHistory } from "vue-router"
import HomeView from "@/views/HomeView.vue"
import MainView from "@/views/MainView.vue"
import GameView from "@/views/GameView.vue"
import DashboardView from "@/views/DashboardView.vue"
import WorldView from "@/views/WorldView.vue"
import ErrorView from "@/views/ErrorView.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView
    },
    {
      path: "/menu",
      name: "menu",
      component: MainView
    },
    {
      path: "/play",
      name: "play",
      component: GameView
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: DashboardView
    },
    {
      path: "/world",
      name: "world",
      component: WorldView
    },
    {
      path: "/:pathMatch(.*)*",
      name: "ErrorView",
      component: ErrorView
    }
  ]
})

export default router
