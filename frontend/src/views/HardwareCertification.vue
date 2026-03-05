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

      <section class="tabs" v-if="stage !== 'technician'">
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

      <section v-else class="distribution-top-steps">
        <button :class="['step-btn', distributionStep === 'screening' ? 'active' : '']" @click="switchDistributionStep('screening')">受试者筛选阶段</button>
        <button :class="['step-btn', distributionStep === 'inspection' ? 'active' : '']" @click="switchDistributionStep('inspection')">影像数据检查阶段</button>
        <button :class="['step-btn', distributionStep === 'reading' ? 'active' : '']" @click="switchDistributionStep('reading')">阅片阶段</button>
        <button :class="['step-btn', distributionStep === 'quality' ? 'active' : '']" @click="switchDistributionStep('quality')">质量抽查</button>
      </section>

      <section class="main-grid">
        <div class="panel detail-panel">
          <template v-if="stage === 'hardware'">
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
          </template>

          <template v-else-if="stage === 'technician' && distributionStep === 'inspection'">
            <section class="task-header">
              <p class="sub">任务列表</p>
              <div class="task-actions">
                <button :disabled="!selectedTaskIds.length" @click="distributeSelected('batch')">批量分发</button>
                <button :disabled="!selectedTaskIds.length" @click="distributeSelected('smart')">智能分发</button>
              </div>
            </section>

            <div class="task-list">
              <article class="task-card" v-for="item in pagedImages" :key="item.id">
                <label>
                  <input type="checkbox" :checked="selectedTaskIds.includes(item.id)" @change="toggleTaskSelection(item.id)" />
                  <div class="task-main">
                    <strong>{{ item.sample }}</strong>
                    <small>患者：{{ item.patient }}｜年龄：{{ item.age }}岁｜检查类型：{{ item.type }}</small>
                    <small>{{ item.date }} · {{ item.imageCount }}张影像</small>
                  </div>
                </label>
                <button class="detail-link" @click="showTaskDetail(item.id)">查看详情</button>
              </article>
            </div>

            <p v-if="distributionMessage" class="distribution-message">{{ distributionMessage }}</p>

            <div class="pager" role="navigation" aria-label="影像分页">
              <button :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">‹</button>
              <button
                v-for="page in imageTotalPages"
                :key="page"
                :class="['page-btn', page === currentPage ? 'active' : '']"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
              <button :disabled="currentPage === imageTotalPages" @click="goToPage(currentPage + 1)">›</button>
            </div>
          </template>

          <template v-else-if="stage === 'technician'">
            <h3>影像分发</h3>
            <p class="sub">当前阶段：{{ distributionStepLabel }}</p>
            <div class="stage-placeholder">
              <p v-if="distributionStep === 'screening'">已进入受试者筛选阶段，可在顶部自由切换到“影像数据检查阶段”。</p>
              <p v-else-if="distributionStep === 'reading'">当前为阅片阶段，影像将进入医生阅片审核流程。</p>
              <p v-else>当前为质量抽查阶段，可对已分发影像进行质量追踪。</p>
            </div>
          </template>

          <template v-else>
            <h3>认证证书</h3>
            <p class="sub">进入证书颁发界面后自动生成证书编号与二维码，可添加印章和电子签名</p>

            <article class="certificate-card">
              <header>
                <h4>认证证书</h4>
                <p>SJTURC眼科在线工作管理系统</p>
              </header>
              <section class="certificate-main">
                <p>兹证明</p>
                <h2>眼科影像分析系统</h2>
                <small>项目ID: {{ projectId }}</small>
                <p>已通过硬件认证和技师认证，符合SJTURC眼科系统认证标准。</p>
              </section>
              <section class="certificate-meta">
                <div><span>颁证日期</span><strong>{{ issueDate }}</strong></div>
                <div><span>有效期至</span><strong>{{ expiryDate }}</strong></div>
                <div><span>证书编号</span><strong>{{ certificateNo }}</strong></div>
              </section>
              <section class="certificate-signature">
                <div>
                  <div class="qr-box">
                    <span v-for="(cell, index) in qrCells" :key="index" :class="['qr-cell', cell ? 'on' : '']"></span>
                  </div>
                  <small>扫码验真</small>
                </div>
                <div class="signature-box">
                  <div class="controls">
                    <label><input type="checkbox" v-model="useSeal" /> 添加印章</label>
                    <label><input type="checkbox" v-model="useSign" /> 添加电子签名</label>
                  </div>
                  <div class="marks">
                    <span v-if="useSeal" class="seal">SJTURC印章</span>
                    <span v-if="useSign" class="sign">授权签名</span>
                  </div>
                </div>
              </section>
            </article>

            <div class="certificate-actions">
              <button class="issue" :disabled="certificateIssued" @click="issueCertificate">{{ certificateIssued ? '已颁发' : '颁发证书' }}</button>
              <div>
                <button class="secondary" @click="printCertificate">打印证书</button>
                <button class="secondary" @click="downloadCertificate">下载证书</button>
              </div>
            </div>
          </template>
        </div>

        <div class="panel opinion-panel" v-if="stage !== 'technician'">
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

        <div class="panel task-detail-panel" v-else>
          <h3>任务详情</h3>
          <template v-if="distributionStep === 'inspection' && activeTaskDetail">
            <dl>
              <dt>病例</dt><dd>{{ activeTaskDetail.sample }}</dd>
              <dt>患者</dt><dd>{{ activeTaskDetail.patient }}</dd>
              <dt>年龄</dt><dd>{{ activeTaskDetail.age }} 岁</dd>
              <dt>检查类型</dt><dd>{{ activeTaskDetail.type }}</dd>
              <dt>采集时间</dt><dd>{{ activeTaskDetail.date }}</dd>
              <dt>影像数量</dt><dd>{{ activeTaskDetail.imageCount }} 张</dd>
            </dl>
            <button class="notify" @click="viewFullImage">查看完整影像</button>
          </template>
          <p v-else-if="distributionStep !== 'inspection'" class="empty">请先进入“影像数据检查阶段”，再查看任务详情。</p>
          <p v-else class="empty">点击“查看详情”后可在此查看任务详情。</p>
          <p v-if="formMessage" class="form-message">{{ formMessage }}</p>
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

