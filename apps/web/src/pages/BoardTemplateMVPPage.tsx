import { FormEvent, useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { apiClient } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'

type HeaderWidgetType = 'input' | 'textarea' | 'attachment' | 'select' | 'date'

type ListResponse<T> = {
  items: T[]
}

type ProjectTemplate = {
  id: number
  name: string
  code: string
  description: string
}

type ProjectTemplateHeaderField = {
  id: number
  projectTemplateId: number
  name: string
  code: string
  widgetType: HeaderWidgetType
  required: boolean
  position: number
}

type StageTemplate = {
  id: number
  projectTemplateId: number
  name: string
  code: string
  description: string
  position: number
}

const headerWidgetTypes: HeaderWidgetType[] = ['input', 'textarea', 'attachment', 'select', 'date']

function maxPosition<T extends { position: number }>(items: T[]): number {
  let max = 0
  for (const item of items) {
    if (item.position > max) {
      max = item.position
    }
  }
  return max
}

function readErrorMessage(error: unknown): string | null {
  return error instanceof Error ? error.message : null
}

export function BoardTemplateMVPPage() {
  const { token, user } = useAuth()
  const { t } = useI18n()
  const queryClient = useQueryClient()

  const [selectedTemplateID, setSelectedTemplateID] = useState<number | null>(null)

  const [headerFieldName, setHeaderFieldName] = useState('')
  const [headerFieldCode, setHeaderFieldCode] = useState('')
  const [headerFieldWidgetType, setHeaderFieldWidgetType] = useState<HeaderWidgetType>('input')
  const [headerFieldRequired, setHeaderFieldRequired] = useState(false)

  const [stageName, setStageName] = useState('')
  const [stageCode, setStageCode] = useState('')
  const [stageNameDrafts, setStageNameDrafts] = useState<Record<number, string>>({})

  const [actionError, setActionError] = useState<string | null>(null)

  const templatesQueryKey = ['board-template-mvp', 'project-templates'] as const
  const headerFieldsQueryKey = ['board-template-mvp', 'header-fields', selectedTemplateID] as const
  const stagesQueryKey = ['board-template-mvp', 'stages', selectedTemplateID] as const

  const templatesQuery = useQuery({
    queryKey: templatesQueryKey,
    queryFn: () => apiClient.get<ListResponse<ProjectTemplate>>('/api/tmpl/project-templates', token ?? undefined),
    enabled: Boolean(token),
  })

  const headerFieldsQuery = useQuery({
    queryKey: headerFieldsQueryKey,
    queryFn: () =>
      apiClient.get<ListResponse<ProjectTemplateHeaderField>>(
        `/api/tmpl/project-templates/${selectedTemplateID}/header-fields`,
        token ?? undefined,
      ),
    enabled: Boolean(token) && selectedTemplateID !== null,
  })

  const stagesQuery = useQuery({
    queryKey: stagesQueryKey,
    queryFn: () =>
      apiClient.get<ListResponse<StageTemplate>>(
        `/api/tmpl/stage-templates?projectTemplateId=${selectedTemplateID}`,
        token ?? undefined,
      ),
    enabled: Boolean(token) && selectedTemplateID !== null,
  })

  const createHeaderFieldMutation = useMutation({
    mutationFn: (payload: {
      name: string
      code: string
      widgetType: HeaderWidgetType
      required: boolean
      position: number
    }) => {
      if (selectedTemplateID === null) {
        throw new Error('project template is required')
      }

      return apiClient.post<ProjectTemplateHeaderField>(
        `/api/tmpl/project-templates/${selectedTemplateID}/header-fields`,
        payload,
        token ?? undefined,
      )
    },
  })

  const patchHeaderFieldMutation = useMutation({
    mutationFn: (payload: { fieldID: number; data: Partial<ProjectTemplateHeaderField> }) => {
      if (selectedTemplateID === null) {
        throw new Error('project template is required')
      }

      return apiClient.patch<ProjectTemplateHeaderField>(
        `/api/tmpl/project-templates/${selectedTemplateID}/header-fields/${payload.fieldID}`,
        payload.data,
        token ?? undefined,
      )
    },
  })

  const createStageMutation = useMutation({
    mutationFn: (payload: { name: string; code: string; position: number }) => {
      if (selectedTemplateID === null) {
        throw new Error('project template is required')
      }

      return apiClient.post<StageTemplate>(
        '/api/tmpl/stage-templates',
        {
          projectTemplateId: selectedTemplateID,
          ...payload,
        },
        token ?? undefined,
      )
    },
  })

  const patchStageMutation = useMutation({
    mutationFn: (payload: { stageID: number; data: Partial<StageTemplate> }) =>
      apiClient.patch<StageTemplate>(`/api/tmpl/stage-templates/${payload.stageID}`, payload.data, token ?? undefined),
  })

  const templates = templatesQuery.data?.items ?? []
  const headerFields = headerFieldsQuery.data?.items ?? []
  const stages = stagesQuery.data?.items ?? []

  useEffect(() => {
    if (templates.length === 0) {
      setSelectedTemplateID(null)
      return
    }

    const selectedStillExists = templates.some((item) => item.id === selectedTemplateID)
    if (!selectedStillExists) {
      setSelectedTemplateID(templates[0].id)
    }
  }, [selectedTemplateID, templates])

  useEffect(() => {
    setStageNameDrafts((prev) => {
      const next: Record<number, string> = {}
      for (const stage of stages) {
        next[stage.id] = prev[stage.id] ?? stage.name
      }
      return next
    })
  }, [stages])

  const canMutate = user?.role !== 'viewer'
  const mutationPending =
    createHeaderFieldMutation.isPending ||
    patchHeaderFieldMutation.isPending ||
    createStageMutation.isPending ||
    patchStageMutation.isPending

  const queryError = useMemo(() => {
    const errors = [templatesQuery.error, headerFieldsQuery.error, stagesQuery.error]
    for (const error of errors) {
      const message = readErrorMessage(error)
      if (message !== null) {
        return message
      }
    }
    return null
  }, [templatesQuery.error, headerFieldsQuery.error, stagesQuery.error])

  const mutationError = useMemo(() => {
    const errors = [
      createHeaderFieldMutation.error,
      patchHeaderFieldMutation.error,
      createStageMutation.error,
      patchStageMutation.error,
    ]
    for (const error of errors) {
      const message = readErrorMessage(error)
      if (message !== null) {
        return message
      }
    }
    return null
  }, [
    createHeaderFieldMutation.error,
    patchHeaderFieldMutation.error,
    createStageMutation.error,
    patchStageMutation.error,
  ])

  async function refreshHeaderFields() {
    await queryClient.invalidateQueries({ queryKey: headerFieldsQueryKey })
  }

  async function refreshStages() {
    await queryClient.invalidateQueries({ queryKey: stagesQueryKey })
  }

  async function onCreateHeaderField(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (selectedTemplateID === null) {
      return
    }

    const name = headerFieldName.trim()
    const code = headerFieldCode.trim()
    if (name === '' || code === '') {
      return
    }

    setActionError(null)
    try {
      await createHeaderFieldMutation.mutateAsync({
        name,
        code,
        widgetType: headerFieldWidgetType,
        required: headerFieldRequired,
        position: maxPosition(headerFields) + 1,
      })
      setHeaderFieldName('')
      setHeaderFieldCode('')
      setHeaderFieldWidgetType('input')
      setHeaderFieldRequired(false)
      await refreshHeaderFields()
    } catch (error) {
      setActionError(readErrorMessage(error))
    }
  }

  async function onMoveHeaderField(index: number, direction: -1 | 1) {
    if (!canMutate || selectedTemplateID === null) {
      return
    }

    const targetIndex = index + direction
    if (targetIndex < 0 || targetIndex >= headerFields.length) {
      return
    }

    const currentField = headerFields[index]
    const targetField = headerFields[targetIndex]

    setActionError(null)
    try {
      await patchHeaderFieldMutation.mutateAsync({
        fieldID: currentField.id,
        data: { position: targetField.position },
      })
      await patchHeaderFieldMutation.mutateAsync({
        fieldID: targetField.id,
        data: { position: currentField.position },
      })
      await refreshHeaderFields()
    } catch (error) {
      setActionError(readErrorMessage(error))
    }
  }

  async function onCreateStage(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (selectedTemplateID === null) {
      return
    }

    const name = stageName.trim()
    const code = stageCode.trim()
    if (name === '' || code === '') {
      return
    }

    setActionError(null)
    try {
      await createStageMutation.mutateAsync({
        name,
        code,
        position: maxPosition(stages) + 1,
      })
      setStageName('')
      setStageCode('')
      await refreshStages()
    } catch (error) {
      setActionError(readErrorMessage(error))
    }
  }

  async function onSaveStageName(stage: StageTemplate) {
    const nameDraft = (stageNameDrafts[stage.id] ?? stage.name).trim()
    if (nameDraft === '' || nameDraft === stage.name) {
      return
    }

    setActionError(null)
    try {
      await patchStageMutation.mutateAsync({
        stageID: stage.id,
        data: { name: nameDraft },
      })
      await refreshStages()
    } catch (error) {
      setActionError(readErrorMessage(error))
    }
  }

  async function onMoveStage(index: number, direction: -1 | 1) {
    if (!canMutate) {
      return
    }

    const targetIndex = index + direction
    if (targetIndex < 0 || targetIndex >= stages.length) {
      return
    }

    const currentStage = stages[index]
    const targetStage = stages[targetIndex]

    setActionError(null)
    try {
      await patchStageMutation.mutateAsync({
        stageID: currentStage.id,
        data: { position: targetStage.position },
      })
      await patchStageMutation.mutateAsync({
        stageID: targetStage.id,
        data: { position: currentStage.position },
      })
      await refreshStages()
    } catch (error) {
      setActionError(readErrorMessage(error))
    }
  }

  return (
    <section className="stack-lg">
      <header className="page-header">
        <div>
          <h1>{t('boardMvp.title')}</h1>
          <p className="muted">{t('boardMvp.subtitle')}</p>
        </div>
      </header>

      <section className="panel stack-md">
        <h2>{t('boardMvp.templateSelect')}</h2>
        {templatesQuery.isPending && <p className="muted">{t('boardMvp.loadingTemplates')}</p>}
        {templates.length === 0 && !templatesQuery.isPending && <p className="muted">{t('boardMvp.noTemplates')}</p>}
        <div className="stack-sm">
          {templates.map((template) => (
            <button
              key={template.id}
              type="button"
              className={selectedTemplateID === template.id ? 'project-item active' : 'project-item'}
              onClick={() => setSelectedTemplateID(template.id)}
            >
              <span>{template.name}</span>
              <small className="muted">#{template.id}</small>
            </button>
          ))}
        </div>
      </section>

      {!canMutate && <p className="muted">{t('boardMvp.viewerReadOnlyHint')}</p>}
      {(queryError ?? mutationError ?? actionError) && <p className="error-text">{queryError ?? mutationError ?? actionError}</p>}

      {selectedTemplateID !== null && (
        <div className="panel-grid board-template-mvp-grid">
          <section className="panel stack-md">
            <h2>{t('boardMvp.headerLabels')}</h2>

            <form onSubmit={onCreateHeaderField} className="mvp-form-grid">
              <input
                placeholder={t('boardMvp.headerFieldName')}
                value={headerFieldName}
                onChange={(event) => setHeaderFieldName(event.target.value)}
                disabled={!canMutate}
              />
              <input
                placeholder={t('boardMvp.headerFieldCode')}
                value={headerFieldCode}
                onChange={(event) => setHeaderFieldCode(event.target.value)}
                disabled={!canMutate}
              />
              <select
                value={headerFieldWidgetType}
                onChange={(event) => setHeaderFieldWidgetType(event.target.value as HeaderWidgetType)}
                disabled={!canMutate}
              >
                {headerWidgetTypes.map((widgetType) => (
                  <option key={widgetType} value={widgetType}>
                    {widgetType}
                  </option>
                ))}
              </select>
              <label className="mvp-inline-checkbox">
                <input
                  type="checkbox"
                  checked={headerFieldRequired}
                  onChange={(event) => setHeaderFieldRequired(event.target.checked)}
                  disabled={!canMutate}
                />
                <span>{t('boardMvp.required')}</span>
              </label>
              <button
                type="submit"
                className="btn-primary"
                disabled={!canMutate || mutationPending || headerFieldName.trim() === '' || headerFieldCode.trim() === ''}
              >
                {t('boardMvp.addLabel')}
              </button>
            </form>

            {headerFieldsQuery.isPending && <p className="muted">{t('boardMvp.loadingHeaderFields')}</p>}
            {!headerFieldsQuery.isPending && headerFields.length === 0 && <p className="muted">{t('boardMvp.noHeaderFields')}</p>}

            <div className="stack-sm">
              {headerFields.map((field, index) => (
                <article key={field.id} className="list-item stack-sm">
                  <div className="row-between">
                    <strong>{field.name}</strong>
                    <small className="muted">#{field.position}</small>
                  </div>
                  <div className="tmpl-row-meta">
                    <span>{field.code}</span>
                    <span>{field.widgetType}</span>
                    <span>{field.required ? t('boardMvp.required') : t('boardMvp.optional')}</span>
                  </div>
                  <div className="tmpl-actions">
                    <button
                      type="button"
                      className="btn-secondary"
                      disabled={!canMutate || mutationPending || index === 0}
                      onClick={() => {
                        void onMoveHeaderField(index, -1)
                      }}
                    >
                      {t('boardMvp.moveUp')}
                    </button>
                    <button
                      type="button"
                      className="btn-secondary"
                      disabled={!canMutate || mutationPending || index === headerFields.length - 1}
                      onClick={() => {
                        void onMoveHeaderField(index, 1)
                      }}
                    >
                      {t('boardMvp.moveDown')}
                    </button>
                  </div>
                </article>
              ))}
            </div>
          </section>

          <section className="panel stack-md">
            <h2>{t('boardMvp.stages')}</h2>

            <form onSubmit={onCreateStage} className="mvp-form-grid">
              <input
                placeholder={t('boardMvp.stageName')}
                value={stageName}
                onChange={(event) => setStageName(event.target.value)}
                disabled={!canMutate}
              />
              <input
                placeholder={t('boardMvp.stageCode')}
                value={stageCode}
                onChange={(event) => setStageCode(event.target.value)}
                disabled={!canMutate}
              />
              <button
                type="submit"
                className="btn-primary"
                disabled={!canMutate || mutationPending || stageName.trim() === '' || stageCode.trim() === ''}
              >
                {t('boardMvp.addStage')}
              </button>
            </form>

            {stagesQuery.isPending && <p className="muted">{t('boardMvp.loadingStages')}</p>}
            {!stagesQuery.isPending && stages.length === 0 && <p className="muted">{t('boardMvp.noStages')}</p>}

            <div className="stack-sm">
              {stages.map((stage, index) => (
                <article key={stage.id} className="list-item stack-sm">
                  <div className="row-between">
                    <small className="muted">#{stage.position}</small>
                    <div className="tmpl-actions">
                      <button
                        type="button"
                        className="btn-secondary"
                        disabled={!canMutate || mutationPending || index === 0}
                        onClick={() => {
                          void onMoveStage(index, -1)
                        }}
                      >
                        {t('boardMvp.moveUp')}
                      </button>
                      <button
                        type="button"
                        className="btn-secondary"
                        disabled={!canMutate || mutationPending || index === stages.length - 1}
                        onClick={() => {
                          void onMoveStage(index, 1)
                        }}
                      >
                        {t('boardMvp.moveDown')}
                      </button>
                    </div>
                  </div>

                  <div className="mvp-stage-edit-row">
                    <input
                      value={stageNameDrafts[stage.id] ?? ''}
                      onChange={(event) =>
                        setStageNameDrafts((prev) => ({
                          ...prev,
                          [stage.id]: event.target.value,
                        }))
                      }
                      disabled={!canMutate}
                    />
                    <button
                      type="button"
                      className="btn-primary"
                      disabled={!canMutate || mutationPending || (stageNameDrafts[stage.id] ?? '').trim() === ''}
                      onClick={() => {
                        void onSaveStageName(stage)
                      }}
                    >
                      {t('boardMvp.save')}
                    </button>
                  </div>

                  <div className="tmpl-row-meta">
                    <span>{stage.code}</span>
                    <span>{stage.description || t('common.noDescription')}</span>
                  </div>
                </article>
              ))}
            </div>
          </section>
        </div>
      )}
    </section>
  )
}
