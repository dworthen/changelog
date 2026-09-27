import { describe, expect, test } from 'bun:test'
import { parseChangeType } from './changeType'

describe('parseChangeType', () => {
  test('accepts an exact canonical change type', () => {
    const result = parseChangeType('Add')
    expect(result).toEqual({ ok: true, value: 'Add' })
  })

  test('normalizes lowercase input to the canonical change type', () => {
    const result = parseChangeType('fix')
    expect(result).toEqual({ ok: true, value: 'Fix' })
  })

  test('normalizes mixed-case input and surrounding whitespace', () => {
    const result = parseChangeType('  dEpReCaTe  ')
    expect(result).toEqual({ ok: true, value: 'Deprecate' })
  })

  test('rejects an unknown change type with the valid options listed', () => {
    const result = parseChangeType('bogus')
    expect(result.ok).toBe(false)
    if (!result.ok) {
      expect(result.error).toContain('bogus')
      expect(result.error).toContain('Add')
      expect(result.error).toContain('Internal')
    }
  })

  test('rejects an empty string', () => {
    const result = parseChangeType('')
    expect(result.ok).toBe(false)
  })
})