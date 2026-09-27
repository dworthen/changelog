import { CHANGE_TYPES, type ChangeType } from './types'

export type ParseChangeTypeResult =
  | { ok: false; error: string }
  | { ok: true; value: ChangeType }

export type ParseChangeType = (value: string) => ParseChangeTypeResult

export const parseChangeType: ParseChangeType = (value) => {
  const normalized = value.trim().toLowerCase()
  const match = (Object.keys(CHANGE_TYPES) as ChangeType[]).find(
    (changeType) => changeType.toLowerCase() === normalized,
  )

  if (match == null) {
    const valid = Object.keys(CHANGE_TYPES).join(', ')
    return {
      ok: false,
      error: `Invalid change type: '${value}'. Valid types are: ${valid}.`,
    }
  }

  return { ok: true, value: match }
}