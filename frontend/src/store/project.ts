import { defineStore } from 'pinia';

export type ProjectItem = {
  id: number;
  name: string;
  owner: string;
  center: string;
  state: '进行中' | '待开始' | '已完成';
  stateClass: 'running' | 'pending' | 'completed';
  desc: string;
  date: string;
  members: number;
};

export type TodoItem = {
  key: string;
  projectId: number;
  projectName: string;
  taskName: string;
};

type ProjectState = {
  recentProjects: ProjectItem[];
  todoItems: TodoItem[];
  nextProjectId: number;
};

const STORAGE_KEY = 'oct-project-state';

const STANDARD_TASKS = ['认证', '分发影像数据', '阅片审核'] as const;

function getProjectStateByTodos(hasStarted: boolean, todoCount: number): Pick<ProjectItem, 'state' | 'stateClass'> {
  if (!hasStarted) return { state: '待开始', stateClass: 'pending' };
  if (todoCount === 0) return { state: '已完成', stateClass: 'completed' };
  return { state: '进行中', stateClass: 'running' };
}

function defaultState(): ProjectState {
  return {
    recentProjects: [
      {
        id: 1,
        name: '青光眼早期诊断研究',
        owner: '张医生',
        center: '上海眼科中心',
        state: '进行中',
        stateClass: 'running',
        desc: '研究青光眼早期诊断的新方法和技术',
        date: '2023-06-15',
        members: 8,
      },
      {
        id: 2,
        name: '视网膜病变图像分析',
        owner: '李医生',
        center: '虹桥临床中心',
        state: '待开始',
        stateClass: 'pending',
        desc: '构建视网膜病变分级模型并优化阅片流程',
        date: '2023-07-01',
        members: 12,
      },
    ],
    todoItems: [
      { key: '1-认证', projectId: 1, projectName: '青光眼早期诊断研究', taskName: '认证' },
      { key: '1-分发影像数据', projectId: 1, projectName: '青光眼早期诊断研究', taskName: '分发影像数据' },
      { key: '2-阅片审核', projectId: 2, projectName: '视网膜病变图像分析', taskName: '阅片审核' },
    ],
    nextProjectId: 3,
  };
}

function loadState(): ProjectState {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return defaultState();

  try {
    const parsed = JSON.parse(raw) as ProjectState;
    if (!Array.isArray(parsed.recentProjects) || !Array.isArray(parsed.todoItems)) {
      return defaultState();
    }
    return parsed;
  } catch {
    return defaultState();
  }
}

export const useProjectStore = defineStore('project', {
  state: () => loadState(),
  getters: {
    pendingCount: (state) => state.todoItems.length,
  },
  actions: {
    persist() {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({
          recentProjects: this.recentProjects,
          todoItems: this.todoItems,
          nextProjectId: this.nextProjectId,
        }),
      );
    },
    addProject(payload: Omit<ProjectItem, 'id' | 'state' | 'stateClass'>) {
      const project: ProjectItem = {
        id: this.nextProjectId,
        name: payload.name,
        owner: payload.owner,
        center: payload.center,
        desc: payload.desc,
        date: payload.date,
        members: payload.members,
        state: '待开始',
        stateClass: 'pending',
      };
      this.recentProjects.unshift(project);
      this.nextProjectId += 1;
      this.persist();
    },

    syncProjectState(projectId: number) {
      const project = this.recentProjects.find((item) => item.id === projectId);
      if (!project) return;

      const todoCount = this.todoItems.filter((item) => item.projectId === projectId).length;
      const hasStarted = project.state !== '待开始' || todoCount > 0;
      const next = getProjectStateByTodos(hasStarted, todoCount);
      project.state = next.state;
      project.stateClass = next.stateClass;
    },

    startProjectTasks(project: ProjectItem) {
      STANDARD_TASKS.forEach((taskName) => {
        const key = `${project.id}-${taskName}`;
        if (this.todoItems.some((item) => item.key === key)) return;
        this.todoItems.push({
          key,
          projectId: project.id,
          projectName: project.name,
          taskName,
        });
      });
      this.syncProjectState(project.id);
      this.persist();
    },

    completeTask(taskKey: string) {
      const idx = this.todoItems.findIndex((item) => item.key === taskKey);
      if (idx < 0) return;

      const task = this.todoItems[idx];
      this.todoItems.splice(idx, 1);
      this.syncProjectState(task.projectId);
      this.persist();
    },

    updateProject(projectId: number, payload: Omit<ProjectItem, 'id' | 'state' | 'stateClass'>) {
      const idx = this.recentProjects.findIndex((item) => item.id === projectId);
      if (idx < 0) return;

      const current = this.recentProjects[idx];
      const updated: ProjectItem = {
        ...current,
        name: payload.name,
        owner: payload.owner,
        center: payload.center,
        desc: payload.desc,
        date: payload.date,
        members: payload.members,
      };
      this.recentProjects.splice(idx, 1, updated);

      this.todoItems = this.todoItems.map((todo) =>
        todo.projectId === projectId
          ? {
              ...todo,
              projectName: updated.name,
            }
          : todo,
      );

      this.persist();
    },

    addTask(project: ProjectItem, taskName: string, dedupeKey?: string) {
      const key = dedupeKey || `${project.id}-${taskName}`;
      if (this.todoItems.some((item) => item.key === key)) return;

      this.todoItems.unshift({
        key,
        projectId: project.id,
        projectName: project.name,
        taskName,
      });
      this.syncProjectState(project.id);
      this.persist();
    },
  },
});
