import Vue from 'vue'
import Router from 'vue-router'
import LoginView from '@/views/LoginView'
import DeviceView from '@/views/DeviceView'

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/device',
      name: 'device',
      component: DeviceView,
      meta: {
        requiresAuth: true
      }
    }
  ]
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const token = localStorage.getItem('token')
  if (requiresAuth && !token) {
    next('/')
  } else {
    next()
  }
})

export default router
