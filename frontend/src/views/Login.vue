<template>
  <div class="login-page">
    <section class="brand-panel">
      <div class="brand-content">
        <h1>SJTURC眼科</h1>
        <p>在线工作管理系统</p>
        <div class="eye-illustration" aria-hidden="true">
          <div class="eye-ring outer"></div>
          <div class="eye-ring mid"></div>
          <div class="eye-ring inner"></div>
          <div class="pupil"></div>
        </div>
      </div>
      <footer>© 2025 SJTURC眼科研究中心. 保留所有权利.</footer>
    </section>

    <section class="form-panel">
      <div class="login-shell">
        <h2>欢迎登录</h2>

        <form @submit.prevent="onLogin">
          <label for="username">用户名</label>
          <div class="field">
            <span class="icon">👤</span>
            <input id="username" v-model="username" placeholder="请输入用户名" required />
          </div>

          <label for="password">密码</label>
          <div class="field">
            <span class="icon">🔒</span>
            <input id="password" v-model="password" :type="showPassword ? 'text' : 'password'" placeholder="请输入密码" required />
            <button class="toggle" type="button" @click="showPassword = !showPassword">
              {{ showPassword ? '🙈' : '👁️' }}
            </button>
          </div>

          <button type="submit" class="submit-btn">登录</button>
          <div v-if="error" class="error">{{ error }}</div>
        </form>

        <p class="support">遇到登录问题? 请联系 <a href="#" @click.prevent>技术支持</a></p>
      </div>
    </section>

    <div v-if="showPasswordErrorModal" class="modal-mask" @click.self="closePasswordErrorModal">
      <div class="theme-modal" role="dialog" aria-modal="true" aria-labelledby="password-error-title">
        <h3 id="password-error-title">登录失败</h3>
        <p>密码错误，请重新输入</p>
        <button class="modal-btn" @click="closePasswordErrorModal">我知道了</button>
      </div>
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
const showPassword = ref(false);
const remember = ref(false);
const showPasswordErrorModal = ref(false);
const router = useRouter();
const userStore = useUserStore();

function closePasswordErrorModal() {
  showPasswordErrorModal.value = false;
}

async function onLogin() {
  error.value = '';
  showPasswordErrorModal.value = false;

  try {
    const res = await login({ username: username.value, password: password.value });
    const { code, msg, data } = res.data || {};

    if (code !== 0) {
      if (msg === '登录失败') {
        error.value = '密码错误，请重新输入';
        showPasswordErrorModal.value = true;
      } else {
        error.value = msg || '登录失败';
      }
      return;
    }

    if (!data?.token || !data?.id) {
      error.value = '登录返回数据异常';
      return;
    }

    userStore.setToken(data.token);
    router.push(`/user/${data.id}`);
  } catch (e: any) {
    error.value = e.response?.data?.msg || '登录失败';
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1fr 1.1fr;
  background: #eef2f6;
}
.brand-panel {
  background: linear-gradient(160deg, #045ea7, #0a4f96);
  color: #fff;
  padding: 46px 44px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.brand-content h1 {
  margin: 0;
  font-size: 52px;
}
.brand-content p {
  margin: 12px 0 0;
  font-size: 34px;
  opacity: 0.96;
}
.eye-illustration {
  margin: 60px auto 0;
  width: 320px;
  height: 320px;
  position: relative;
}
.eye-ring {
  position: absolute;
  border-radius: 50%;
  inset: 0;
  margin: auto;
}
.eye-ring.outer {
  width: 320px;
  height: 320px;
  background: radial-gradient(circle, rgba(124, 201, 255, 0.45), rgba(124, 201, 255, 0.18));
}
.eye-ring.mid {
  width: 260px;
  height: 260px;
  background: radial-gradient(circle, #74a9ff, #4176df);
}
.eye-ring.inner {
  width: 88px;
  height: 88px;
  background: #031329;
  top: 116px;
}
.pupil {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.78);
  position: absolute;
  top: 126px;
  left: 122px;
}
.brand-panel footer {
  font-size: 20px;
  opacity: 0.9;
}
.form-panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}
.login-shell {
  width: min(620px, 100%);
}
.badge {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #d6ebfb;
  color: #177ac2;
  display: grid;
  place-items: center;
  margin: 0 auto 22px;
}
h2 {
  margin: 0;
  text-align: center;
  font-size: 50px;
  color: #10203c;
}
.subtitle {
  text-align: center;
  color: #3a4a62;
  margin: 10px 0 38px;
  font-size: 30px;
}
form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
label {
  font-size: 30px;
  color: #1d2d47;
}
.field {
  display: flex;
  align-items: center;
  background: #fff;
  border: 1px solid #c9d3df;
  border-radius: 10px;
  padding: 0 14px;
}
.field .icon { margin-right: 10px; opacity: 0.7; }
input {
  border: none;
  flex: 1;
  height: 54px;
  outline: none;
  font-size: 26px;
  background: transparent;
}
.toggle {
  border: none;
  background: transparent;
  cursor: pointer;
}
.helper-row {
  margin-top: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #3f4f68;
  font-size: 24px;
}
.helper-row a { color: #157ec6; text-decoration: none; }
.remember { display: flex; align-items: center; gap: 6px; }
.submit-btn {
  margin-top: 14px;
  height: 56px;
  border: none;
  border-radius: 10px;
  background: #077cbc;
  color: #fff;
  font-size: 30px;
  cursor: pointer;
}
.error { color: #d93c3c; margin-top: 6px; font-size: 24px; }
.support {
  text-align: center;
  margin-top: 24px;
  color: #3f4f68;
  font-size: 24px;
}
.support a { color: #157ec6; text-decoration: none; }

.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(64, 158, 255, 0.18);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.theme-modal {
  width: min(88vw, 360px);
  background: #fff;
  border: 2px solid #409eff;
  border-radius: 14px;
  box-shadow: 0 10px 28px rgba(64, 158, 255, 0.28);
  padding: 22px 20px 18px;
  text-align: center;
}
.theme-modal h3 {
  margin: 0 0 10px;
  font-size: 1.2rem;
  color: #337ecc;
}
.theme-modal p { margin: 0; color: #4a5d76; }
.modal-btn {
  margin-top: 18px;
  width: 100%;
  border: none;
  border-radius: 8px;
  background: #409eff;
  color: #fff;
  padding: 10px 0;
}

@media (max-width: 960px) {
  .login-page { grid-template-columns: 1fr; }
  .brand-panel { display: none; }
  h2 { font-size: 34px; }
  .subtitle, label, input, .helper-row, .submit-btn, .support, .error { font-size: 18px; }
}
</style>
