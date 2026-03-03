<template>
  <div class="project-page">
    <aside class="sidebar">
      <div class="brand">
        <div class="logo">👁</div>
        <div>
          <h2>SJTURC</h2>
          <p>眼科工作管理系统</p>
        </div>
      </div>

      <nav class="menu">
        <RouterLink :to="`/user/${user?.id || route.params.id}`">🏠 首页</RouterLink>
        <a class="active" href="#" @click.prevent>📂 项目管理</a>
        <a href="#" @click.prevent>📊 查询统计</a>
        <a href="#" @click.prevent>✅ 质控管理</a>
        <a href="#" @click.prevent>⚙️ 系统设置</a>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h1>项目管理</h1>
        <div class="user-box">
          <span>👤 {{ user?.username || '项目经理' }}</span>
          <button @click="onLogout">退出登录</button>
        </div>
      </header>

      <section class="notice">
        <p>今天是 {{ today }}，您有 <strong>{{ todoItems.length }}</strong> 项待处理任务</p>
        <button>+ 新建任务</button>
      </section>

      <section class="main-grid">
        <div class="panel project-panel">
          <header>
            <h3>最近项目</h3>
          </header>
          <article v-for="project in recentProjects" :key="project.name" class="project-card">
            <div class="project-head">
              <h4>{{ project.name }}</h4>
              <span :class="['state', project.stateClass]">{{ project.state }}</span>
            </div>
            <p>{{ project.desc }}</p>
            <footer>
              <span>{{ project.date }}</span>
              <span>{{ project.members }}人参与</span>
              <button>查看详情</button>
              <button class="primary">处理任务</button>
            </footer>
          </article>
        </div>

        <div class="panel todo-panel">
          <header>
            <h3>待处理</h3>
          </header>
          <ul>
            <li v-for="item in todoItems" :key="item">{{ item }}</li>
          </ul>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { getUser } from '../api/user';
import { useUserStore } from '../store/user';
import type { User } from '../types/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);

const recentProjects = [
  {
    name: '青光眼早期诊断研究',
    state: '进行中',
    stateClass: 'running',
    desc: '研究青光眼早期诊断的新方法和技术',
    date: '2023-06-15',
    members: 8,
  },
  {
    name: '视网膜病变图像分析',
    state: '待开始',
    stateClass: 'pending',
    desc: '研究青光眼早期诊断的新方法和技术',
    date: '2023-07-01',
    members: 12,
  },
];

const todoItems = ['认证', '分发影像数据', '阅片审核'];
const today = computed(() => new Date().toLocaleDateString('zh-CN'));

onMounted(async () => {
  if (!userStore.token) {
    router.push('/login');
    return;
  }

  try {
    const res = await getUser(route.params.id as string, userStore.token);
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
  } catch {
    user.value = null;
  }
});

function onLogout() {
  userStore.logout();
  router.push('/login');
}
</script>

<style scoped>
.project-page { min-height: 100vh; display: flex; background: #eef2f6; color: #1f2937; }
.sidebar { width: 250px; background: #3f8fdb; color: #fff; padding: 22px 16px; }
.brand { display: flex; align-items: center; gap: 10px; border-bottom: 1px solid rgba(255,255,255,.25); padding-bottom: 14px; }
.logo { width: 48px; height: 48px; border-radius: 8px; background: #fff; color: #3f8fdb; display: grid; place-items: center; font-size: 24px; }
.brand h2 { margin: 0; }
.brand p { margin: 4px 0 0; opacity: .9; }
.menu { margin-top: 18px; display: grid; gap: 12px; }
.menu a { color: #ecf5ff; text-decoration: none; padding: 10px 12px; border-radius: 8px; }
.menu a.active, .menu a:hover { background: rgba(255,255,255,.2); }
.content { flex: 1; padding: 18px 20px; }
.topbar { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; align-items: center; }
.topbar h1 { margin: 0; }
.user-box { display: flex; align-items: center; gap: 10px; }
.user-box button { border: 1px solid #d1d5db; background: #fff; border-radius: 6px; padding: 6px 10px; cursor: pointer; }
.notice { margin-top: 14px; background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; align-items: center; }
.notice strong { color: #ef4444; }
.notice button { background: #3f8fdb; color: #fff; border: none; border-radius: 8px; padding: 10px 20px; }
.main-grid { margin-top: 14px; display: grid; grid-template-columns: 2fr 1fr; gap: 14px; }
.panel { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; overflow: hidden; }
.panel > header { padding: 14px 16px; border-bottom: 1px solid #e8edf5; background: #f8fafc; }
.panel h3 { margin: 0; color: #4f77d9; }
.project-card { margin: 14px; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px; }
.project-head { display: flex; justify-content: space-between; align-items: center; }
.project-head h4 { margin: 0; font-size: 30px; }
.project-card p { margin: 10px 0; color: #4b5563; }
.state { padding: 4px 14px; border-radius: 999px; color: #fff; font-size: 14px; }
.state.running { background: #3b82f6; }
.state.pending { background: #f59e0b; }
.project-card footer { display: grid; grid-template-columns: auto auto 1fr auto; gap: 12px; align-items: center; color: #6b7280; }
.project-card footer button { justify-self: end; border: 1px solid #d1d5db; background: #fff; border-radius: 8px; padding: 8px 12px; cursor: pointer; }
.project-card footer button.primary { background: #dbe7ff; color: #1f3b8f; border-color: #c8d7ff; margin-left: 8px; }
.todo-panel ul { list-style: none; margin: 0; padding: 12px; display: grid; gap: 10px; }
.todo-panel li { border: 1px solid #dfe5f1; border-radius: 8px; padding: 16px; font-size: 26px; }

@media (max-width: 1100px) {
  .project-page { flex-direction: column; }
  .sidebar { width: auto; }
  .main-grid { grid-template-columns: 1fr; }
  .project-card footer { grid-template-columns: 1fr; }
  .project-card footer button { justify-self: start; }
}
</style>
