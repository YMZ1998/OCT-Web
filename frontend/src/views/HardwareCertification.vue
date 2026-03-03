<template>
  <div class="cert-page">
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
        <RouterLink class="active" :to="`/projects/${route.params.id}`">📂 项目管理</RouterLink>
        <a href="#" @click.prevent>📊 查询统计</a>
        <a href="#" @click.prevent>✅ 质控管理</a>
        <a href="#" @click.prevent>⚙️ 系统设置</a>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h1>项目管理-硬件认证</h1>
        <div class="user-box">
          <span>👤 {{ user?.username}}</span>
        </div>
      </header>

      <section class="tabs">
        <button :class="['tab', stage === 'hardware' ? 'active' : '']">
          <span class="tab-icon">🧩</span>
          <span>硬件认证</span>
        </button>
        <button :class="['tab', stage === 'technician' ? 'active' : '']">
          <span class="tab-icon">👨‍🔧</span>
          <span>技师认证</span>
        </button>
        <button :class="['tab', stage === 'certificate' ? 'active' : '']">
          <span class="tab-icon">🏅</span>
          <span>证书颁发</span>
        </button>
      </section>

      <section class="main-grid">
        <div class="panel detail-panel">
          <h3>硬件认证详情</h3>
          <p class="sub">设备信息</p>
          <div class="detail-box">
            <div><small>设备名称</small><strong>眼科光学相干断层扫描仪（OCT）</strong></div>
            <div><small>设备型号</small><strong>OCT-SJTU-3000</strong></div>
            <div><small>序列号</small><strong>SJTU-OCT-2026-0102</strong></div>
            <div><small>制造商</small><strong>上海视研医疗科技</strong></div>
          </div>

          <p class="sub">附件文档</p>
          <div class="file-item" v-for="doc in docs" :key="doc.name">
            <div>
              <strong>{{ doc.name }}</strong>
              <small>{{ doc.size }} · {{ doc.date }}</small>
            </div>
            <button>下载</button>
          </div>
        </div>

        <div class="panel opinion-panel">
          <h3>审核意见</h3>
          <textarea v-model.trim="opinion" placeholder="医生在查看硬件详情后请输入审核意见"></textarea>

          <div class="action-row">
            <button class="danger" @click="review(false)">认证不通过</button>
            <button class="success" @click="review(true)">认证通过</button>
          </div>
          <button class="notify" @click="sendNotification">发送通知</button>
          <p v-if="formMessage" class="form-message">{{ formMessage }}</p>

          <div class="status-box">
            <p>当前流程：<strong>{{ stageLabel }}</strong></p>
            <p v-if="lastDecision">最近结论：{{ lastDecision }}</p>
          </div>

          <div class="msg-list">
            <h4>已发送至试验中心技师账户</h4>
            <ul>
              <li v-for="(msg, idx) in messages" :key="`${idx}-${msg.time}`">
                <strong>{{ msg.result }}</strong>
                <p>{{ msg.content }}</p>
                <small>{{ msg.time }}</small>
              </li>
            </ul>
            <p v-if="!messages.length" class="empty">暂无通知记录</p>
          </div>
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

type ReviewMessage = {
  result: '通过' | '不通过' | '通知';
  content: string;
  time: string;
};

type FlowState = {
  stage: 'hardware' | 'technician' | 'certificate';
  messages: ReviewMessage[];
  lastDecision: string;
};

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);

const opinion = ref('');
const formMessage = ref('');
const stage = ref<FlowState['stage']>('hardware');
const messages = ref<ReviewMessage[]>([]);
const lastDecision = ref('');

const docs = [
  { name: '硬件设备清单.pdf', size: '2.5MB', date: '2026-02-28' },
  { name: '校准证书合集.pdf', size: '4.8MB', date: '2026-02-28' },
];

const projectId = computed(() => String(route.query.projectId || 'default'));
const flowKey = computed(() => `oct-hardware-flow-${projectId.value}`);
const stageLabel = computed(() => {
  if (stage.value === 'technician') return '技师认证';
  if (stage.value === 'certificate') return '证书颁发';
  return '硬件认证';
});

function nowText() {
  return new Date().toLocaleString('zh-CN');
}

function persistState() {
  const payload: FlowState = {
    stage: stage.value,
    messages: messages.value,
    lastDecision: lastDecision.value,
  };
  localStorage.setItem(flowKey.value, JSON.stringify(payload));
}

