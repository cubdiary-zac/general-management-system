import { describe, expect, it } from 'vitest'

import { nextTaskStatus } from './pm-utils'

describe('nextTaskStatus', () => {
  it('returns next status in ordered chain', () => {
    expect(nextTaskStatus('todo')).toBe('in_progress')
    expect(nextTaskStatus('in_progress')).toBe('in_review')
    expect(nextTaskStatus('in_review')).toBe('done')
  })

  it('returns null at final status', () => {
    expect(nextTaskStatus('done')).toBeNull()
  })
})
