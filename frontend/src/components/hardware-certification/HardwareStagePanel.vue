<template>
  <div>
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
</template>

<script setup lang="ts">
type HardwareDoc = {
  name: string;
  size: string;
  date: string;
};

defineProps<{
  docs: HardwareDoc[];
}>();
</script>