const projectId = computed(() => String(route.query.projectId || 'XXXXXXXXXX'));
const flowKey = computed(() => `oct-hardware-flow-${projectId.value}`);
const stageLabel = computed(() => {
  if (stage.value === 'technician') return '技师认证';
  if (stage.value === 'certificate') return '证书颁发';
  return '硬件认证';
});
const pageTitle = computed(() => {
  if (stage.value === 'technician') return '项目管理-分发影像数据';
  return `项目管理-${stageLabel.value}`;
});
const distributionStepLabel = computed(() => {
  if (distributionStep.value === 'inspection') return '影像数据检查阶段';
  if (distributionStep.value === 'reading') return '阅片阶段';
  if (distributionStep.value === 'quality') return '质量抽查';
  return '受试者筛选阶段';
});
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
  if (nextStep !== 'inspection') {
    selectedTaskIds.value = [];
    activeTaskId.value = null;
    distributionMessage.value = '';
  }
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
  formMessage.value = '';
}

function distributeSelected(mode: 'batch' | 'smart') {
  if (distributionStep.value !== 'inspection') {
    distributionMessage.value = '请先进入影像数据检查阶段。';
    return;
  }
  if (!selectedTaskIds.value.length) return;
  distributionMessage.value = mode === 'batch'
    ? `已将 ${selectedTaskIds.value.length} 份影像资料批量分发给指定医生。`
    : `已将 ${selectedTaskIds.value.length} 份影像资料智能分发给适配医生。`;
  selectedTaskIds.value = [];
}

function viewFullImage() {
  if (distributionStep.value !== 'inspection') {
    formMessage.value = '请先进入影像数据检查阶段。';
    return;
  }
  if (!activeTaskDetail.value) {
    formMessage.value = '请先选择任务并查看详情。';
    return;
  }
  formMessage.value = `已打开 ${activeTaskDetail.value.sample} 的完整影像。`;
  showImagePreview.value = true;
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

<style scoped>
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
.task-actions { display: flex; gap: 10px; }
.task-actions button { border: 1px solid #c8d7ff; background: #e8f0ff; color: #1f3b8f; border-radius: 8px; padding: 7px 12px; cursor: pointer; }
.task-actions button:disabled { opacity: .5; cursor: not-allowed; }
.task-list { display: grid; gap: 10px; }
.task-card { border: 1px solid #d2dae6; border-radius: 8px; padding: 10px; display: flex; justify-content: space-between; align-items: center; }
.task-card label { display: flex; align-items: flex-start; gap: 10px; flex: 1; }
.task-main { display: grid; gap: 3px; }
.task-main small { color: #64748b; }
.detail-link { border: none; background: transparent; color: #2563eb; cursor: pointer; }
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
.task-detail-panel dl { display: grid; grid-template-columns: 70px 1fr; gap: 8px; margin: 0; }
.task-detail-panel dt { color: #64748b; }
.task-detail-panel dd { margin: 0; }
.stage-placeholder { border: 1px dashed #d2dae6; border-radius: 8px; padding: 14px; color: #334155; display: grid; gap: 10px; }
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
