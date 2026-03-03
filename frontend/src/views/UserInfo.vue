<template>
  <div class="dashboard-page">
    <aside class="sidebar">
      <div class="brand">
        <div class="logo">👁</div>
        <div>
          <h2>SJTURC</h2>
          <p>眼科工作管理系统</p>
        </div>
      </div>

      <nav class="menu">
        <a class="active" href="#">
          <span class="menu-icon">🏠</span>
          <span class="menu-label">首页</span>
        </a>
        <RouterLink :to="`/projects/${user?.id || route.params.id}`">
          <span class="menu-icon">📂</span>
          <span class="menu-label">项目管理</span>
        </RouterLink>
        <a href="#">
          <span class="menu-icon">📊</span>
          <span class="menu-label">查询统计</span>
        </a>
        <a href="#">
          <span class="menu-icon">✅</span>
          <span class="menu-label">质控管理</span>
        </a>
        <a href="#">
          <span class="menu-icon">⚙️</span>
          <span class="menu-label">系统设置</span>
        </a>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h1>首页</h1>
        <div class="user-box">
          <span>👤 {{ user?.username}}</span>
          <button @click="onLogout">退出登录</button>
        </div>
      </header>

      <section class="banner">
        <p>今天是 {{ today }}，您有 <strong>{{ pendingCount }}</strong> 项待处理任务</p>
        <button @click="openTaskModal">+ 新建任务</button>
      </section>

      <section class="section-head">
        <h3>任务概况</h3>
        <a href="#" @click.prevent>查看全部</a>
      </section>
      <section class="task-grid">
        <div v-for="item in taskStats" :key="item.name" class="task-card">
          <div class="row">
            <span>{{ item.icon }} {{ item.name }}</span>
            <em :class="item.statusClass">{{ item.status }}</em>
          </div>
          <div class="row minor">
            <span>已完成 {{ item.done }}/{{ item.total }}</span>
            <span>{{ item.progress }}%</span>
          </div>
          <div class="progress"><span :style="{ width: `${item.progress}%`, background: item.color }"></span></div>
        </div>
      </section>

      <section class="section-head">
        <h3>任务分布</h3>
        <a href="#" @click.prevent>查看详情</a>
      </section>
      <section class="dist-panel">
        <div class="donut"></div>
        <ul>
          <li v-for="item in taskStats" :key="`${item.name}-legend`">
            <span class="dot" :style="{ background: item.color }"></span>{{ item.name }}
          </li>
        </ul>
      </section>

      <section class="section-head">
        <h3>近期任务</h3>
      </section>
      <section class="table-panel">
        <table>
          <thead>
            <tr>
              <th>项目名称</th>
              <th>试验中心</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in recentTasks" :key="item.key">
              <td>{{ item.projectName }}</td>
              <td>{{ projectCenterMap[item.projectId] || '-' }}</td>
              <td>{{ taskTimeText(item.key) }}</td>
              <td>{{ item.taskName }}</td>
            </tr>
            <tr v-if="!recentTasks.length">
              <td colspan="4">暂无近期任务</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="profile-panel">
        <h3>个人资料</h3>
        <div class="profile-form">
          <label>邮箱 <input v-model="form.email" type="email" placeholder="请输入邮箱" /></label>
          <label>性别
            <select v-model="form.gender">
              <option value="">未填写</option>
              <option value="male">男</option>
              <option value="female">女</option>
            </select>
          </label>
          <label>年龄 <input v-model.number="form.age" type="number" min="0" max="200" placeholder="请输入年龄" /></label>
          <button @click="onSave" :disabled="saving">{{ saving ? '保存中...' : '保存到数据库' }}</button>
        </div>
        <p v-if="message" class="message">{{ message }}</p>
      </section>

      <div v-if="showTaskModal" class="modal-mask" @click.self="closeTaskModal">
        <section class="modal-card" role="dialog" aria-modal="true" aria-labelledby="task-title">
          <h3 id="task-title">新建任务</h3>
          <form class="task-form" @submit.prevent="submitTask">
            <label>选择项目
              <select v-model.number="newTask.projectId" required>
                <option :value="0" disabled>请选择项目</option>
                <option v-for="project in projectStore.recentProjects" :key="project.id" :value="project.id">{{ project.name }}</option>
              </select>
            </label>
            <label>任务名称
              <input v-model.trim="newTask.taskName" placeholder="请输入任务名称" required />
            </label>
            <div class="task-actions">
              <button type="button" @click="closeTaskModal">取消</button>
              <button type="submit" class="primary">创建任务</button>
            </div>
          </form>
        </section>
      </div>

      <footer class="page-footer">
        <span>SJTURC眼科在吸纳工作管理系统</span>
        <span>版本1.0　<a href="#" @click.prevent>帮助</a></span>
      </footer>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { getUser, updateUser } from '../api/user';
