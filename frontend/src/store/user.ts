import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('oct-token') || '',
    userInfo: JSON.parse(localStorage.getItem('oct-user-info') || 'null') as any,
  }),
  actions: {
    setToken(token: string) {
      this.token = token;
      localStorage.setItem('oct-token', token);
    },
    setUserInfo(info: any) {
      this.userInfo = info;
      localStorage.setItem('oct-user-info', JSON.stringify(info));
    },
    logout() {
      this.token = '';
      this.userInfo = null;
      localStorage.removeItem('oct-token');
      localStorage.removeItem('oct-user-info');
    },
  },
});
