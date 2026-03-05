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
        <h1>{{ pageTitle }}</h1>
        <div class="user-box">
          <span>👤 {{ user?.username }}</span>
        </div>
      </header>

      <section class="tabs">
        <button :class="tabClass('hardware')" @click="switchStage('hardware')">
          <span class="tab-icon">🧩</span>
          <span>硬件认证</span>
        </button>
        <button :class="tabClass('technician')" @click="switchStage('technician')">
          <span class="tab-icon">👨‍🔧</span>
          <span>技师认证</span>
        </button>
        <button :class="tabClass('certificate')" @click="switchStage('certificate')">
          <span class="tab-icon">🏅</span>
          <span>证书颁发</span>
        </button>
      </section>



      <section class="main-grid">
        <div class="panel detail-panel">
          <HardwareStagePanel v-if="stage === 'hardware'" :docs="docs" />

          <TechnicianStagePanel
            v-else-if="stage === 'technician'"
            :distribution-step="distributionStep"
            :selected-task-ids="selectedTaskIds"
            :paged-images="pagedImages"
            :distribution-message="distributionMessage"
            :active-task-id="activeTaskId"
            :current-page="currentPage"
            :image-total-pages="imageTotalPages"
            @distribute="distributeSelected"
            @toggle-selection="toggleTaskSelection"
            @show-task-detail="showTaskDetail"
            @go-page="goToPage"
          />

          <CertificateStagePanel
            v-else
            :project-id="projectId"
            :issue-date="issueDate"
            :expiry-date="expiryDate"
            :certificate-no="certificateNo"
            :qr-cells="qrCells"
            :use-seal="useSeal"
            :use-sign="useSign"
            :certificate-issued="certificateIssued"
            @update:seal="updateUseSeal"
            @update:sign="updateUseSign"
            @issue="issueCertificate"
            @print="printCertificate"
            @download="downloadCertificate"
          />
        </div>

        <div class="panel opinion-panel">
          <h3>审核意见</h3>
          <textarea v-model.trim="opinion" :placeholder="opinionPlaceholder"></textarea>

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

    <div v-if="showImagePreview && activeTaskDetail" class="modal-mask" @click.self="showImagePreview = false">
      <section class="modal-card" role="dialog" aria-modal="true" aria-labelledby="image-preview-title">
        <h3 id="image-preview-title">完整影像预览</h3>
        <p class="sub">{{ activeTaskDetail.sample }} · {{ activeTaskDetail.type }}</p>
        <div class="image-preview">{{ activeTaskDetail.sample }} 影像图</div>
        <div class="modal-actions">
          <button class="primary" @click="showImagePreview = false">关闭</button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { getUser } from '../api/user';
import CertificateStagePanel from '../components/hardware-certification/CertificateStagePanel.vue';
import HardwareStagePanel from '../components/hardware-certification/HardwareStagePanel.vue';
import TechnicianStagePanel from '../components/hardware-certification/TechnicianStagePanel.vue';
import { useUserStore } from '../store/user';
import type { User } from '../types/user';

type ReviewMessage = {
  result: '通过' | '不通过' | '通知';
  content: string;
  time: string;
};

type FlowState = {
  stage: 'hardware' | 'technician' | 'certificate';
  distributionStep: 'screening' | 'inspection' | 'reading' | 'quality';
  messages: ReviewMessage[];
  lastDecision: string;
  currentPage: number;
  certificateNo: string;
  certificateIssued: boolean;
  useSeal: boolean;
  useSign: boolean;
  passedStages: Record<'hardware' | 'technician' | 'certificate', boolean>;
};

type ImageItem = {
  id: number;
  sample: string;
  patient: string;
  age: number;
  date: string;
  imageCount: number;
  type: string;
};

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);

const opinion = ref('');
const formMessage = ref('');
const managerOpinion = ref('');
const showFullReport = ref(false);
const stage = ref<FlowState['stage']>('technician');
const messages = ref<ReviewMessage[]>([]);
const lastDecision = ref('');
const currentPage = ref(1);
const distributionStep = ref<FlowState['distributionStep']>('screening');
const selectedTaskIds = ref<number[]>([]);
const activeTaskId = ref<number | null>(null);
const distributionMessage = ref('');
const showImagePreview = ref(false);
const certificateNo = ref('');
const certificateIssued = ref(false);
const useSeal = ref(false);
const useSign = ref(false);
const passedStages = ref<Record<FlowState['stage'], boolean>>({
  hardware: false,
  technician: false,
  certificate: false,
});
const pageSize = 12;


