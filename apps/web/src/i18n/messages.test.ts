import { describe, expect, it } from 'vitest'

import { messages } from './messages'

describe('i18n message packs', () => {
  it('keeps zh-CN and en-US key sets in sync', () => {
    const zhKeys = Object.keys(messages['zh-CN']).sort()
    const enKeys = Object.keys(messages['en-US']).sort()

    expect(zhKeys).toEqual(enKeys)
  })
})
