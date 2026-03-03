import axios from 'axios';

export function register(data: { username: string; password: string; email: string }) {
  return axios.post('/api/v1/register', data);
}

export function login(data: { username: string; password: string }) {
  return axios.post('/api/v1/login', data);
}

export function getUser(id: string, token: string) {
  return axios.get(`/api/v1/user/${id}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function updateUser(id: string, data: any, token: string) {
  return axios.put(`/api/v1/user/${id}`, data, {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function deleteUser(id: string, token: string) {
  return axios.delete(`/api/v1/user/${id}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
}

export function updateUserProjects(id: string, projectState: any, token: string) {
  return axios.put(`/api/v1/user/${id}/projects`, { project_state: projectState }, {
    headers: { Authorization: `Bearer ${token}` }
  });
}