const docs = [
  { name: '设备配置说明.pdf', size: '2.1MB', date: '2026-01-09' },
  { name: '计量校准报告.pdf', size: '1.6MB', date: '2026-01-08' },
  { name: '设备维护记录.pdf', size: '0.9MB', date: '2026-01-07' },
];

const imageData: ImageItem[] = Array.from({ length: 36 }, (_, idx) => ({
  id: idx + 1,
  sample: `病例#${String(idx + 1).padStart(6, '0')}`,
  patient: `张${String.fromCharCode(65 + (idx % 26))}`,
  age: 45 + (idx % 25),
  date: `2026-02-${String((idx % 28) + 1).padStart(2, '0')}`,
  imageCount: 12 + (idx % 8),
  type: 'OCT',
}));


const imageTotalPages = computed(() => Math.ceil(imageData.length / pageSize));
const pagedImages = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return imageData.slice(start, start + pageSize);
});
const activeTaskDetail = computed(() => imageData.find((item) => item.id === activeTaskId.value) || null);
const activeTaskProgress = computed(() => {
  const id = activeTaskDetail.value?.id || 0;
  const percent = 65 + (id % 30);
  return {
    percent,
    flow: [
      { label: '任务创建', time: '2026-02-10 09:15' },
      { label: '初级读片完成', time: '2026-02-10 11:20' },
      { label: '主管读片完成', time: '2026-02-10 14:35' },
      { label: '项目经理审核中', time: '待处理' },
    ],
  };
});
const activeTaskReport = computed(() => {
  const sample = activeTaskDetail.value?.sample || '当前任务';
  return {
    junior: `${sample} 黄斑区可见轻度水肿，建议结合随访观察。`,
    seniorReport: `${sample} OCT影像层次清晰，黄斑中心凹结构基本完整，未见明显出血征象。`,
    seniorOpinion: '建议纳入下一阶段随访，维持当前治疗方案并加强复查频率。',
  };
});

const projectId = computed(() => String(route.query.projectId || 'XXXXXXXXXX'));
const flowKey = computed(() => `oct-hardware-flow-${projectId.value}`);
const stageLabel = computed(() => {
  if (stage.value === 'technician') return '技师认证';
  if (stage.value === 'certificate') return '证书颁发';
  return '硬件认证';
});
const pageTitle = '项目管理-认证';
const opinionPlaceholder = computed(() => {
  if (stage.value === 'technician') return '医生在查看影像数据后请输入审核意见';
  if (stage.value === 'hardware') return '医生在查看硬件详情后请输入审核意见';
  return '请补充证书颁发阶段相关意见';
});
const issueDate = computed(() => new Date().toISOString().slice(0, 10));
const expiryDate = computed(() => {
  const d = new Date();
  d.setFullYear(d.getFullYear() + 1);
  return d.toISOString().slice(0, 10);
});
const tabClass = (target: FlowState['stage']) => ({
  tab: true,
  active: stage.value === target,
  passed: passedStages.value[target],
});

const qrCells = computed(() => {
  const source = certificateNo.value || `${projectId.value}-PENDING`;
  const chars = source.replace(/[^A-Za-z0-9]/g, '');
  return Array.from({ length: 64 }, (_, idx) => {
    const ch = chars[idx % Math.max(chars.length, 1)] || '0';
    return (ch.charCodeAt(0) + idx) % 3 !== 0;
  });
});

function nowText() {
  return new Date().toLocaleString('zh-CN');
}

function ensureCertificateMeta() {
  if (!certificateNo.value) {
    const rand = Math.random().toString(36).slice(2, 8).toUpperCase();
    certificateNo.value = `OCT-${new Date().getFullYear()}-${rand}`;
  }
}

function persistState() {
  const payload: FlowState = {
    stage: stage.value,
    distributionStep: distributionStep.value,
    messages: messages.value,
    lastDecision: lastDecision.value,
    currentPage: currentPage.value,
    certificateNo: certificateNo.value,
    certificateIssued: certificateIssued.value,
    useSeal: useSeal.value,
    useSign: useSign.value,
    passedStages: passedStages.value,
  };
  localStorage.setItem(flowKey.value, JSON.stringify(payload));
}

function switchDistributionStep(nextStep: FlowState['distributionStep']) {
  distributionStep.value = nextStep;
  persistState();
}

function updateUseSeal(next: boolean) {
  useSeal.value = next;
  persistState();
}

function updateUseSign(next: boolean) {
  useSign.value = next;
  persistState();
}

