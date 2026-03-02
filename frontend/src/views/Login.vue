<template>
  <div class="login-wrapper">
    <div class="login-card blue-theme">
      <h2>用户登录</h2>
      <form @submit.prevent="onLogin">
        <input v-model="username" placeholder="用户名" required />
        <input v-model="password" type="password" placeholder="密码" required />
        <button type="submit">登录</button>
      </form>
      <button class="register-btn" @click="goRegister">注册新用户</button>
      <div v-if="error" class="error">{{ error }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { login } from '../api/user';

const username = ref('');
const password = ref('');
const error = ref('');
const router = useRouter();
const userStore = useUserStore();

function goRegister() {
  router.push('/register');
}

async function onLogin() {
  error.value = '';
  try {
    const res = await login({ username: username.value, password: password.value });
    userStore.setToken(res.data.data.token);
    router.push(`/user/${res.data.data.id}`);
  } catch (e: any) {
    error.value = e.response?.data?.msg || '登录失败';
  }
}
</script>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  background: linear-gradient(135deg, #e3f0ff 0%, #f6f8fa 100%);
}
.login-card.blue-theme {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(64,158,255,0.12);
  padding: 40px 48px;
  min-width: 340px;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 2px solid #409eff;
}
.login-card h2 {
  color: #409eff;
  font-size: 1.6em;
  font-weight: 700;
  margin-bottom: 24px;
}
form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 18px;
}
input {
  padding: 10px 14px;
  border: 1px solid #bcdcff;
  border-radius: 8px;
  font-size: 1em;
  outline: none;
  transition: border 0.2s;
}
input:focus {
  border-color: #409eff;
}
button {
  background: linear-gradient(90deg, #409eff 60%, #66b1ff 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 0;
  font-size: 1.08em;
  cursor: pointer;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(64,158,255,0.08);
  transition: background 0.2s;
}
button:hover {
  background: linear-gradient(90deg, #337ecc 60%, #409eff 100%);
}
.register-btn {
  margin-top: 18px;
  width: 100%;
  background: #fff;
  color: #409eff;
  border: 1px solid #409eff;
  border-radius: 8px;
  padding: 10px 0;
  font-size: 1.08em;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}
.register-btn:hover {
  background: #409eff;
  color: #fff;
}
.error {
  color: #e74c3c;
  margin-top: 16px;
  font-size: 1em;
  text-align: center;
}
</style>
