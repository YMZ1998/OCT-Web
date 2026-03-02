import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import UserInfo from '../views/UserInfo.vue';

const routes: Array<RouteRecordRaw> = [
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/user/:id', component: UserInfo },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
