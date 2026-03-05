<template>
  <div class="distribution-page">
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
          <span class="menu-icon">🏠</span>
          <span class="menu-label">首页</span>
        </RouterLink>
        <RouterLink class="active" :to="`/projects/${route.params.id}`">
          <span class="menu-icon">📂</span>
          <span class="menu-label">项目管理</span>
        </RouterLink>
        <a href="#" @click.prevent>
          <span class="menu-icon">📊</span>
          <span class="menu-label">查询统计</span>
        </a>
        <a href="#" @click.prevent>
          <span class="menu-icon">✅</span>
          <span class="menu-label">质控管理</span>
        </a>
        <a href="#" @click.prevent>
          <span class="menu-icon">⚙️</span>
          <span class="menu-label">系统设置</span>
        </a>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h1>项目管理-分发影像数据</h1>
        <div class="user-box">👤 {{ user?.username || '项目经理' }}</div>
      </header>

      <section class="stage-tabs">
        <button v-for="step in steps" :key="step.key" :class="['stage-tab', { active: activeStep === step.key }]" @click="activeStep = step.key">
          <span>{{ step.icon }}</span>
          <span>{{ step.label }}</span>
        </button>
      </section>

      <section class="workspace">
        <div class="task-panel panel">
          <header>
            <h3>任务列表</h3>
            <div class="actions">
              <button class="ghost" @click="batchAssign">批量分发</button>
              <button class="ghost" @click="smartAssign">智能分发</button>
            </div>
          </header>

          <div class="list-wrap">
            <article
              v-for="item in taskList"
              :key="item.id"
              :class="['task-item', { selected: selectedTaskId === item.id }]"
              @click="selectedTaskId = item.id"
            >
              <div class="task-head">
                <h4>{{ item.caseNo }}</h4>
                <span>分配给: {{ item.assignee || '未分配' }}</span>
              </div>
              <p>患者：{{ item.patient }} | 年龄：{{ item.age }}岁 | 检查类型：{{ item.type }}</p>
              <footer>
                <span>🕘 {{ item.date }}</span>
                <span>🖼 {{ item.images }}张影像</span>
                <button @click.stop="selectedTaskId = item.id">查看详情</button>
              </footer>
            </article>
          </div>
        </div>

        <div class="detail-panel panel">
          <h3>任务详情</h3>
          <div v-if="activeTask" class="detail-content">
            <h2>{{ activeTask.caseNo }}</h2>
            <p>受试者：{{ activeTask.patient }}（{{ activeTask.age }}岁）</p>
            <p>检查类型：{{ activeTask.type }}</p>
            <p>分配状态：{{ activeTask.assignee || '未分配' }}</p>
            <p>影像数量：{{ activeTask.images }} 张</p>
          </div>
          <div v-else class="placeholder">请选择左侧任务查看详情</div>

          <button class="full-image-btn" @click="showPreview = true">🔍 查看完整影像</button>
        </div>
      </section>
    </main>

    <div v-if="showPreview && activeTask" class="modal-mask" @click.self="showPreview = false">
      <section class="modal-card" role="dialog" aria-modal="true">
        <h3>完整影像预览</h3>
        <p>{{ activeTask.caseNo }} · {{ activeTask.type }} · {{ activeTask.images }} 张</p>
        <div class="mock-image">完整影像显示区域</div>
        <div class="modal-actions">
          <button class="primary" @click="showPreview = false">关闭</button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { getUser } from '../api/user';
import { useUserStore } from '../store/user';
import type { User } from '../types/user';

type StepKey = 'screening' | 'inspection' | 'reading' | 'quality';
type TaskItem = {
  id: number;
  caseNo: string;
  patient: string;
  age: number;
  type: string;
  date: string;
  images: number;
  assignee: string;
};

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);
const activeStep = ref<StepKey>('screening');
const selectedTaskId = ref<number>(1);
const showPreview = ref(false);

const steps: Array<{ key: StepKey; icon: string; label: string }> = [
  { key: 'screening', icon: '🧑‍⚕️', label: '受试者筛选阶段' },
  { key: 'inspection', icon: '🖼', label: '影像数据检查阶段' },
  { key: 'reading', icon: '👁', label: '阅片阶段' },
  { key: 'quality', icon: '☑', label: '质量抽查' },
];