import { useProjectStore } from '../store/project';
import { useUserStore } from '../store/user';
import type { User } from '../types/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const projectStore = useProjectStore();
const user = ref<User | null>(null);
const saving = ref(false);
const message = ref('');
const form = ref({ email: '', gender: '', age: undefined as number | undefined });
const showTaskModal = ref(false);
const newTask = ref({ projectId: 0, taskName: '' });

const taskStats = computed(() => {
  const todos = projectStore.todoItems;

  const countPending = (keyword: string) =>
    todos.filter((item) => item.taskName.includes(keyword)).length;

  const buildItem = (name: string, total: number, color: string, icon: string, keyword: string) => {
    const pending = countPending(keyword);
    const done = Math.max(total - pending, 0);
    const progress = Math.round((done / total) * 100);
    const status = pending === 0 ? '已完成' : pending === total ? '待处理' : '进行中';
    const statusClass = pending === 0 ? 's-done' : pending === total ? 's-wait' : 's-doing';
    return { name, done, total, progress, color, status, statusClass, icon };
  };

  return [
    buildItem('硬件认证', 6, '#16a34a', '✅', '认证'),
    buildItem('技师认证', 6, '#f59e0b', '📌', '技师'),
    buildItem('证书颁发', 6, '#0284c7', '🏅', '证书'),
    buildItem('阅片分发', 6, '#7c3aed', '👁', '分发'),
    buildItem('阅片审核', 6, '#4f46e5', '🔍', '审核'),
    buildItem('质量抽查', 6, '#ef4444', '⚠️', '抽查'),
  ];
});

const pendingCount = computed(() => projectStore.pendingCount);
const recentTasks = computed(() => projectStore.todoItems.slice(0, 8));
const projectCenterMap = computed(() =>
  projectStore.recentProjects.reduce((acc, item) => {
    acc[item.id] = item.center;
    return acc;
  }, {} as Record<number, string>),
);
const today = computed(() => new Date().toLocaleDateString('zh-CN'));

function taskTimeText(key: string) {
  const [idText] = key.split('-');
  const id = Number(idText);
  const project = projectStore.recentProjects.find((item) => item.id === id);
  return project?.date || today.value;
}

function formatDate(ts: number) {
  return new Date(ts * 1000).toLocaleDateString();
}

function fillFormFromUser() {
  form.value.email = user.value?.email || '';
  form.value.gender = user.value?.gender || '';
  form.value.age = user.value?.age;
}


function openTaskModal() {
  showTaskModal.value = true;
  newTask.value = {
    projectId: projectStore.recentProjects[0]?.id || 0,
    taskName: '',
  };
}

function closeTaskModal() {
  showTaskModal.value = false;
  newTask.value = { projectId: 0, taskName: '' };
}

function submitTask() {
  const project = projectStore.recentProjects.find((item) => item.id === newTask.value.projectId);
  if (!project) return;

  projectStore.addTask(project, newTask.value.taskName);
  closeTaskModal();
}

onMounted(async () => {
  if (!userStore.token) {
    router.push('/login');
    return;
  }

  try {
    const res = await getUser(route.params.id as string, userStore.token);
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
    fillFormFromUser();
  } catch {
    user.value = null;
  }
});

async function onSave() {
  if (!userStore.token || !user.value) return;
  saving.value = true;
  message.value = '';

  try {
    const res = await updateUser(
      user.value.id,
      {
        email: form.value.email,
        gender: form.value.gender,
        age: form.value.age,
      },
      userStore.token,
    );
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
    message.value = '个人信息已更新并同步到数据库';
  } catch {
    message.value = '保存失败，请稍后重试';
  } finally {
    saving.value = false;
  }
}

function onLogout() {
  userStore.logout();
  router.push('/login');
}
</script>

