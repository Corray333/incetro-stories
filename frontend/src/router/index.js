import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/home',
      name: 'home',
      component: HomeView, 
      meta: { requiresAuth: true }
    },
    { path: '/login', name:'login', component: () => import('../views/LoginView.vue') },
    { path: '/profile', name:'profile', component: () => import('../views/ProfileView.vue'),meta: { requiresAuth: true } },
  ]
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    console.log(document.cookie)
     if (!document.cookie) {
       next({ name: 'login' });
     } else {
       next();
     }
  } else {
     next();
  }
 });

export default router