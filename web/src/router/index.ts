import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Info from '../views/Info.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Dashboard from '../views/Dashboard.vue'
import GradeView from '@/views/GradeView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/info',
    name: 'Info', 
    component: Info
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard/grades',
    name: 'GradeView',
    component: GradeView,
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard/attendance',
    name: 'GradeView',
    component: GradeView,
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard/instruments',
    name: 'GradeView',
    component: GradeView,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router