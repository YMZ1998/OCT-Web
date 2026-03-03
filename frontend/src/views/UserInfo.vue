<template>
  <div class="user-info-wrapper">
    <div class="user-card blue-theme">
      <div class="avatar">
        <img :src="`https://api.dicebear.com/7.x/initials/svg?seed=${user?.username || 'U'}`" alt="avatar" />
      </div>
      <div class="info">
        <h2>{{ user?.username || '未登录' }}</h2>
        <ul>
          <li>
            <span class="label">邮箱：</span>
            <template v-if="editing">
              <input v-model="form.email" type="email" placeholder="请输入邮箱" />
            </template>
            <template v-else>{{ user?.email || '未填写' }}</template>
          </li>
          <li>
            <span class="label">性别：</span>
            <template v-if="editing">
              <select v-model="form.gender">
                <option value="">未填写</option>
                <option value="male">男</option>
                <option value="female">女</option>
              </select>
            </template>
            <template v-else>{{ user?.gender === 'male' ? '男' : user?.gender === 'female' ? '女' : (user?.gender || '未填写') }}</template>
          </li>
          <li>
            <span class="label">年龄：</span>
            <template v-if="editing">
              <input v-model.number="form.age" type="number" min="0" max="200" placeholder="请输入年龄" />
            </template>
            <template v-else>{{ user?.age ?? '未填写' }}</template>
          </li>
          <li><span class="label">注册时间：</span>{{ user?.created_at ? formatDate(user.created_at) : '未填写' }}</li>
          <li><span class="label">最近登录：</span>{{ user?.last_login_at ? formatDateTime(user.last_login_at) : '首次登录' }}</li>
        </ul>
      </div>
      <p v-if="message" class="message">{{ message }}</p>
      <div class="actions">
        <button v-if="!editing" @click="startEdit">编辑资料</button>
        <button v-if="editing" @click="onSave" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
        <button v-if="editing" class="secondary" @click="cancelEdit" :disabled="saving">取消</button>
        <button class="danger" @click="onLogout">退出登录</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { getUser, updateUser } from '../api/user';
import type { User } from '../types/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const user = ref<User | null>(null);
const editing = ref(false);
const saving = ref(false);
const message = ref('');
const form = ref({ email: '', gender: '', age: undefined as number | undefined });

function formatDate(ts: number) {
  const d = new Date(ts * 1000);
  return d.toLocaleDateString();
}

function formatDateTime(ts: number) {
  const d = new Date(ts * 1000);
  return d.toLocaleString();
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

function startEdit() {
  fillFormFromUser();
  message.value = '';
  editing.value = true;
}

function cancelEdit() {
  editing.value = false;
  message.value = '';
  fillFormFromUser();
}

async function onSave() {
  if (!userStore.token || !user.value) return;
  saving.value = true;
  message.value = '';
  try {
    const res = await updateUser(user.value.id, {
      email: form.value.email,
      gender: form.value.gender,
      age: form.value.age,
    }, userStore.token);
    user.value = res.data.data;
    userStore.setUserInfo(res.data.data);
    editing.value = false;
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
  max-width: 460px;
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
.info { width: 100%; }
.info h2 {
  text-align: center;
  margin-bottom: 20px;
  font-size: 1.7em;
  font-weight: 700;
  color: #409eff;
}
.info ul { list-style: none; padding: 0; margin: 0; }
.info li {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  font-size: 1.08em;
  color: #333;
}
.label { width: 90px; color: #409eff; font-weight: 600; }
.info input, .info select {
  flex: 1;
  border: 1px solid #bcdcff;
  border-radius: 6px;
  padding: 6px 8px;
}
.message {
  margin-top: 8px;
  color: #409eff;
  font-size: 0.95em;
}
.actions {
  margin-top: 20px;
  width: 100%;
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}
.actions button {
  background: linear-gradient(90deg, #409eff 60%, #66b1ff 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  font-size: 1em;
  cursor: pointer;
  font-weight: 600;
}
.actions button.secondary { background: #909399; }
.actions button.danger { background: #f56c6c; }
.actions button:disabled { opacity: 0.7; cursor: not-allowed; }
</style>
