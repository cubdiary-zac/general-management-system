import { FormEvent, KeyboardEvent, useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { apiClient } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'
import { nextTaskStatus, orderedStatuses, TaskStatus } from './pm-utils'

type Project = {
  id: number
  name: string
  description: string
  ownerId: number
}

type Task = {
  id: number
  projectId: number
  title: string
  description: string
  status: TaskStatus
}

type TaskTransitionLog = {
  id: number
  taskId: number
  fromStatus: TaskStatus
  toStatus: TaskStatus
  operatorId: number
  createdAt: string
}

type ListResponse<T> = {
  items: T[]
}

type TaskStatusFilter = TaskStatus | 'all'

function formatLogTimestamp(value: string, locale: string): string {
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) {
    return value
  }

  return parsed.toLocaleString(locale)
}

export function PMPage() {
  const { token } = useAuth()
  const { locale, t } = useI18n()
  const queryClient = useQueryClient()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [projectName, setProjectName] = useState('')
  const [projectDescription, setProjectDescription] = useState('')
  const [taskTitle, setTaskTitle] = useState('')
  const [taskDescription, setTaskDescription] = useState('')
  const [statusFilter, setStatusFilter] = useState<TaskStatusFilter>('all')
  const [keyword, setKeyword] = useState('')
  const [selectedTaskId, setSelectedTaskId] = useState<number | null>(null)

  const trimmedKeyword = keyword.trim()
  const taskListQueryKey = ['tasks', selectedProjectId, statusFilter, trimmedKeyword] as const

  function taskStatusLabel(status: TaskStatus): string {
    return t(`status.task.${status}`)
  }

  const projectsQuery = useQuery({
    queryKey: ['projects'],
    queryFn: () => apiClient.get<ListResponse<Project>>('/api/pm/projects', token ?? undefined),
    enabled: Boolean(token),
  })

  const tasksQuery = useQuery({
    queryKey: taskListQueryKey,
    queryFn: () => {
      const params = new URLSearchParams()
      if (selectedProjectId !== null) {
        params.set('projectId', String(selectedProjectId))
      }
      if (statusFilter !== 'all') {
        params.set('status', statusFilter)
      }
      if (trimmedKeyword !== '') {
        params.set('q', trimmedKeyword)
      }

      return apiClient.get<ListResponse<Task>>(`/api/pm/tasks?${params.toString()}`, token ?? undefined)
    },
    enabled: Boolean(token) && selectedProjectId !== null,
  })

  const taskDetailQuery = useQuery({
    queryKey: ['task', selectedTaskId],
    queryFn: () => apiClient.get<Task>(`/api/pm/tasks/${selectedTaskId}`, token ?? undefined),
    enabled: Boolean(token) && selectedTaskId !== null,
  })

  const taskLogsQuery = useQuery({
    queryKey: ['task-logs', selectedTaskId],
    queryFn: () =>
      apiClient.get<ListResponse<TaskTransitionLog>>(`/api/pm/tasks/${selectedTaskId}/logs`, token ?? undefined),
    enabled: Boolean(token) && selectedTaskId !== null,
  })

  useEffect(() => {
    if (!selectedProjectId && projectsQuery.data?.items?.length) {
      setSelectedProjectId(projectsQuery.data.items[0].id)
    }
  }, [projectsQuery.data, selectedProjectId])

  useEffect(() => {
    setSelectedTaskId(null)
  }, [selectedProjectId])

  useEffect(() => {
    if (!tasksQuery.data || selectedTaskId === null) {
      return
    }

    const stillVisible = tasksQuery.data.items.some((task) => task.id === selectedTaskId)
    if (!stillVisible) {
      setSelectedTaskId(null)
    }
  }, [tasksQuery.data, selectedTaskId])

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
      await queryClient.invalidateQueries({ queryKey: taskListQueryKey })
    },
  })

  const patchTaskStatusMutation = useMutation({
    mutationFn: (payload: { id: number; status: TaskStatus }) =>
      apiClient.patch<Task>(`/api/pm/tasks/${payload.id}/status`, { status: payload.status }, token ?? undefined),
    onSuccess: async (_, payload) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: taskListQueryKey }),
        queryClient.invalidateQueries({ queryKey: ['task', payload.id] }),
        queryClient.invalidateQueries({ queryKey: ['task-logs', payload.id] }),
      ])
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

  function onTaskCardKeyDown(event: KeyboardEvent<HTMLDivElement>, taskId: number) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      setSelectedTaskId(taskId)
    }
  }

  return (
    <section className="stack-lg">
      <header className="page-header">
        <div>
          <h1>{t('pm.title')}</h1>
          <p className="muted">{t('pm.subtitle')}</p>
        </div>
      </header>

      <div className="panel-grid">
        <article className="panel stack-md">
          <h2>{t('pm.projects')}</h2>
          <form onSubmit={onCreateProject} className="stack-sm">
            <input
              placeholder={t('pm.newProjectName')}
              value={projectName}
              onChange={(event) => setProjectName(event.target.value)}
            />
            <input
              placeholder={t('common.description')}
              value={projectDescription}
              onChange={(event) => setProjectDescription(event.target.value)}
            />
            <button type="submit" className="btn-primary" disabled={createProjectMutation.isPending}>
              {createProjectMutation.isPending ? t('pm.creatingProject') : t('pm.createProject')}
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
          <h2>{t('pm.tasks')}</h2>
          <form onSubmit={onCreateTask} className="stack-sm">
            <input
              placeholder={t('pm.taskTitle')}
              value={taskTitle}
              onChange={(event) => setTaskTitle(event.target.value)}
              disabled={!selectedProjectId}
            />
            <input
              placeholder={t('pm.taskDescription')}
              value={taskDescription}
              onChange={(event) => setTaskDescription(event.target.value)}
              disabled={!selectedProjectId}
            />
            <button
              type="submit"
              className="btn-primary"
              disabled={!selectedProjectId || createTaskMutation.isPending}
            >
              {createTaskMutation.isPending ? t('pm.creatingTask') : t('pm.createTask')}
            </button>
          </form>

          <div className="task-filters">
            <select
              value={statusFilter}
              onChange={(event) => setStatusFilter(event.target.value as TaskStatusFilter)}
              disabled={!selectedProjectId}
            >
              <option value="all">{t('common.allStatuses')}</option>
              {orderedStatuses.map((status) => (
                <option key={status} value={status}>
                  {taskStatusLabel(status)}
                </option>
              ))}
            </select>
            <input
              placeholder={t('pm.searchTaskByTitle')}
              value={keyword}
              onChange={(event) => setKeyword(event.target.value)}
              disabled={!selectedProjectId}
            />
          </div>
        </article>
      </div>

      <section className="kanban-grid">
        {orderedStatuses.map((status) => (
          <article key={status} className="kanban-column">
            <h3>{taskStatusLabel(status)}</h3>
            <div className="stack-sm">
              {board[status].map((task) => {
                const next = nextTaskStatus(task.status)

                return (
                  <div
                    key={task.id}
                    className={selectedTaskId === task.id ? 'task-card clickable active' : 'task-card clickable'}
                    onClick={() => setSelectedTaskId(task.id)}
                    onKeyDown={(event) => onTaskCardKeyDown(event, task.id)}
                    role="button"
                    tabIndex={0}
                  >
                    <strong>{task.title}</strong>
                    <p className="muted">{task.description || t('common.noDescription')}</p>
                    {next ? (
                      <button
                        type="button"
                        className="btn-secondary"
                        onClick={(event) => {
                          event.stopPropagation()
                          patchTaskStatusMutation.mutate({ id: task.id, status: next })
                        }}
                        disabled={patchTaskStatusMutation.isPending}
                      >
                        {t('pm.moveTo', { status: taskStatusLabel(next) })}
                      </button>
                    ) : (
                      <span className="done-badge">{t('pm.completed')}</span>
                    )}
                  </div>
                )
              })}
              {board[status].length === 0 && <p className="muted">{t('pm.noTasks')}</p>}
            </div>
          </article>
        ))}
      </section>

      <section className="panel task-detail-panel stack-md">
        <div className="task-detail-heading">
          <h2>{t('pm.taskDetail')}</h2>
          {selectedTaskId && <small className="muted">#{selectedTaskId}</small>}
        </div>

        {!selectedTaskId && <p className="muted">{t('pm.selectTaskHint')}</p>}

        {selectedTaskId && taskDetailQuery.isPending && <p className="muted">{t('pm.loadingTaskDetail')}</p>}

        {selectedTaskId && taskDetailQuery.data && (
          <article className="stack-sm">
            <strong>{taskDetailQuery.data.title}</strong>
            <p className="muted">{taskDetailQuery.data.description || t('common.noDescription')}</p>
            <div className="task-meta-row">
              <span>
                {t('common.status')}: {taskStatusLabel(taskDetailQuery.data.status)}
              </span>
              <span>
                {t('common.project')}: #{taskDetailQuery.data.projectId}
              </span>
            </div>
          </article>
        )}

        {selectedTaskId && (
          <article className="stack-sm">
            <h3>{t('pm.transitionLogs')}</h3>
            {taskLogsQuery.isPending && <p className="muted">{t('pm.loadingTransitionLogs')}</p>}
            {!taskLogsQuery.isPending &&
              (taskLogsQuery.data?.items ?? []).map((log) => (
                <div key={log.id} className="log-item">
                  <p>
                    <strong>{taskStatusLabel(log.fromStatus)}</strong> → <strong>{taskStatusLabel(log.toStatus)}</strong>
                  </p>
                  <p className="muted">
                    {t('pm.byUserAt', {
                      userId: log.operatorId,
                      timestamp: formatLogTimestamp(log.createdAt, locale),
                    })}
                  </p>
                </div>
              ))}
            {!taskLogsQuery.isPending && (taskLogsQuery.data?.items?.length ?? 0) === 0 && (
              <p className="muted">{t('pm.noTransitionsYet')}</p>
            )}
          </article>
        )}
      </section>
    </section>
  )
}
