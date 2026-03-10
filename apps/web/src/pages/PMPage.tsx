import { FormEvent, useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { apiClient } from '../api/client'
import { useAuth } from '../context/AuthContext'

type Project = {
  id: number
  name: string
  description: string
  ownerId: number
}

type TaskStatus = 'todo' | 'in_progress' | 'in_review' | 'done'

type Task = {
  id: number
  projectId: number
  title: string
  description: string
  status: TaskStatus
}

type ListResponse<T> = {
  items: T[]
}

const statuses: TaskStatus[] = ['todo', 'in_progress', 'in_review', 'done']

export function PMPage() {
  const { token } = useAuth()
  const queryClient = useQueryClient()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [projectName, setProjectName] = useState('')
  const [projectDescription, setProjectDescription] = useState('')
  const [taskTitle, setTaskTitle] = useState('')
  const [taskDescription, setTaskDescription] = useState('')

  const projectsQuery = useQuery({
    queryKey: ['projects'],
    queryFn: () => apiClient.get<ListResponse<Project>>('/api/pm/projects', token ?? undefined),
    enabled: Boolean(token),
  })

  const tasksQuery = useQuery({
    queryKey: ['tasks', selectedProjectId],
    queryFn: () =>
      apiClient.get<ListResponse<Task>>(
        `/api/pm/tasks?projectId=${selectedProjectId ?? ''}`,
        token ?? undefined,
      ),
    enabled: Boolean(token) && selectedProjectId !== null,
  })

  useEffect(() => {
    if (!selectedProjectId && projectsQuery.data?.items?.length) {
      setSelectedProjectId(projectsQuery.data.items[0].id)
    }
  }, [projectsQuery.data, selectedProjectId])

  const createProjectMutation = useMutation({
    mutationFn: (payload: { name: string; description: string }) =>
      apiClient.post<Project>('/api/pm/projects', payload, token ?? undefined),
    onSuccess: async (created) => {
      setProjectName('')
      setProjectDescription('')
      await queryClient.invalidateQueries({ queryKey: ['projects'] })
      setSelectedProjectId(created.id)
    },
  })

  const createTaskMutation = useMutation({
    mutationFn: (payload: { projectId: number; title: string; description: string }) =>
      apiClient.post<Task>('/api/pm/tasks', payload, token ?? undefined),
    onSuccess: async () => {
      setTaskTitle('')
      setTaskDescription('')
      // partial refresh: only refetch current tasks list
      await queryClient.invalidateQueries({ queryKey: ['tasks', selectedProjectId] })
    },
  })

  const patchTaskStatusMutation = useMutation({
    mutationFn: (payload: { id: number; status: TaskStatus }) =>
      apiClient.patch<Task>(`/api/pm/tasks/${payload.id}/status`, { status: payload.status }, token ?? undefined),
    onSuccess: async () => {
      // partial refresh: invalidate only current task board data
      await queryClient.invalidateQueries({ queryKey: ['tasks', selectedProjectId] })
    },
  })

  const board = useMemo(() => {
    const buckets: Record<TaskStatus, Task[]> = {
      todo: [],
      in_progress: [],
      in_review: [],
      done: [],
    }

    for (const task of tasksQuery.data?.items ?? []) {
      buckets[task.status].push(task)
    }

    return buckets
  }, [tasksQuery.data])

  function onCreateProject(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const name = projectName.trim()
    if (!name) {
      return
    }

    createProjectMutation.mutate({
      name,
      description: projectDescription.trim(),
    })
  }

  function onCreateTask(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!selectedProjectId) {
      return
    }

    const title = taskTitle.trim()
    if (!title) {
      return
    }

    createTaskMutation.mutate({
      projectId: selectedProjectId,
      title,
      description: taskDescription.trim(),
    })
  }

  function nextStatus(status: TaskStatus): TaskStatus | null {
    const idx = statuses.indexOf(status)
    if (idx < 0 || idx === statuses.length - 1) {
      return null
    }
    return statuses[idx + 1]
  }

  return (
    <section className="stack-lg">
      <header className="page-header">
        <div>
          <h1>Project Management</h1>
          <p className="muted">Create projects, track tasks, and progress work quickly.</p>
        </div>
      </header>

      <div className="panel-grid">
        <article className="panel stack-md">
          <h2>Projects</h2>
          <form onSubmit={onCreateProject} className="stack-sm">
            <input
              placeholder="New project name"
              value={projectName}
              onChange={(event) => setProjectName(event.target.value)}
            />
            <input
              placeholder="Description"
              value={projectDescription}
              onChange={(event) => setProjectDescription(event.target.value)}
            />
            <button type="submit" className="btn-primary" disabled={createProjectMutation.isPending}>
              {createProjectMutation.isPending ? 'Creating...' : 'Create Project'}
            </button>
          </form>
          <div className="stack-sm">
            {(projectsQuery.data?.items ?? []).map((project) => (
              <button
                key={project.id}
                type="button"
                className={selectedProjectId === project.id ? 'project-item active' : 'project-item'}
                onClick={() => setSelectedProjectId(project.id)}
              >
                <span>{project.name}</span>
                <small className="muted">#{project.id}</small>
              </button>
            ))}
          </div>
        </article>

        <article className="panel stack-md">
          <h2>Tasks</h2>
          <form onSubmit={onCreateTask} className="stack-sm">
            <input
              placeholder="Task title"
              value={taskTitle}
              onChange={(event) => setTaskTitle(event.target.value)}
              disabled={!selectedProjectId}
            />
            <input
              placeholder="Task description"
              value={taskDescription}
              onChange={(event) => setTaskDescription(event.target.value)}
              disabled={!selectedProjectId}
            />
            <button
              type="submit"
              className="btn-primary"
              disabled={!selectedProjectId || createTaskMutation.isPending}
            >
              {createTaskMutation.isPending ? 'Creating...' : 'Create Task'}
            </button>
          </form>
        </article>
      </div>

      <section className="kanban-grid">
        {statuses.map((status) => (
          <article key={status} className="kanban-column">
            <h3>{status.replace('_', ' ')}</h3>
            <div className="stack-sm">
              {board[status].map((task) => {
                const next = nextStatus(task.status)

                return (
                  <div key={task.id} className="task-card">
                    <strong>{task.title}</strong>
                    <p className="muted">{task.description || 'No description'}</p>
                    {next ? (
                      <button
                        type="button"
                        className="btn-secondary"
                        onClick={() => patchTaskStatusMutation.mutate({ id: task.id, status: next })}
                        disabled={patchTaskStatusMutation.isPending}
                      >
                        Move to {next.replace('_', ' ')}
                      </button>
                    ) : (
                      <span className="done-badge">Completed</span>
                    )}
                  </div>
                )
              })}
              {board[status].length === 0 && <p className="muted">No tasks</p>}
            </div>
          </article>
        ))}
      </section>
    </section>
  )
}
