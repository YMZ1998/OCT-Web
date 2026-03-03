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
        <a class="active" href="#">🏠 首页</a>
        <a href="#">📁 项目管理</a>
        <a href="#">📊 查询统计</a>
        <a href="#">✅ 质控管理</a>
        <a href="#">⚙️ 系统设置</a>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h1>首页</h1>
        <div class="user-box">
          <span>👤 {{ user?.username || '项目经理' }}</span>
          <button @click="onLogout">退出登录</button>
        </div>
      </header>

      <section class="banner">
        <p>今天是 {{ today }}，您有 <strong>{{ pendingCount }}</strong> 项待处理任务</p>
        <button>+ 新建任务</button>
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
            <tr>
              <td>儿童眼底病变星筛查项目</td>
              <td>{{ user?.email || 'xxxxxxx' }}</td>
              <td>{{ user?.created_at ? formatDate(user.created_at) : '2023-06-10' }}</td>
              <td><a href="#" @click.prevent>查看</a></td>
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

      <footer class="page-footer">
        <span>SJTURC眼科在吸纳工作管理系统</span>
        <span>版本1.0　<a href="#" @click.prevent>帮助</a></span>
      </footer>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getUser, updateUser } from '../api/user';
import { useUserStore } from '../store/user';
import type { User } from '../types/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);
const saving = ref(false);
const message = ref('');
const form = ref({ email: '', gender: '', age: undefined as number | undefined });

const taskStats = [
  { name: '硬件认证', done: 25, total: 25, progress: 100, color: '#16a34a', status: '已完成', statusClass: 's-done', icon: '✅' },
  { name: '技师认证', done: 0, total: 15, progress: 0, color: '#f59e0b', status: '待处理', statusClass: 's-wait', icon: '📌' },
  { name: '证书颁发', done: 12, total: 20, progress: 60, color: '#0284c7', status: '进行中', statusClass: 's-doing', icon: '🏅' },
  { name: '阅片分发', done: 18, total: 30, progress: 60, color: '#7c3aed', status: '进行中', statusClass: 's-doing', icon: '👁' },
  { name: '阅片审核', done: 8, total: 20, progress: 40, color: '#4f46e5', status: '进行中', statusClass: 's-doing', icon: '🔍' },
  { name: '质量抽查', done: 2, total: 10, progress: 20, color: '#ef4444', status: '紧急', statusClass: 's-urgent', icon: '⚠️' },
];

const pendingCount = computed(() => taskStats.filter((i) => i.status !== '已完成').length);
const today = computed(() => new Date().toLocaleDateString('zh-CN'));

function formatDate(ts: number) {
  return new Date(ts * 1000).toLocaleDateString();
}

function fillFormFromUser() {
  form.value.email = user.value?.email || '';
  form.value.gender = user.value?.gender || '';
  form.value.age = user.value?.age;
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
.menu a { color: #eaf4ff; text-decoration: none; padding: 10px 12px; border-radius: 8px; }
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
