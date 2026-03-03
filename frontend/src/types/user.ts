export interface ProjectItem {
  id: number;
  name: string;
  owner: string;
  center: string;
  state: '进行中' | '待开始' | '已完成';
  stateClass: 'running' | 'pending' | 'completed';
  desc: string;
  date: string;
  members: number;
}

export interface TodoItem {
  key: string;
  projectId: number;
  projectName: string;
  taskName: string;
}

export interface ProjectState {
  recentProjects: ProjectItem[];
  todoItems: TodoItem[];
  nextProjectId: number;
}

export interface User {
  id: string;
  username: string;
  email: string;
  gender?: string;
  age?: number;
  created_at?: number;
  last_login_at?: number;
  project_state?: ProjectState;
}