<style scoped>
.dashboard-page { min-height: 100vh; background: #eef2f6; display: flex; color: #1f2937; }
.sidebar { width: 250px; background: #3989d4; color: #fff; padding: 22px 16px; display: flex; flex-direction: column; }
.brand { display: flex; align-items: center; gap: 10px; padding: 10px; border-bottom: 1px solid rgba(255,255,255,.24); }
.logo { width: 48px; height: 48px; border-radius: 8px; background: #fff; color: #3989d4; display: grid; place-items: center; font-size: 24px; }
.brand h2 { margin: 0; font-size: 28px; }
.brand p { margin: 3px 0 0; font-size: 16px; opacity: .9; }
.menu { margin-top: 18px; display: flex; flex-direction: column; gap: 10px; }
.menu a { color: #eaf4ff; text-decoration: none; padding: 10px 12px; border-radius: 8px; display: flex; align-items: center; gap: 8px; }
.menu-icon { width: 20px; height: 20px; display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; font-size: 16px; line-height: 1; }
.menu-label { line-height: 20px; }
.menu a.active, .menu a:hover { background: rgba(255,255,255,.2); }
.content { flex: 1; padding: 18px 20px; }
.topbar { background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; align-items: center; }
.topbar h1 { margin: 0; font-size: 32px; }
.user-box { display: flex; align-items: center; gap: 10px; }
.user-box button { border: 1px solid #d1d5db; background: #fff; border-radius: 6px; padding: 6px 10px; cursor: pointer; }
.banner { margin-top: 14px; background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; align-items: center; }
.banner p { margin: 0; }
.banner strong { color: #ef4444; }
.banner button { background: #3989d4; color: #fff; border: none; border-radius: 8px; padding: 10px 20px; cursor: pointer; }
.section-head { margin-top: 18px; display: flex; justify-content: space-between; align-items: center; }
.section-head h3 { margin: 0; font-size: 22px; }
.section-head a { color: #3b82f6; text-decoration: none; }
.task-grid { margin-top: 10px; display: grid; grid-template-columns: repeat(3, minmax(0,1fr)); gap: 12px; }
.task-card { background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; padding: 10px 12px; }
.row { display: flex; justify-content: space-between; align-items: center; }
.row.minor { font-size: 13px; margin-top: 8px; color: #4b5563; }
.task-card em { font-style: normal; border-radius: 999px; padding: 2px 8px; font-size: 12px; }
.s-done { background: #dcfce7; color: #15803d; }
.s-wait { background: #fef3c7; color: #b45309; }
.s-doing { background: #dbeafe; color: #1d4ed8; }
.s-urgent { background: #fee2e2; color: #b91c1c; }
.progress { margin-top: 8px; height: 8px; border-radius: 999px; background: #e5e7eb; overflow: hidden; }
.progress span { display: block; height: 100%; border-radius: inherit; }
.dist-panel { margin-top: 10px; background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; min-height: 180px; display: flex; align-items: center; justify-content: space-around; }
.donut { width: 160px; height: 160px; border-radius: 50%; background: conic-gradient(#16a34a 0 30%, #f59e0b 30% 42%, #0284c7 42% 58%, #7c3aed 58% 76%, #4f46e5 76% 90%, #ef4444 90% 100%); position: relative; }
.donut::after { content: ''; position: absolute; inset: 30px; border-radius: 50%; background: #fff; }
.dist-panel ul { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
.dot { display: inline-block; width: 10px; height: 10px; border-radius: 50%; margin-right: 8px; }
.table-panel { margin-top: 10px; background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; overflow: hidden; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px; border-bottom: 1px solid #edf0f5; text-align: center; }
.profile-panel { margin-top: 14px; background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; padding: 14px; }
.profile-panel h3 { margin: 0 0 10px; }
.profile-form { display: grid; grid-template-columns: repeat(4, minmax(0,1fr)); gap: 10px; align-items: end; }
.profile-form label { display: flex; flex-direction: column; gap: 4px; font-size: 14px; }
.profile-form input, .profile-form select { border: 1px solid #d1d5db; border-radius: 6px; padding: 8px; }
.profile-form button { background: #2563eb; color: #fff; border: none; border-radius: 6px; padding: 9px 12px; cursor: pointer; }
.message { margin: 10px 0 0; color: #2563eb; }
.page-footer { margin-top: 14px; background: #fff; border: 1px solid #d5dbe5; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; }


.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.35);
  display: grid;
  place-items: center;
  z-index: 30;
}
.modal-card {
  width: min(520px, calc(100vw - 32px));
  background: #fff;
  border-radius: 12px;
  border: 1px solid #dbe3ef;
  padding: 16px;
}
.modal-card h3 { margin: 0 0 12px; }
.task-form { display: grid; gap: 10px; }
.task-form label { display: grid; gap: 6px; font-size: 14px; }
.task-form input, .task-form select { border: 1px solid #cfd8e3; border-radius: 8px; padding: 8px 10px; }
.task-actions { display: flex; justify-content: flex-end; gap: 10px; }
.task-actions button { border: 1px solid #cad5e4; border-radius: 8px; padding: 8px 14px; cursor: pointer; background: #fff; }
.task-actions .primary { background: #3f8fdb; border-color: #3f8fdb; color: #fff; }

@media (max-width: 1100px) {
  .dashboard-page { flex-direction: column; }
  .sidebar { width: auto; }
  .task-grid { grid-template-columns: repeat(2, minmax(0,1fr)); }
  .profile-form { grid-template-columns: 1fr; }
}
@media (max-width: 760px) {
  .task-grid { grid-template-columns: 1fr; }
  .dist-panel { flex-direction: column; padding: 16px 0; }
  .topbar, .banner, .page-footer { flex-direction: column; gap: 10px; align-items: flex-start; }
}
.actions button.secondary { background: #909399; }
.actions button.danger { background: #f56c6c; }
.actions button:disabled { opacity: 0.7; cursor: not-allowed; }
</style>