function sendToTechnicianAccount(result: ReviewMessage['result'], content: string) {
  messages.value.unshift({
    result,
    content,
    time: nowText(),
  });
  persistState();
}

function review(pass: boolean) {
  if (!opinion.value) {
    formMessage.value = '请先填写审核意见，再提交结论。';
    return;
  }

  formMessage.value = '';

  if (pass) {
    lastDecision.value = '认证通过';
    stage.value = 'technician';
    sendToTechnicianAccount('通过', `【认证通过】${opinion.value}。流程已流转到“技师认证”。`);
  } else {
    lastDecision.value = '认证不通过';
    stage.value = 'hardware';
    sendToTechnicianAccount('不通过', `【认证不通过】${opinion.value}。请试验中心技师整改后重新提交。`);
  }
}

function sendNotification() {
  formMessage.value = '';
  const content = opinion.value || '请关注当前硬件认证任务并及时处理。';
  sendToTechnicianAccount('通知', content);
}

onMounted(async () => {
  if (!userStore.token) {
    router.push('/login');
    return;
  }

  const cached = localStorage.getItem(flowKey.value);
  if (cached) {
    try {
      const parsed: FlowState = JSON.parse(cached);
      stage.value = parsed.stage;
      messages.value = parsed.messages;
      lastDecision.value = parsed.lastDecision;
    } catch {
      // ignore invalid cache
    }
  }

  try {
    const res = await getUser(route.params.id as string, userStore.token);
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
  } catch {
    user.value = null;
  }
});
</script>

<style scoped>
.cert-page { min-height: 100vh; display: flex; background: #eef2f6; color: #1f2937; }
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
.tabs { margin: 10px 0 14px; display: flex; gap: 10px; border-bottom: 1px solid #d5dbe5; padding-bottom: 10px; }
.tab {
  border: none;
  background: transparent;
  color: #64748b;
  padding: 8px 12px;
  cursor: default;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border-bottom: 2px solid transparent;
  line-height: 1.2;
}
.tab-icon {
  width: 20px;
  height: 20px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}
.tab.active { color: #2563eb; border-bottom-color: #2563eb; }
.main-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 14px; }
.panel { background: #fff; border: 1px solid #d2dae6; border-radius: 8px; padding: 14px; }
.sub { margin: 12px 0 8px; color: #4b5563; }
.detail-box { border: 1px solid #e2e8f0; border-radius: 8px; padding: 10px; display: grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap: 14px; }
.detail-box div { display: grid; gap: 4px; }
.detail-box small { color: #64748b; }
.file-item { border: 1px solid #e2e8f0; border-radius: 8px; padding: 10px 12px; display: flex; justify-content: space-between; align-items: center; margin-top: 10px; }
.file-item div { display: grid; gap: 4px; }
.file-item small { color: #64748b; }
.file-item button { border: 1px solid #cbd5e1; background: #fff; border-radius: 6px; padding: 6px 12px; }
.opinion-panel textarea { width: 100%; min-height: 260px; border: 1px solid #cbd5e1; border-radius: 8px; padding: 10px; resize: vertical; }
.action-row { margin-top: 10px; display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.action-row button, .notify { border: none; border-radius: 8px; padding: 10px 12px; cursor: pointer; }
.danger { background: #fee2e2; color: #b91c1c; }
.success { background: #dcfce7; color: #15803d; }
.notify { margin-top: 10px; width: 100%; background: #3f8fdb; color: #fff; }
.form-message { margin: 8px 0 0; color: #b91c1c; }
.status-box { margin-top: 10px; padding: 10px; border: 1px solid #e2e8f0; border-radius: 8px; color: #334155; }
.status-box p { margin: 0 0 4px; }
.msg-list { margin-top: 10px; border-top: 1px dashed #dbe3ef; padding-top: 10px; }
.msg-list h4 { margin: 0 0 8px; }
.msg-list ul { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
.msg-list li { border: 1px solid #e2e8f0; border-radius: 8px; padding: 8px; }
.msg-list li p { margin: 6px 0; }
.msg-list li small { color: #64748b; }
.empty { color: #64748b; }

@media (max-width: 1100px) {
  .cert-page { flex-direction: column; }
  .sidebar { width: auto; }
  .main-grid { grid-template-columns: 1fr; }
  .detail-box { grid-template-columns: 1fr; }
}
</style>
