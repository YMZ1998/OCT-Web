<template>
  <div class="user-info-wrapper">
    <div class="user-card blue-theme">
      <div class="avatar">
        <img :src="`https://api.dicebear.com/7.x/initials/svg?seed=${user?.username || 'U'}`" alt="avatar" />
      </div>
      <div class="info">
        <h2>{{ user?.username || '未登录' }}</h2>
        <ul>
          <li><span class="label">邮箱：</span>{{ user?.email || '未填写' }}</li>
          <li><span class="label">性别：</span>{{ user?.gender === 'male' ? '男' : user?.gender === 'female' ? '女' : (user?.gender || '未填写') }}</li>
          <li><span class="label">年龄：</span>{{ user?.age ?? '未填写' }}</li>
          <li><span class="label">注册时间：</span>{{ user?.created_at ? formatDate(user.created_at) : '未填写' }}</li>
          <li><span class="label">最近登录：</span>{{ user?.last_login_at ? formatDateTime(user.last_login_at) : '首次登录' }}</li>
        </ul>
      </div>
      <div class="actions">
        <button @click="onLogout">退出登录</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { getUser } from '../api/user';
import type { User } from '../types/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);

function formatDate(ts: number) {
  const d = new Date(ts * 1000);
  return d.toLocaleDateString();
}

function formatDateTime(ts: number) {
  const d = new Date(ts * 1000);
  return d.toLocaleString();
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
.user-info-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  background: linear-gradient(135deg, #e3f0ff 0%, #f6f8fa 100%);
}
.user-card.blue-theme {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(64,158,255,0.12);
  padding: 40px 48px;
  min-width: 340px;
  max-width: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 2px solid #409eff;
}
.avatar {
  width: 90px;
  height: 90px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 18px;
  background: #e3f0ff;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #409eff;
}
.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.info {
  width: 100%;
}
.info h2 {
  text-align: center;
  margin-bottom: 20px;
  font-size: 1.7em;
  font-weight: 700;
  color: #409eff;
}
.info ul {
  list-style: none;
  padding: 0;
  margin: 0;
}
.info li {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  font-size: 1.08em;
  color: #333;
}
.label {
  width: 90px;
  color: #409eff;
  font-weight: 600;
}
.actions {
  margin-top: 28px;
  width: 100%;
  display: flex;
  justify-content: center;
}
.actions button {
  background: linear-gradient(90deg, #409eff 60%, #66b1ff 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 32px;
  font-size: 1.08em;
  cursor: pointer;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(64,158,255,0.08);
  transition: background 0.2s;
}
.actions button:hover {
  background: linear-gradient(90deg, #337ecc 60%, #409eff 100%);
}
</style>
