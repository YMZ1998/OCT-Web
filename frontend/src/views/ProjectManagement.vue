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
        <RouterLink :to="`/user/${user?.id || route.params.id}`">
          <span class="menu-icon">🏠</span>
          <span class="menu-label">首页</span>
        </RouterLink>
        <a class="active" href="#" @click.prevent>
          <span class="menu-icon">📂</span>
          <span class="menu-label">项目管理</span>
        </a>
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
        <h1>项目管理</h1>
        <div class="user-box">
          <span>👤 {{ user?.username }}</span>
          <button @click="onLogout">退出登录</button>
        </div>
      </header>

      <section class="notice">
        <p>
          今天是 {{ today }}，您有 <strong>{{ pendingCount }}</strong> 项待处理任务，
          共管理 <strong>{{ recentProjects.length }}</strong> 个最近项目
        </p>
        <button @click="openCreateModal">+ 新建项目</button>
      </section>

      <section class="main-grid">
        <div class="panel project-panel">
          <header>
            <h3>最近项目（{{ recentProjects.length }}）</h3>
          </header>
          <article v-for="project in recentProjects" :key="project.id" class="project-card">
            <div class="project-head">
              <h4>{{ project.name }}</h4>
              <span :class="['state', project.stateClass]">{{ project.state }}</span>
            </div>
            <p>{{ project.desc }}</p>
            <footer>
              <span>📅 {{ project.date }}</span>
              <span>👥 {{ project.members }}人参与</span>
              <button @click="openDetailModal(project)">查看详情</button>
              <button @click="openEditModal(project)">编辑项目</button>
              <button class="primary" @click="handleTask(project)">处理任务</button>
            </footer>
          </article>
        </div>

        <div class="panel todo-panel">
          <header>
            <h3>待处理（{{ pendingCount }}）</h3>
          </header>
          <div class="todo-groups" v-if="groupedTodoItems.length">
            <section class="todo-group" v-for="group in groupedTodoItems" :key="group.projectId">
              <h4>{{ group.projectName }}</h4>
              <ul>
                <li v-for="item in group.items" :key="item.key">
                  <div class="todo-title">{{ item.taskName }}</div>
                  <div class="todo-actions">
                    <button class="todo-action" @click="openCertification(item)">进入处理</button>
                    <button class="todo-done" @click="completeTask(item)">完成任务</button>
                  </div>
                </li>
              </ul>
            </section>
          </div>
          <p v-if="!todoItems.length" class="empty">暂无待处理任务</p>
        </div>
      </section>
    </main>

    <div v-if="showCreateModal" class="modal-mask" @click.self="closeCreateModal">
      <section class="modal-card" role="dialog" aria-modal="true" aria-labelledby="create-title">
        <h3 id="create-title">{{ editingProjectId ? '编辑项目' : '新建项目' }}</h3>
        <form class="form-grid" @submit.prevent="submitCreateProject">
          <label>
            项目名称
            <input v-model.trim="newProjectForm.name" required placeholder="请输入项目名称" />
          </label>
          <label>
            项目负责人
            <input v-model.trim="newProjectForm.owner" required placeholder="请输入负责人" />
          </label>
          <label>
            试验中心
            <input v-model.trim="newProjectForm.center" required placeholder="请输入试验中心" />
          </label>
          <label>
            开始日期
            <input v-model="newProjectForm.date" type="date" required />
          </label>
          <label>
            参与人数
            <input v-model.number="newProjectForm.members" type="number" min="1" required />
          </label>
          <label class="full">
            项目描述
            <textarea v-model.trim="newProjectForm.desc" rows="3" required placeholder="请输入项目描述"></textarea>
          </label>
          <div class="modal-actions full">
            <button type="button" class="ghost" @click="closeCreateModal">取消</button>
            <button type="submit" class="primary">{{ editingProjectId ? '保存修改' : '创建项目' }}</button>
          </div>
        </form>
      </section>
    </div>

    <div v-if="detailProject" class="modal-mask" @click.self="detailProject = null">
      <section class="modal-card" role="dialog" aria-modal="true" aria-labelledby="detail-title">
        <h3 id="detail-title">项目详情</h3>
        <dl class="detail-grid">
          <dt>项目名称</dt><dd>{{ detailProject.name }}</dd>
          <dt>项目负责人</dt><dd>{{ detailProject.owner }}</dd>
          <dt>试验中心</dt><dd>{{ detailProject.center }}</dd>
          <dt>开始日期</dt><dd>{{ detailProject.date }}</dd>
          <dt>参与人数</dt><dd>{{ detailProject.members }} 人</dd>
          <dt>项目状态</dt><dd>{{ detailProject.state }}</dd>
          <dt>项目描述</dt><dd>{{ detailProject.desc }}</dd>
        </dl>
        <div class="modal-actions">
          <button class="primary" @click="detailProject = null">关闭</button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { getUser, updateUserProjects } from '../api/user';
