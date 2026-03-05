<template>
  <div>
    <template v-if="distributionStep === 'inspection'">
      <section class="task-header">
        <p class="sub">任务列表</p>
        <div class="task-actions">
          <button :disabled="!selectedTaskIds.length" @click="$emit('distribute', 'batch')">批量分发</button>
          <button :disabled="!selectedTaskIds.length" @click="$emit('distribute', 'smart')">智能分发</button>
        </div>
      </section>

      <div class="task-list">
        <article class="task-card" v-for="item in pagedImages" :key="item.id">
          <label>
            <input type="checkbox" :checked="selectedTaskIds.includes(item.id)" @change="$emit('toggle-selection', item.id)" />
            <div class="task-main">
              <strong>{{ item.sample }}</strong>
              <small>患者：{{ item.patient }}｜年龄：{{ item.age }}岁｜检查类型：{{ item.type }}</small>
              <small>{{ item.date }} · {{ item.imageCount }}张影像</small>
            </div>
          </label>
          <button class="detail-link" @click="$emit('show-task-detail', item.id)">查看详情</button>
        </article>
      </div>

      <p v-if="distributionMessage" class="distribution-message">{{ distributionMessage }}</p>
    </template>

    <template v-else>
      <section class="task-header">
        <p class="sub">影像任务</p>
      </section>

      <div class="tech-grid">
        <article
          class="tech-card"
          v-for="item in pagedImages"
          :key="item.id"
          :class="{ active: activeTaskId === item.id }"
          @click="$emit('show-task-detail', item.id)"
        >
          <div class="tech-thumb">🖼️</div>
          <div class="tech-meta">
            <strong>{{ item.sample }}</strong>
            <small>眼底照片</small>
          </div>
        </article>
      </div>
    </template>

    <div class="pager" role="navigation" aria-label="影像分页">
      <button :disabled="currentPage === 1" @click="$emit('go-page', currentPage - 1)">‹</button>
      <button
        v-for="page in imageTotalPages"
        :key="page"
        :class="['page-btn', page === currentPage ? 'active' : '']"
        @click="$emit('go-page', page)"
      >
        {{ page }}
      </button>
      <button :disabled="currentPage === imageTotalPages" @click="$emit('go-page', currentPage + 1)">›</button>
    </div>
  </div>
</template>

<script setup lang="ts">
type DistributionStep = 'screening' | 'inspection' | 'reading' | 'quality';

type ImageItem = {
  id: number;
  sample: string;
  patient: string;
  age: number;
  date: string;
  imageCount: number;
  type: string;
};

defineProps<{
  distributionStep: DistributionStep;
  selectedTaskIds: number[];
  pagedImages: ImageItem[];
  distributionMessage: string;
  activeTaskId: number | null;
  currentPage: number;
  imageTotalPages: number;
}>();

defineEmits<{
  (e: 'distribute', mode: 'batch' | 'smart'): void;
  (e: 'toggle-selection', id: number): void;
  (e: 'show-task-detail', id: number): void;
  (e: 'go-page', page: number): void;
}>();
</script>