const taskList = ref<TaskItem[]>([
  { id: 1, caseNo: '病例#XXXXXXXXXX', patient: '张', age: 68, type: 'OTC', date: '2026-02-28', images: 16, assignee: '' },
  { id: 2, caseNo: '病例#XXXXXXXXXY', patient: '李', age: 72, type: 'OTC', date: '2026-02-28', images: 14, assignee: '' },
  { id: 3, caseNo: '病例#XXXXXXXXXZ', patient: '王', age: 65, type: 'OTC', date: '2026-02-28', images: 18, assignee: '' },
]);

const activeTask = computed(() => taskList.value.find((item) => item.id === selectedTaskId.value) || null);

function batchAssign() {
  taskList.value = taskList.value.map((task, index) => ({ ...task, assignee: task.assignee || `技师${index + 1}` }));
}

function smartAssign() {
  taskList.value = taskList.value.map((task) => ({ ...task, assignee: task.assignee || 'AI建议技师A' }));
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
.distribution-page { min-height: 100vh; display: flex; background: #eef2f6; color: #1f2937; }
.sidebar { width: 250px; background: #3f8fdb; color: #fff; padding: 22px 16px; }
.brand { display: flex; align-items: center; gap: 10px; border-bottom: 1px solid rgba(255,255,255,.25); padding-bottom: 14px; }
.logo { width: 48px; height: 48px; border-radius: 8px; background: #fff; color: #3f8fdb; display: grid; place-items: center; font-size: 24px; }
.brand h2 { margin: 0; font-size: 28px; line-height: 1.15; font-weight: 700; }
.brand p { margin: 3px 0 0; font-size: 16px; line-height: 1.2; opacity: .9; }
.menu { margin-top: 18px; display: grid; gap: 12px; }
.menu a { color: #ecf5ff; text-decoration: none; padding: 10px 12px; border-radius: 8px; display: flex; align-items: center; gap: 8px; }
.menu a.active, .menu a:hover { background: rgba(255,255,255,.2); }
.content { flex: 1; padding: 18px 20px; }
.topbar { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; }
.stage-tabs { margin-top: 12px; background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 10px; display: flex; gap: 10px; flex-wrap: wrap; }
.stage-tab { border: none; background: #f1f5f9; border-radius: 8px; padding: 8px 12px; display: flex; gap: 8px; cursor: pointer; }
.stage-tab.active { background: #dbeafe; color: #1d4ed8; }
.workspace { margin-top: 12px; display: grid; grid-template-columns: 2fr 1fr; gap: 12px; }
.panel { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 12px; }
.task-panel > header { display: flex; align-items: center; justify-content: space-between; }
.actions { display: flex; gap: 8px; }
.ghost { border: 1px solid #cbd5e1; background: #f8fafc; border-radius: 8px; padding: 6px 10px; cursor: pointer; }
.list-wrap { margin-top: 8px; display: grid; gap: 10px; }
.task-item { border: 1px solid #e2e8f0; border-radius: 8px; padding: 10px; cursor: pointer; }
.task-item.selected { border-color: #93c5fd; background: #f8fbff; }
.task-head { display: flex; justify-content: space-between; }
.task-head h4 { margin: 0; }
.task-item p { margin: 8px 0; color: #64748b; }
.task-item footer { display: flex; gap: 10px; align-items: center; color: #64748b; }
.task-item footer button { margin-left: auto; border: none; background: transparent; color: #2563eb; cursor: pointer; }
.detail-panel { display: flex; flex-direction: column; }
.detail-content { flex: 1; display: grid; gap: 8px; }
.placeholder { flex: 1; display: grid; place-items: center; color: #94a3b8; }
.full-image-btn { margin-top: auto; border: none; border-radius: 8px; background: #3f8fdb; color: #fff; padding: 12px; cursor: pointer; }
.modal-mask { position: fixed; inset: 0; background: rgba(15, 23, 42, 0.35); display: grid; place-items: center; }
.modal-card { width: min(620px, calc(100vw - 32px)); background: #fff; border-radius: 12px; border: 1px solid #dbe3ef; padding: 16px; }
.mock-image { margin-top: 12px; border: 1px dashed #94a3b8; border-radius: 10px; height: 220px; display: grid; place-items: center; color: #64748b; }
.modal-actions { margin-top: 12px; display: flex; justify-content: flex-end; }
.primary { border: none; border-radius: 8px; background: #3f8fdb; color: #fff; padding: 8px 14px; cursor: pointer; }

@media (max-width: 1100px) {
  .distribution-page { flex-direction: column; }
  .sidebar { width: auto; }
  .workspace { grid-template-columns: 1fr; }
}
</style>