import { useProjectStore, type ProjectItem, type ProjectStateData, type TodoItem } from '../store/project';
import { useUserStore } from '../store/user';
import type { User } from '../types/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const projectStore = useProjectStore();
const user = ref<User | null>(null);

const showCreateModal = ref(false);
const detailProject = ref<ProjectItem | null>(null);
const editingProjectId = ref<number | null>(null);

const recentProjects = computed(() => projectStore.recentProjects);
const todoItems = computed(() => projectStore.todoItems);
const pendingCount = computed(() => projectStore.pendingCount);
const today = computed(() => new Date().toLocaleDateString('zh-CN'));
const groupedTodoItems = computed(() => {
  const map = new Map<number, { projectId: number; projectName: string; items: TodoItem[] }>();
  todoItems.value.forEach((item) => {
    if (!map.has(item.projectId)) {
      map.set(item.projectId, {
        projectId: item.projectId,
        projectName: item.projectName,
        items: [],
      });
    }
    map.get(item.projectId)?.items.push(item);
  });
  return Array.from(map.values());
});

const newProjectForm = ref({
  name: '',
  owner: '',
  center: '',
  date: '',
  members: 1,
  desc: '',
});


async function syncProjectsToServer() {
  if (!userStore.token) return;
  const userId = String(user.value?.id || route.params.id || userStore.userInfo?.id || '');
  if (!userId) return;

  const payload: ProjectStateData = {
    recentProjects: projectStore.recentProjects,
    todoItems: projectStore.todoItems,
    nextProjectId: projectStore.nextProjectId,
  };

  try {
    await updateUserProjects(userId, payload, userStore.token);
  } catch {
    // 后端同步失败时保留前端状态，避免阻塞用户操作
  }
}

function resetNewProjectForm() {
  editingProjectId.value = null;
  newProjectForm.value = {
    name: '',
    owner: '',
    center: '',
    date: '',
    members: 1,
    desc: '',
  };
}

function openCreateModal() {
  showCreateModal.value = true;
}

function closeCreateModal() {
  showCreateModal.value = false;
  resetNewProjectForm();
}

async function submitCreateProject() {
  const payload = {
    name: newProjectForm.value.name,
    owner: newProjectForm.value.owner,
    center: newProjectForm.value.center,
    desc: newProjectForm.value.desc,
    date: newProjectForm.value.date,
    members: newProjectForm.value.members,
  };

  if (editingProjectId.value) {
    projectStore.updateProject(editingProjectId.value, payload);
  } else {
    projectStore.addProject(payload);
  }

  await syncProjectsToServer();
  closeCreateModal();
  void syncProjectsToServer();
}

async function handleTask(project: ProjectItem) {
  projectStore.startProjectTasks(project);
  void syncProjectsToServer();
}

async function completeTask(item: TodoItem) {
  projectStore.completeTask(item.key);
  void syncProjectsToServer();
}


function openEditModal(project: ProjectItem) {
  editingProjectId.value = project.id;
  newProjectForm.value = {
    name: project.name,
    owner: project.owner,
    center: project.center,
    date: project.date,
    members: project.members,
    desc: project.desc,
  };
  showCreateModal.value = true;
}

function openDetailModal(project: ProjectItem) {
  detailProject.value = project;
}

function openCertification(item: TodoItem) {
  router.push(`/projects/${route.params.id}/hardware?projectId=${item.projectId}&task=${encodeURIComponent(item.taskName)}`);
}

