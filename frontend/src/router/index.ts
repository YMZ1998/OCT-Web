import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import UserInfo from '../views/UserInfo.vue';
import ProjectManagement from '../views/ProjectManagement.vue';
import HardwareCertification from '../views/HardwareCertification.vue';

const routes: Array<RouteRecordRaw> = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/user/:id', component: UserInfo },
  { path: '/projects/:id', component: ProjectManagement },
  { path: '/projects/:id/hardware', component: HardwareCertification },
  { path: '/:pathMatch(.*)*', redirect: '/login' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
