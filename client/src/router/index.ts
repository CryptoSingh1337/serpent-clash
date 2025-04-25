import { createRouter, createWebHistory } from "vue-router"
import HomeView from "@/views/HomeView.vue"
import MainView from "@/views/MainView.vue"
import GameView from "@/views/GameView.vue"
import PixiJsView from "@/views/PixiJsView.vue"
import GameViewV2 from "@/views/GameViewV2.vue"
import DashboardView from "@/views/DashboardView.vue"
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
      path: "/pixi",
      name: "pixi-js",
      component: PixiJsView
    },
    {
      path: "/play/v2",
      name: "play-v2",
      component: GameViewV2
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: DashboardView
    },
    {
      path: "/:pathMatch(.*)*",
      name: "ErrorView",
      component: ErrorView
    }
  ]
})

export default router
