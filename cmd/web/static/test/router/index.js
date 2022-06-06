import {
  createRouter,
  createWebHistory
} from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
  }
]

const router = createRouter({
  linkActiveClass: 'active',
  history: createWebHistory(process.env.BASE_URL),
  routes: routes,
  scrollBehavior () {
    return {
      top: 0
    }
  }
})

export default router
