import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import GameView from '@/views/GameView.vue'
import HistoryView from '@/views/HistoryView.vue'
import LeaderboardView from '@/views/LeaderboardView.vue'

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/game', name: 'game', component: GameView },
  { path: '/history', name: 'history', component: HistoryView },
  { path: '/leaderboard', name: 'leaderboard', component: LeaderboardView }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router