function goToPage(page: number) {
  if (page < 1 || page > imageTotalPages.value) return;
  currentPage.value = page;
  persistState();
}

function toggleTaskSelection(id: number) {
  if (selectedTaskIds.value.includes(id)) {
    selectedTaskIds.value = selectedTaskIds.value.filter((taskId) => taskId !== id);
    return;
  }
  selectedTaskIds.value.push(id);
}

function showTaskDetail(id: number) {
  activeTaskId.value = id;
  showFullReport.value = false;
  formMessage.value = '';
}

function viewFullReport(id?: number) {
  if (id) activeTaskId.value = id;
  if (!activeTaskDetail.value) {
    formMessage.value = '请先选择任务。';
    return;
  }
  showFullReport.value = true;
  formMessage.value = '';
}

function distributeSelected(mode: 'batch' | 'smart') {
  if (!selectedTaskIds.value.length) return;
  distributionMessage.value = mode === 'batch'
    ? `已将 ${selectedTaskIds.value.length} 份影像资料批量分发给指定医生。`
    : `已将 ${selectedTaskIds.value.length} 份影像资料智能分发给适配医生。`;
  selectedTaskIds.value = [];
}

function viewFullImage() {
  if (!activeTaskDetail.value) {
    formMessage.value = '请先选择任务并查看详情。';
    return;
  }
  formMessage.value = `已打开 ${activeTaskDetail.value.sample} 的完整影像。`;
  showImagePreview.value = true;
}

function submitReadingDecision(pass: boolean) {
  if (!activeTaskDetail.value) {
    formMessage.value = '请先选择任务后再提交审核结论。';
    return;
  }
  if (!managerOpinion.value) {
    formMessage.value = '请先填写报告意见。';
    return;
  }
  const prefix = pass ? '通过' : '不通过';
  const content = `【阅片审核${prefix}】${activeTaskDetail.value.sample}：${managerOpinion.value}`;
  sendToTechnicianAccount(pass ? '通过' : '不通过', content);
  formMessage.value = `已提交${prefix}结论并通知试验中心。`;
}

function switchStage(nextStage: FlowState['stage']) {
  stage.value = nextStage;
  if (nextStage === 'certificate') ensureCertificateMeta();
  if (nextStage === 'technician') {
    distributionStep.value = 'screening';
  }
  if (nextStage !== 'technician') {
    selectedTaskIds.value = [];
    distributionMessage.value = '';
    activeTaskId.value = null;
    showImagePreview.value = false;
  }
  persistState();
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
    passedStages.value[stage.value] = true;
    if (stage.value === 'hardware') {
      stage.value = 'technician';
      sendToTechnicianAccount('通过', `【硬件认证通过】${opinion.value}。流程已流转到“技师认证”。`);
      return;
    }
    stage.value = 'certificate';
    ensureCertificateMeta();
    sendToTechnicianAccount('通过', `【认证通过】${opinion.value}。流程已流转到“证书颁发”。`);
  } else {
    lastDecision.value = '认证不通过';
    passedStages.value[stage.value] = false;
    if (stage.value === 'hardware') {
      stage.value = 'hardware';
      sendToTechnicianAccount('不通过', `【硬件认证不通过】${opinion.value}。请先完成设备整改后再进入技师认证。`);
      return;
    }
    stage.value = 'technician';
    sendToTechnicianAccount('不通过', `【认证不通过】${opinion.value}。请试验中心技师整改后重新提交。`);
  }
}

function issueCertificate() {
  if (stage.value !== 'certificate') {
    formMessage.value = '请先进入证书颁发阶段。';
    return;
  }
  ensureCertificateMeta();
  certificateIssued.value = true;
  passedStages.value.certificate = true;
  sendToTechnicianAccount('通知', `【证书已颁发】证书编号：${certificateNo.value}，请试验中心技师账户查收。`);
}

function printCertificate() {
  window.print();
}

