import { describe, expect, it } from 'vitest'

import { nextLeadStatus } from './crm-utils'

describe('nextLeadStatus', () => {
  it('returns the next status in the quick advance chain', () => {
    expect(nextLeadStatus('new')).toBe('contacted')
    expect(nextLeadStatus('contacted')).toBe('qualified')
    expect(nextLeadStatus('qualified')).toBe('won')
  })

  it('returns null for terminal statuses', () => {
    expect(nextLeadStatus('won')).toBeNull()
    expect(nextLeadStatus('lost')).toBeNull()
  })
})
