export type LeadStatus = 'new' | 'contacted' | 'qualified' | 'won' | 'lost'

export const orderedLeadStatuses: LeadStatus[] = ['new', 'contacted', 'qualified', 'won', 'lost']

const advanceChain: LeadStatus[] = ['new', 'contacted', 'qualified', 'won']

export function nextLeadStatus(status: LeadStatus): LeadStatus | null {
  const idx = advanceChain.indexOf(status)
  if (idx < 0 || idx === advanceChain.length - 1) {
    return null
  }

  return advanceChain[idx + 1]
}
