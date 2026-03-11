import { FormEvent, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { apiClient } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'
import { LeadStatus, nextLeadStatus, orderedLeadStatuses } from './crm-utils'

type Customer = {
  id: number
  name: string
  company: string
  email: string
  phone: string
  ownerId: number
}

type Lead = {
  id: number
  customerId?: number
  name: string
  source: string
  status: LeadStatus
  amount?: number
  ownerId: number
}

type ListResponse<T> = {
  items: T[]
}

type SummaryResponse = {
  counts: Partial<Record<LeadStatus, number>>
}

type LeadStatusFilter = LeadStatus | 'all'

export function CRMPage() {
  const { token } = useAuth()
  const { locale, t } = useI18n()
  const queryClient = useQueryClient()

  const [customerName, setCustomerName] = useState('')
  const [customerCompany, setCustomerCompany] = useState('')
  const [customerEmail, setCustomerEmail] = useState('')
  const [customerPhone, setCustomerPhone] = useState('')

  const [leadName, setLeadName] = useState('')
  const [leadSource, setLeadSource] = useState('')
  const [leadCustomerId, setLeadCustomerId] = useState('')
  const [leadAmount, setLeadAmount] = useState('')

  const [statusFilter, setStatusFilter] = useState<LeadStatusFilter>('all')
  const [keyword, setKeyword] = useState('')

  const trimmedKeyword = keyword.trim()
  const leadListQueryKey = ['crm-leads', statusFilter, trimmedKeyword] as const

  const currencyFormatter = useMemo(
    () =>
      new Intl.NumberFormat(locale, {
        style: 'currency',
        currency: 'USD',
      }),
    [locale],
  )

  function leadStatusLabel(status: LeadStatus): string {
    return t(`status.lead.${status}`)
  }

  const customersQuery = useQuery({
    queryKey: ['crm-customers'],
    queryFn: () => apiClient.get<ListResponse<Customer>>('/api/crm/customers', token ?? undefined),
    enabled: Boolean(token),
  })

  const leadsQuery = useQuery({
    queryKey: leadListQueryKey,
    queryFn: () => {
      const params = new URLSearchParams()
      if (statusFilter !== 'all') {
        params.set('status', statusFilter)
      }
      if (trimmedKeyword !== '') {
        params.set('q', trimmedKeyword)
      }

      const query = params.toString()
      const path = query ? `/api/crm/leads?${query}` : '/api/crm/leads'
      return apiClient.get<ListResponse<Lead>>(path, token ?? undefined)
    },
    enabled: Boolean(token),
  })

  const summaryQuery = useQuery({
    queryKey: ['crm-summary'],
    queryFn: () => apiClient.get<SummaryResponse>('/api/crm/summary', token ?? undefined),
    enabled: Boolean(token),
  })

  const createCustomerMutation = useMutation({
    mutationFn: (payload: { name: string; company: string; email: string; phone: string }) =>
      apiClient.post<Customer>('/api/crm/customers', payload, token ?? undefined),
    onSuccess: async () => {
      setCustomerName('')
      setCustomerCompany('')
      setCustomerEmail('')
      setCustomerPhone('')
      await queryClient.invalidateQueries({ queryKey: ['crm-customers'] })
    },
  })

  const createLeadMutation = useMutation({
    mutationFn: (payload: { customerId?: number; name: string; source: string; amount?: number }) =>
      apiClient.post<Lead>('/api/crm/leads', payload, token ?? undefined),
    onSuccess: async () => {
      setLeadName('')
      setLeadSource('')
      setLeadCustomerId('')
      setLeadAmount('')
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: leadListQueryKey }),
        queryClient.invalidateQueries({ queryKey: ['crm-summary'] }),
      ])
    },
  })

  const patchLeadStatusMutation = useMutation({
    mutationFn: (payload: { id: number; status: LeadStatus }) =>
      apiClient.patch<Lead>(`/api/crm/leads/${payload.id}/status`, { status: payload.status }, token ?? undefined),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: leadListQueryKey }),
        queryClient.invalidateQueries({ queryKey: ['crm-summary'] }),
      ])
    },
  })

  const customerByID = useMemo(() => {
    const map = new Map<number, Customer>()
    for (const customer of customersQuery.data?.items ?? []) {
      map.set(customer.id, customer)
    }

    return map
  }, [customersQuery.data])

  const summaryCounts = useMemo(() => {
    const counts: Record<LeadStatus, number> = {
      new: 0,
      contacted: 0,
      qualified: 0,
      won: 0,
      lost: 0,
    }

    const payload = summaryQuery.data?.counts ?? {}
    for (const status of orderedLeadStatuses) {
      const value = payload[status]
      if (typeof value === 'number') {
        counts[status] = value
      }
    }

    return counts
  }, [summaryQuery.data])

  const createCustomerError = createCustomerMutation.error instanceof Error ? createCustomerMutation.error.message : null
  const createLeadError = createLeadMutation.error instanceof Error ? createLeadMutation.error.message : null
  const patchLeadError = patchLeadStatusMutation.error instanceof Error ? patchLeadStatusMutation.error.message : null

  function onCreateCustomer(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const name = customerName.trim()
    if (!name) {
      return
    }

    createCustomerMutation.mutate({
      name,
      company: customerCompany.trim(),
      email: customerEmail.trim(),
      phone: customerPhone.trim(),
    })
  }

  function onCreateLead(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const name = leadName.trim()
    const source = leadSource.trim()
    if (!name || !source) {
      return
    }

    let customerId: number | undefined
    if (leadCustomerId !== '') {
      const parsedCustomerID = Number.parseInt(leadCustomerId, 10)
      if (Number.isNaN(parsedCustomerID) || parsedCustomerID <= 0) {
        return
      }
      customerId = parsedCustomerID
    }

    let amount: number | undefined
    const trimmedAmount = leadAmount.trim()
    if (trimmedAmount !== '') {
      const parsedAmount = Number.parseFloat(trimmedAmount)
      if (Number.isNaN(parsedAmount) || parsedAmount < 0) {
        return
      }
      amount = parsedAmount
    }

    createLeadMutation.mutate({
      customerId,
      name,
      source,
      amount,
    })
  }

  return (
    <section className="stack-lg">
      <header className="page-header">
        <div>
          <h1>{t('crm.title')}</h1>
          <p className="muted">{t('crm.subtitle')}</p>
        </div>
      </header>

      <section className="summary-grid">
        {orderedLeadStatuses.map((status) => (
          <article key={status} className="summary-card">
            <p className="muted">{leadStatusLabel(status)}</p>
            <strong>{summaryCounts[status]}</strong>
          </article>
        ))}
      </section>

      <div className="panel-grid">
        <article className="panel stack-md">
          <h2>{t('crm.customers')}</h2>
          <form onSubmit={onCreateCustomer} className="stack-sm">
            <input
              placeholder={t('crm.customerName')}
              value={customerName}
              onChange={(event) => setCustomerName(event.target.value)}
            />
            <input
              placeholder={t('common.company')}
              value={customerCompany}
              onChange={(event) => setCustomerCompany(event.target.value)}
            />
            <input
              placeholder={t('common.email')}
              value={customerEmail}
              onChange={(event) => setCustomerEmail(event.target.value)}
            />
            <input
              placeholder={t('common.phone')}
              value={customerPhone}
              onChange={(event) => setCustomerPhone(event.target.value)}
            />
            <button type="submit" className="btn-primary" disabled={createCustomerMutation.isPending}>
              {createCustomerMutation.isPending ? t('crm.creatingCustomer') : t('crm.createCustomer')}
            </button>
            {createCustomerError && <p className="error-text">{createCustomerError}</p>}
          </form>

          {customersQuery.isPending && <p className="muted">{t('crm.loadingCustomers')}</p>}

          <div className="stack-sm simple-list">
            {(customersQuery.data?.items ?? []).map((customer) => (
              <article key={customer.id} className="list-item stack-sm">
                <div className="row-between">
                  <strong>{customer.name}</strong>
                  <small className="muted">#{customer.id}</small>
                </div>
                <p className="muted">{customer.company || t('crm.noCompany')}</p>
                <p className="muted">
                  {customer.email || t('crm.noEmail')} | {customer.phone || t('crm.noPhone')}
                </p>
              </article>
            ))}
            {(customersQuery.data?.items?.length ?? 0) === 0 && !customersQuery.isPending && (
              <p className="muted">{t('crm.noCustomersYet')}</p>
            )}
          </div>
        </article>

        <article className="panel stack-md">
          <h2>{t('crm.leads')}</h2>
          <form onSubmit={onCreateLead} className="stack-sm">
            <input placeholder={t('crm.leadName')} value={leadName} onChange={(event) => setLeadName(event.target.value)} />
            <input
              placeholder={t('crm.leadSource')}
              value={leadSource}
              onChange={(event) => setLeadSource(event.target.value)}
            />
            <select value={leadCustomerId} onChange={(event) => setLeadCustomerId(event.target.value)}>
              <option value="">{t('crm.unlinkedCustomer')}</option>
              {(customersQuery.data?.items ?? []).map((customer) => (
                <option key={customer.id} value={String(customer.id)}>
                  {customer.name}
                </option>
              ))}
            </select>
            <input
              type="number"
              min={0}
              step="0.01"
              placeholder={t('crm.amountOptional')}
              value={leadAmount}
              onChange={(event) => setLeadAmount(event.target.value)}
            />
            <button type="submit" className="btn-primary" disabled={createLeadMutation.isPending}>
              {createLeadMutation.isPending ? t('crm.creatingLead') : t('crm.createLead')}
            </button>
            {createLeadError && <p className="error-text">{createLeadError}</p>}
          </form>

          <div className="task-filters">
            <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as LeadStatusFilter)}>
              <option value="all">{t('common.allStatuses')}</option>
              {orderedLeadStatuses.map((status) => (
                <option key={status} value={status}>
                  {leadStatusLabel(status)}
                </option>
              ))}
            </select>
            <input
              placeholder={t('crm.searchLeadByName')}
              value={keyword}
              onChange={(event) => setKeyword(event.target.value)}
            />
          </div>

          {leadsQuery.isPending && <p className="muted">{t('crm.loadingLeads')}</p>}
          {patchLeadError && <p className="error-text">{patchLeadError}</p>}

          <div className="stack-sm simple-list">
            {(leadsQuery.data?.items ?? []).map((lead) => {
              const nextStatus = nextLeadStatus(lead.status)
              const customer = lead.customerId ? customerByID.get(lead.customerId) : null

              return (
                <article key={lead.id} className="list-item stack-sm">
                  <div className="row-between">
                    <strong>{lead.name}</strong>
                    <span className={`status-badge status-${lead.status}`}>{leadStatusLabel(lead.status)}</span>
                  </div>
                  <p className="muted">
                    {t('common.source')}: {lead.source}
                  </p>
                  <p className="muted">
                    {t('common.customer')}: {customer?.name ?? t('crm.unlinkedCustomer')}
                  </p>
                  {typeof lead.amount === 'number' && (
                    <p className="muted">
                      {t('common.amount')}: {currencyFormatter.format(lead.amount)}
                    </p>
                  )}
                  <div className="lead-actions">
                    {nextStatus ? (
                      <button
                        type="button"
                        className="btn-secondary"
                        onClick={() => patchLeadStatusMutation.mutate({ id: lead.id, status: nextStatus })}
                        disabled={patchLeadStatusMutation.isPending}
                      >
                        {t('crm.advanceTo', { status: leadStatusLabel(nextStatus) })}
                      </button>
                    ) : (
                      <small className="muted">{t('common.terminalStatus')}</small>
                    )}
                  </div>
                </article>
              )
            })}
            {(leadsQuery.data?.items?.length ?? 0) === 0 && !leadsQuery.isPending && (
              <p className="muted">{t('crm.noLeadsFound')}</p>
            )}
          </div>
        </article>
      </div>
    </section>
  )
}
