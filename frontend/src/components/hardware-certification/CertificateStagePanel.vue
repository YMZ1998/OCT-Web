<template>
  <div>
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
            <label><input type="checkbox" :checked="useSeal" @change="onSealChange" /> 添加印章</label>
            <label><input type="checkbox" :checked="useSign" @change="onSignChange" /> 添加电子签名</label>
          </div>
          <div class="marks">
            <span v-if="useSeal" class="seal">SJTURC印章</span>
            <span v-if="useSign" class="sign">授权签名</span>
          </div>
        </div>
      </section>
    </article>

    <div class="certificate-actions">
      <button class="issue" :disabled="certificateIssued" @click="$emit('issue')">{{ certificateIssued ? '已颁发' : '颁发证书' }}</button>
      <div>
        <button class="secondary" @click="$emit('print')">打印证书</button>
        <button class="secondary" @click="$emit('download')">下载证书</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  projectId: string;
  issueDate: string;
  expiryDate: string;
  certificateNo: string;
  qrCells: boolean[];
  useSeal: boolean;
  useSign: boolean;
  certificateIssued: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:seal', value: boolean): void;
  (e: 'update:sign', value: boolean): void;
  (e: 'issue'): void;
  (e: 'print'): void;
  (e: 'download'): void;
}>();

function onSealChange(event: Event) {
  emit('update:seal', (event.target as HTMLInputElement).checked);
}

function onSignChange(event: Event) {
  emit('update:sign', (event.target as HTMLInputElement).checked);
}
</script>
