<template>
  <div class="review-page">
    <aside class="sidebar">
      <div class="brand">
        <div class="logo">👁</div>
        <div>
          <h2>SJTURC</h2>
          <p>眼科工作管理系统</p>
        </div>
      </div>

      <nav class="menu">
        <RouterLink :to="`/user/${user?.id || route.params.id}`">
          <span>🏠 首页</span>
        </RouterLink>
        <RouterLink class="active" :to="`/projects/${route.params.id}`">
          <span>📂 项目管理</span>
        </RouterLink>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h1>项目管理-阅片审核</h1>
        <div>👤 {{ user?.username || '项目经理' }}</div>
      </header>

      <section class="panel">
        <h3>审核任务</h3>
        <article v-for="item in reviewTasks" :key="item.id" :class="['task', { active: activeTaskId === item.id }]" @click="activeTaskId = item.id">
          <div class="row">
            <strong>{{ item.caseNo }}</strong>
            <span>{{ item.status }}</span>
          </div>
          <p>{{ item.summary }}</p>
        </article>
      </section>

      <section class="panel detail" v-if="activeTask">
        <h3>审核详情</h3>
        <p><strong>病例：</strong>{{ activeTask.caseNo }}</p>
        <p><strong>初审意见：</strong>{{ activeTask.summary }}</p>
        <label>
          审核结论
          <textarea v-model.trim="opinion" placeholder="请输入阅片审核意见"></textarea>
        </label>
        <div class="actions">
          <button class="reject" @click="decision('驳回')">驳回</button>
          <button class="pass" @click="decision('通过')">审核通过</button>
        </div>
        <p v-if="message" class="message">{{ message }}</p>
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

type ReviewTask = { id: number; caseNo: string; status: string; summary: string };

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);
const activeTaskId = ref(1);
const opinion = ref('');
const message = ref('');

const reviewTasks = ref<ReviewTask[]>([
  { id: 1, caseNo: '病例#REVIEW-001', status: '待审核', summary: '黄斑中心厚度偏高，建议重点复核。' },
  { id: 2, caseNo: '病例#REVIEW-002', status: '待审核', summary: '层次结构清晰，建议快速通过。' },
]);

const activeTask = computed(() => reviewTasks.value.find((item) => item.id === activeTaskId.value) || null);

function decision(result: '驳回' | '通过') {
  if (!activeTask.value) return;
  message.value = `${activeTask.value.caseNo} 已${result}，时间：${new Date().toLocaleString('zh-CN')}`;
}

onMounted(async () => {
  if (!userStore.token) {
    router.push('/login');
    return;
  }

  try {
    const res = await getUser(route.params.id as string, userStore.token);
    user.value = res.data.data;
  } catch {
    user.value = null;
  }
});
</script>

<style scoped>
.review-page { min-height: 100vh; display: flex; background: #eef2f6; }
.sidebar { width: 250px; background: #3f8fdb; color: #fff; padding: 22px 16px; }
.brand { display: flex; align-items: center; gap: 10px; border-bottom: 1px solid rgba(255,255,255,.25); padding-bottom: 14px; }
.logo { width: 48px; height: 48px; border-radius: 8px; background: #fff; color: #3f8fdb; display: grid; place-items: center; font-size: 24px; }
.menu { margin-top: 18px; display: grid; gap: 12px; }
.menu a { color: #ecf5ff; text-decoration: none; padding: 10px 12px; border-radius: 8px; }
.menu a.active, .menu a:hover { background: rgba(255,255,255,.2); }
.content { flex: 1; padding: 18px 20px; display: grid; gap: 12px; }
.topbar, .panel { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px; }
.task { border: 1px solid #dbe3ef; border-radius: 8px; padding: 10px; cursor: pointer; margin-top: 10px; }
.task.active { border-color: #93c5fd; background: #f8fbff; }
.row { display: flex; justify-content: space-between; }
textarea { margin-top: 6px; width: 100%; min-height: 100px; border: 1px solid #cbd5e1; border-radius: 8px; padding: 8px; }
.actions { margin-top: 10px; display: flex; gap: 8px; }
button { border: none; border-radius: 8px; padding: 8px 12px; cursor: pointer; }
.reject { background: #fee2e2; color: #991b1b; }
.pass { background: #dcfce7; color: #166534; }
.message { margin-top: 8px; color: #334155; }
</style>
