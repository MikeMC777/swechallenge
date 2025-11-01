import { createRouter, createWebHistory } from 'vue-router'

const Home = () => import('@/views/Home.vue')
const StockDetail = () => import('@/views/StockDetail.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/stocks/:symbol', component: StockDetail, props: true },
  ],
})

export default router
