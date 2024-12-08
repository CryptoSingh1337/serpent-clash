import { createRouter, createWebHistory } from "vue-router"
import HomeView from "@/views/HomeView.vue"
import MainView from "@/views/MainView.vue"
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
      component: MainView
    },
    {
      path: "/:pathMatch(.*)*",
      name: "ErrorView",
      component: ErrorView
    }
  ]
})

export default router
