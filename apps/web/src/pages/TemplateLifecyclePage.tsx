import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { apiClient } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'

type TemplateStatus = 'draft' | 'published'
type TemplateStatusFilter = TemplateStatus | 'all'

type ProjectTemplate = {
  id: number
  industryTemplateId: number
  name: string
  code: string
  description: string
  version: number
  status: TemplateStatus
}

type ListResponse<T> = {
  items: T[]
}

type ProjectTemplateLifecycle = {
  projectTemplate: ProjectTemplate
  industryTemplate: {
    id: number
    name: string
    code: string
    status: TemplateStatus
    version: number
  }
  counts: {
    publishedStages: number
    publishedForms: number
    publishedFields: number
    runtimeProjects: number
  }
  runtimeByStageStatus: {
    pending: number
    active: number
    done: number
  }
  guidance: string
}

type NextProjectTemplateVersionResponse = {
  projectTemplate: ProjectTemplate
}

type InstantiateProjectTemplateResponse = {
  project: {
    id: number
  }
}

function generateRuntimeProjectName(templateName: string): string {
  const now = new Date()
  const stamp = now.toISOString().replace(/[:.]/g, '-')
  return `${templateName} Runtime ${stamp}`
}

export function TemplateLifecyclePage() {
  const { token, user } = useAuth()
  const { t } = useI18n()
  const queryClient = useQueryClient()

  const [statusFilter, setStatusFilter] = useState<TemplateStatusFilter>('all')
  const [industryTemplateIDFilter, setIndustryTemplateIDFilter] = useState('')
  const [versionFilter, setVersionFilter] = useState('')
  const [selectedLifecycleTemplateID, setSelectedLifecycleTemplateID] = useState<number | null>(null)

  const normalizedIndustryTemplateID = industryTemplateIDFilter.trim()
  const normalizedVersion = versionFilter.trim()
  const templateListQueryKey = ['tmpl-project-templates', statusFilter, normalizedIndustryTemplateID, normalizedVersion] as const

  const templatesQuery = useQuery({
    queryKey: templateListQueryKey,
    queryFn: () => {
      const params = new URLSearchParams()
      if (statusFilter !== 'all') {
        params.set('status', statusFilter)
      }
      if (normalizedIndustryTemplateID !== '') {
        params.set('industryTemplateId', normalizedIndustryTemplateID)
      }
      if (normalizedVersion !== '') {
        params.set('version', normalizedVersion)
      }

      const query = params.toString()
      const path = query ? `/api/tmpl/project-templates?${query}` : '/api/tmpl/project-templates'
      return apiClient.get<ListResponse<ProjectTemplate>>(path, token ?? undefined)
    },
    enabled: Boolean(token),
  })

  const lifecycleQuery = useQuery({
    queryKey: ['tmpl-lifecycle', selectedLifecycleTemplateID],
    queryFn: () =>
      apiClient.get<ProjectTemplateLifecycle>(
        `/api/tmpl/project-templates/${selectedLifecycleTemplateID}/lifecycle`,
        token ?? undefined,
      ),
    enabled: Boolean(token) && selectedLifecycleTemplateID !== null,
  })

  const publishMutation = useMutation({
    mutationFn: (templateID: number) =>
      apiClient.post<ProjectTemplate>(`/api/tmpl/project-templates/${templateID}/publish`, {}, token ?? undefined),
    onSuccess: async (_, templateID) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: templateListQueryKey }),
        queryClient.invalidateQueries({ queryKey: ['tmpl-lifecycle', templateID] }),
      ])
    },
  })

  const unpublishMutation = useMutation({
    mutationFn: (templateID: number) =>
      apiClient.post<ProjectTemplate>(`/api/tmpl/project-templates/${templateID}/unpublish`, {}, token ?? undefined),
    onSuccess: async (_, templateID) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: templateListQueryKey }),
        queryClient.invalidateQueries({ queryKey: ['tmpl-lifecycle', templateID] }),
      ])
    },
  })

  const nextVersionMutation = useMutation({
    mutationFn: (templateID: number) =>
      apiClient.post<NextProjectTemplateVersionResponse>(
        `/api/tmpl/project-templates/${templateID}/next-version`,
        {},
        token ?? undefined,
      ),
    onSuccess: async (_, templateID) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: templateListQueryKey }),
        queryClient.invalidateQueries({ queryKey: ['tmpl-lifecycle', templateID] }),
      ])
    },
  })

  const instantiateMutation = useMutation({
    mutationFn: (payload: { templateID: number; templateName: string }) =>
      apiClient.post<InstantiateProjectTemplateResponse>(
        `/api/tmpl/project-templates/${payload.templateID}/instantiate`,
        {
          name: generateRuntimeProjectName(payload.templateName),
          description: t('tmpl.instantiateGeneratedDescription', { templateName: payload.templateName }),
        },
        token ?? undefined,
      ),
    onSuccess: async (_, payload) => {
      await queryClient.invalidateQueries({ queryKey: ['tmpl-lifecycle', payload.templateID] })
    },
  })

  const templates = templatesQuery.data?.items ?? []
  const selectedLifecycleTemplate = useMemo(
    () => templates.find((item) => item.id === selectedLifecycleTemplateID) ?? null,
    [selectedLifecycleTemplateID, templates],
  )

  const canMutateTemplates = user?.role !== 'viewer'
  const actionPending =
    publishMutation.isPending ||
    unpublishMutation.isPending ||
    nextVersionMutation.isPending ||
    instantiateMutation.isPending

  const actionError = useMemo(() => {
    const candidates = [publishMutation.error, unpublishMutation.error, nextVersionMutation.error, instantiateMutation.error]
    for (const item of candidates) {
      if (item instanceof Error) {
        return item.message
      }
    }
    return null
  }, [instantiateMutation.error, nextVersionMutation.error, publishMutation.error, unpublishMutation.error])

  const templatesError = templatesQuery.error instanceof Error ? templatesQuery.error.message : null
  const lifecycleError = lifecycleQuery.error instanceof Error ? lifecycleQuery.error.message : null

  function statusLabel(status: TemplateStatus): string {
    return t(`tmpl.status.${status}`)
  }

  return (
    <section className="stack-lg">
      <header className="page-header">
        <div>
          <h1>{t('tmpl.title')}</h1>
          <p className="muted">{t('tmpl.subtitle')}</p>
        </div>
      </header>

      <section className="panel stack-md">
        <div className="tmpl-filter-grid">
          <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as TemplateStatusFilter)}>
            <option value="all">{t('common.allStatuses')}</option>
            <option value="draft">{statusLabel('draft')}</option>
            <option value="published">{statusLabel('published')}</option>
          </select>
          <input
            type="number"
            min={1}
            placeholder={t('tmpl.filterIndustryTemplateId')}
            value={industryTemplateIDFilter}
            onChange={(event) => setIndustryTemplateIDFilter(event.target.value)}
          />
          <input
            type="number"
            min={1}
            placeholder={t('tmpl.filterVersion')}
            value={versionFilter}
            onChange={(event) => setVersionFilter(event.target.value)}
          />
        </div>

        {templatesQuery.isPending && <p className="muted">{t('tmpl.loadingTemplates')}</p>}
        {templatesError && <p className="error-text">{templatesError}</p>}
        {actionError && <p className="error-text">{actionError}</p>}
        {!canMutateTemplates && <p className="muted">{t('tmpl.viewerReadOnlyHint')}</p>}

        <div className="stack-sm simple-list">
          {templates.map((item) => (
            <article key={item.id} className="list-item stack-sm">
              <div className="row-between">
                <strong>{item.name}</strong>
                <span className={`status-badge status-${item.status}`}>{statusLabel(item.status)}</span>
              </div>
              <div className="tmpl-row-meta">
                <span>
                  {t('tmpl.code')}: {item.code}
                </span>
                <span>
                  {t('tmpl.version')}: {item.version}
                </span>
                <span>
                  {t('tmpl.industryTemplateId')}: #{item.industryTemplateId}
                </span>
              </div>
              <div className="tmpl-actions">
                {item.status === 'draft' ? (
                  <button
                    type="button"
                    className="btn-secondary"
                    disabled={!canMutateTemplates || actionPending}
                    onClick={() => publishMutation.mutate(item.id)}
                  >
                    {publishMutation.isPending ? t('tmpl.publishing') : t('tmpl.publish')}
                  </button>
                ) : (
                  <button
                    type="button"
                    className="btn-secondary"
                    disabled={!canMutateTemplates || actionPending}
                    onClick={() => unpublishMutation.mutate(item.id)}
                  >
                    {unpublishMutation.isPending ? t('tmpl.unpublishing') : t('tmpl.unpublish')}
                  </button>
                )}
                <button
                  type="button"
                  className="btn-secondary"
                  disabled={!canMutateTemplates || actionPending}
                  onClick={() => nextVersionMutation.mutate(item.id)}
                >
                  {nextVersionMutation.isPending ? t('tmpl.creatingNextVersion') : t('tmpl.nextVersion')}
                </button>
                <button
                  type="button"
                  className={selectedLifecycleTemplateID === item.id ? 'btn-primary' : 'btn-secondary'}
                  onClick={() => setSelectedLifecycleTemplateID(item.id)}
                >
                  {t('tmpl.viewLifecycle')}
                </button>
                <button
                  type="button"
                  className="btn-primary"
                  disabled={!canMutateTemplates || actionPending}
                  onClick={() => instantiateMutation.mutate({ templateID: item.id, templateName: item.name })}
                >
                  {instantiateMutation.isPending ? t('tmpl.instantiating') : t('tmpl.instantiate')}
                </button>
              </div>
            </article>
          ))}
          {templates.length === 0 && !templatesQuery.isPending && <p className="muted">{t('tmpl.noTemplates')}</p>}
        </div>
      </section>

      <section className="panel stack-md">
        <div className="row-between">
          <h2>{t('tmpl.lifecycleTitle')}</h2>
          {selectedLifecycleTemplate && <small className="muted">#{selectedLifecycleTemplate.id}</small>}
        </div>

        {!selectedLifecycleTemplateID && <p className="muted">{t('tmpl.lifecycleSelectHint')}</p>}
        {selectedLifecycleTemplateID && lifecycleQuery.isPending && <p className="muted">{t('tmpl.loadingLifecycle')}</p>}
        {lifecycleError && <p className="error-text">{lifecycleError}</p>}

        {selectedLifecycleTemplateID && lifecycleQuery.data && (
          <div className="stack-md">
            <div className="tmpl-row-meta">
              <span>
                {t('tmpl.projectTemplate')}: {lifecycleQuery.data.projectTemplate.name}
              </span>
              <span>
                {t('tmpl.status')}: {statusLabel(lifecycleQuery.data.projectTemplate.status)}
              </span>
            </div>
            <div className="tmpl-row-meta">
              <span>
                {t('tmpl.industry')}: {lifecycleQuery.data.industryTemplate.name}
              </span>
              <span>
                {t('tmpl.code')}: {lifecycleQuery.data.industryTemplate.code}
              </span>
              <span>
                {t('tmpl.version')}: {lifecycleQuery.data.industryTemplate.version}
              </span>
            </div>
            <div className="tmpl-count-grid">
              <article className="summary-card">
                <p className="muted">{t('tmpl.publishedStages')}</p>
                <strong>{lifecycleQuery.data.counts.publishedStages}</strong>
              </article>
              <article className="summary-card">
                <p className="muted">{t('tmpl.publishedForms')}</p>
                <strong>{lifecycleQuery.data.counts.publishedForms}</strong>
              </article>
              <article className="summary-card">
                <p className="muted">{t('tmpl.publishedFields')}</p>
                <strong>{lifecycleQuery.data.counts.publishedFields}</strong>
              </article>
              <article className="summary-card">
                <p className="muted">{t('tmpl.runtimeProjects')}</p>
                <strong>{lifecycleQuery.data.counts.runtimeProjects}</strong>
              </article>
              <article className="summary-card">
                <p className="muted">{t('tmpl.stagePending')}</p>
                <strong>{lifecycleQuery.data.runtimeByStageStatus.pending}</strong>
              </article>
              <article className="summary-card">
                <p className="muted">{t('tmpl.stageActive')}</p>
                <strong>{lifecycleQuery.data.runtimeByStageStatus.active}</strong>
              </article>
              <article className="summary-card">
                <p className="muted">{t('tmpl.stageDone')}</p>
                <strong>{lifecycleQuery.data.runtimeByStageStatus.done}</strong>
              </article>
            </div>
            <article className="log-item">
              <p>
                <strong>{t('tmpl.guidance')}</strong>
              </p>
              <p className="muted">{lifecycleQuery.data.guidance}</p>
            </article>
          </div>
        )}
      </section>
    </section>
  )
}