function downloadCertificate() {
  ensureCertificateMeta();
  const content = [
    '认证证书',
    `项目ID: ${projectId.value}`,
    `证书编号: ${certificateNo.value}`,
    `颁证日期: ${issueDate.value}`,
    `有效期至: ${expiryDate.value}`,
    `印章: ${useSeal.value ? '已添加' : '未添加'}`,
    `电子签名: ${useSign.value ? '已添加' : '未添加'}`,
  ].join('\n');
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${certificateNo.value}.txt`;
  a.click();
  URL.revokeObjectURL(url);
}

function sendNotification() {
  formMessage.value = '';
  const content = opinion.value || '请关注当前硬件认证任务并及时处理。';
  sendToTechnicianAccount('通知', content);
}

function syncStageByTaskQuery() {
  const task = String(route.query.task || '');
  if (!task) return;

  if (task.includes('认证')) {
    stage.value = 'technician';
    distributionStep.value = 'screening';
    return;
  }

  if (task.includes('分发') || task.includes('检查')) {
    stage.value = 'technician';
    distributionStep.value = 'inspection';
    return;
  }

  if (task.includes('阅片')) {
    stage.value = 'technician';
    distributionStep.value = 'reading';
    return;
  }

  stage.value = 'technician';
  distributionStep.value = 'screening';
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
      distributionStep.value = parsed.distributionStep || 'screening';
      messages.value = parsed.messages;
      lastDecision.value = parsed.lastDecision;
      currentPage.value = parsed.currentPage || 1;
      certificateNo.value = parsed.certificateNo || '';
      certificateIssued.value = parsed.certificateIssued || false;
      useSeal.value = parsed.useSeal || false;
      useSign.value = parsed.useSign || false;
      passedStages.value = {
        hardware: parsed.passedStages?.hardware || false,
        technician: parsed.passedStages?.technician || false,
        certificate: parsed.passedStages?.certificate || false,
      };
    } catch {
      // ignore invalid cache
    }
  }

  if (stage.value === 'certificate') ensureCertificateMeta();
  syncStageByTaskQuery();

  try {
    const res = await getUser(route.params.id as string, userStore.token);
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
  } catch {
    user.value = null;
  }
});

watch(
  () => route.query.task,
  () => {
    syncStageByTaskQuery();
  },
);
</script>

<style>
.cert-page { min-height: 100vh; display: flex; background: #eef2f6; color: #1f2937; }
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
.tabs { margin: 10px 0 14px; display: flex; gap: 10px; border-bottom: 1px solid #d5dbe5; padding-bottom: 10px; }
.tab {
  border: none;
  background: transparent;
  color: #64748b;
  padding: 8px 12px;
  cursor: pointer;
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
.tab.passed { color: #16a34a; border-bottom-color: #16a34a; }
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
.distribution-top-steps { margin: 10px 0 14px; display: flex; gap: 20px; border-bottom: 1px solid #d2dae6; padding: 4px 0 10px; color: #475569; }
.step-btn { border: none; background: transparent; color: #475569; font: inherit; cursor: pointer; padding: 0; }
.step-btn.active { color: #2563eb; font-weight: 600; }
.task-header { display: flex; justify-content: space-between; align-items: center; }
.tech-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.tech-card { border: 1px solid #d2dae6; border-radius: 8px; overflow: hidden; cursor: pointer; background: #fff; }
.tech-card.active { border-color: #3f8fdb; box-shadow: 0 0 0 2px rgba(63,143,219,.15) inset; }
.tech-thumb { height: 110px; background: #eef2f6; display: grid; place-items: center; font-size: 36px; color: #94a3b8; }
.tech-meta { padding: 8px 10px; display: grid; gap: 2px; }
.tech-meta small { color: #94a3b8; }
.task-actions { display: flex; gap: 10px; }
.task-actions button { border: 1px solid #c8d7ff; background: #e8f0ff; color: #1f3b8f; border-radius: 8px; padding: 7px 12px; cursor: pointer; }
.task-actions button:disabled { opacity: .5; cursor: not-allowed; }
.task-list { display: grid; gap: 10px; }
.task-card { border: 1px solid #d2dae6; border-radius: 8px; padding: 10px; display: flex; justify-content: space-between; align-items: center; }
.task-card label { display: flex; align-items: flex-start; gap: 10px; flex: 1; }
.task-main { display: grid; gap: 3px; }
.task-main small { color: #64748b; }
.detail-link { border: none; background: transparent; color: #2563eb; cursor: pointer; }
.task-card-actions { display: flex; flex-direction: column; gap: 6px; align-items: flex-end; }
.task-detail-panel { display: grid; gap: 12px; align-content: start; }
.reading-preview-card, .task-progress-card, .report-card, .manager-review-card { border: 1px solid #dbe3ef; border-radius: 8px; padding: 10px; background: #fff; }
.reading-preview-card header, .report-card header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.task-progress-card h4, .report-card h4, .manager-review-card h4 { margin: 0 0 8px; }
.progress-row { display: flex; justify-content: space-between; margin-bottom: 6px; color: #334155; }
.progress-track { height: 8px; background: #e2e8f0; border-radius: 999px; overflow: hidden; }
.progress-track span { display: block; height: 100%; background: #3f8fdb; }
.task-progress-card ul { list-style: none; padding: 0; margin: 10px 0 0; display: grid; gap: 6px; }
.task-progress-card li { display: flex; justify-content: space-between; gap: 12px; }
.task-progress-card li small { color: #64748b; }
.manager-review-card textarea { width: 100%; min-height: 100px; border: 1px solid #cbd5e1; border-radius: 8px; padding: 8px; resize: vertical; }
.distribution-message { margin: 12px 0 0; color: #166534; }
.pager { margin-top: 12px; display: flex; justify-content: center; gap: 8px; }
.pager button { border: 1px solid #cbd5e1; border-radius: 6px; min-width: 32px; height: 32px; background: #fff; color: #334155; cursor: pointer; }
.pager button:disabled { cursor: not-allowed; opacity: .5; }
.pager .page-btn.active { background: #3f8fdb; color: #fff; border-color: #3f8fdb; }
.certificate-card { border: 1px solid #d2dae6; border-radius: 8px; overflow: hidden; background: #fff; }
.certificate-card header { text-align: center; padding: 18px 12px; border-bottom: 1px solid #e2e8f0; }
.certificate-card header h4 { margin: 0; font-size: 30px; }
.certificate-card header p { margin: 6px 0 0; color: #64748b; }
.certificate-main { text-align: center; padding: 20px 16px; display: grid; gap: 6px; }
.certificate-main h2 { margin: 0; font-size: 36px; }
.certificate-main p { margin: 0; color: #334155; }
.certificate-main small { color: #64748b; }
.certificate-meta { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); padding: 14px 16px; border-top: 1px solid #e2e8f0; }
.certificate-meta div { display: grid; gap: 4px; }
.certificate-meta span { color: #64748b; }
.certificate-signature { display: flex; justify-content: space-between; align-items: flex-end; border-top: 1px solid #e2e8f0; padding: 14px 16px; }
.qr-box { width: 88px; height: 88px; display: grid; grid-template-columns: repeat(8, 1fr); gap: 2px; background: #e2e8f0; padding: 4px; }
.qr-cell { background: #fff; }
.qr-cell.on { background: #1f2937; }
.signature-box { text-align: right; display: grid; gap: 8px; }
.controls { display: flex; gap: 14px; color: #475569; }
.marks { min-height: 38px; display: flex; justify-content: flex-end; gap: 10px; }
.seal, .sign { border: 1px solid #cbd5e1; border-radius: 6px; padding: 6px 8px; }
.seal { color: #dc2626; border-color: #fecaca; background: #fff1f2; }
.sign { color: #2563eb; border-color: #bfdbfe; background: #eff6ff; }
.certificate-actions { margin-top: 12px; display: flex; justify-content: space-between; }
.certificate-actions .issue, .certificate-actions .secondary { border: none; border-radius: 8px; padding: 10px 16px; cursor: pointer; }
.certificate-actions .issue { background: #dbeafe; color: #1d4ed8; }
.certificate-actions .issue:disabled { opacity: .6; cursor: not-allowed; }
.certificate-actions .secondary { background: #e2e8f0; color: #334155; margin-left: 8px; }
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
.image-preview { height: 280px; border: 1px solid #d2dae6; border-radius: 8px; background: linear-gradient(135deg, #f8fbff, #e6eefb); display: grid; place-items: center; color: #1e3a8a; font-size: 22px; }
.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.35);
  display: grid;
  place-items: center;
  z-index: 40;
}
.modal-card {
  width: min(720px, calc(100vw - 32px));
  background: #fff;
  border-radius: 12px;
  border: 1px solid #dbe3ef;
  padding: 16px;
}
.modal-actions { margin-top: 12px; display: flex; justify-content: flex-end; }
.modal-actions .primary { border: none; border-radius: 8px; padding: 8px 16px; background: #3f8fdb; color: #fff; cursor: pointer; }

@media (max-width: 1100px) {
  .cert-page { flex-direction: column; }
  .sidebar { width: auto; }
  .main-grid { grid-template-columns: 1fr; }
  .detail-box { grid-template-columns: 1fr; }
  .distribution-top-steps { flex-wrap: wrap; gap: 10px; }
  .certificate-meta { grid-template-columns: 1fr; gap: 10px; }
  .certificate-signature { flex-direction: column; align-items: flex-start; gap: 10px; }
  .controls { flex-direction: column; gap: 6px; }
  .certificate-actions { flex-direction: column; gap: 10px; }
}
</style>
