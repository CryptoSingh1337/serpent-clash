import { createRouter, createWebHistory } from "vue-router"
import HomeView from "@/views/HomeView.vue"
import PlayView from "@/views/PlayView.vue"
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
      path: "/play",
      name: "play",
      component: PlayView
    },
    {
      path: "/:pathMatch(.*)*",
      name: "ErrorView",
      component: ErrorView
    }
  ]
})

export default router