onMounted(async () => {
  if (!userStore.token) {
    router.push('/login');
    return;
  }

  projectStore.initForUser(String(route.params.id || userStore.userInfo?.id || 'guest'));

  try {
    const res = await getUser(route.params.id as string, userStore.token);
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
    const userId = String(res.data.data?.id || route.params.id || 'guest');
    projectStore.initForUser(userId);

    if (res.data.data?.project_state) {
      projectStore.replaceState(res.data.data.project_state);
    } else {
      await syncProjectsToServer();
    }
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
.brand h2 { margin: 0; font-size: 28px; line-height: 1.15; font-weight: 700; }
.brand p { margin: 3px 0 0; font-size: 16px; line-height: 1.2; opacity: .9; }
.menu { margin-top: 18px; display: grid; gap: 12px; }
.menu a { color: #ecf5ff; text-decoration: none; padding: 10px 12px; border-radius: 8px; display: flex; align-items: center; gap: 8px; line-height: 20px; }
.menu-icon { width: 20px; height: 20px; display: inline-flex; align-items: center; justify-content: flex-start; flex-shrink: 0; font-size: 16px; line-height: 1; text-align: left; }
.menu-label { line-height: 20px; }
.menu a.active, .menu a:hover { background: rgba(255,255,255,.2); }
.content { flex: 1; padding: 18px 20px; }
.topbar { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; align-items: center; }
.topbar h1 { margin: 0; }
.user-box { display: flex; align-items: center; gap: 10px; }
.user-box button { border: 1px solid #d1d5db; background: #fff; border-radius: 6px; padding: 6px 10px; cursor: pointer; }
.notice { margin-top: 14px; background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px 18px; display: flex; justify-content: space-between; align-items: center; }
.notice strong { color: #ef4444; }
.notice button { background: #3f8fdb; color: #fff; border: none; border-radius: 8px; padding: 10px 20px; cursor: pointer; }
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
.state.completed { background: #10b981; }
.project-card footer { display: grid; grid-template-columns: auto auto 1fr auto auto auto; gap: 10px; align-items: center; color: #6b7280; }
.project-card footer button { justify-self: end; border: 1px solid #d1d5db; background: #fff; border-radius: 8px; padding: 8px 12px; cursor: pointer; }
.project-card footer button.primary { background: #dbe7ff; color: #1f3b8f; border-color: #c8d7ff; }
.todo-groups { padding: 12px; display: grid; gap: 12px; }
.todo-group { border: 1px dashed #d5deea; border-radius: 8px; padding: 10px; }
.todo-group h4 { margin: 0 0 8px; color: #334155; font-size: 16px; }
.todo-panel ul { list-style: none; margin: 0; padding: 0; display: grid; gap: 10px; }
.todo-panel li { border: 1px solid #dfe5f1; border-radius: 8px; padding: 14px; display: grid; gap: 8px; }
.todo-title { font-size: 20px; color: #1f2937; }
.todo-panel li small { color: #64748b; }
.todo-actions { display: flex; gap: 8px; }
.todo-action { justify-self: start; border: 1px solid #c8d7ff; background: #edf3ff; color: #1f3b8f; border-radius: 6px; padding: 6px 10px; cursor: pointer; }
.todo-done { border: 1px solid #86efac; background: #dcfce7; color: #166534; border-radius: 6px; padding: 6px 10px; cursor: pointer; }
.empty { margin: 10px 12px 14px; color: #64748b; }

.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.35);
  display: grid;
  place-items: center;
  z-index: 30;
}
.modal-card {
  width: min(760px, calc(100vw - 32px));
  background: #fff;
  border-radius: 12px;
  border: 1px solid #dbe3ef;
  padding: 16px;
}
.modal-card h3 { margin: 0 0 14px; }
.form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.form-grid label { display: flex; flex-direction: column; gap: 6px; font-size: 14px; }
.form-grid input, .form-grid textarea {
  border: 1px solid #cfd8e3;
  border-radius: 8px;
  padding: 8px 10px;
  font: inherit;
}
.full { grid-column: 1 / -1; }
.modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 8px; }
.modal-actions button { border: 1px solid #cad5e4; border-radius: 8px; padding: 8px 14px; cursor: pointer; }
.modal-actions .primary { background: #3f8fdb; border-color: #3f8fdb; color: #fff; }
.modal-actions .ghost { background: #fff; }
.detail-grid {
  display: grid;
  grid-template-columns: 110px 1fr;
  gap: 8px 10px;
  margin: 0;
}
.detail-grid dt { color: #64748b; }
.detail-grid dd { margin: 0; color: #0f172a; }

@media (max-width: 1100px) {
  .project-page { flex-direction: column; }
  .sidebar { width: auto; }
  .main-grid { grid-template-columns: 1fr; }
  .project-card footer { grid-template-columns: 1fr; }
  .project-card footer button { justify-self: start; }
  .form-grid { grid-template-columns: 1fr; }
}
</style>
