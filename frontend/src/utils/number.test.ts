import { describe, expect, it } from 'vitest'
import { formatNumber, parseNumber } from './number'

describe('parseNumber', () => {
  it.each([
    ['42', 42],
    ['-8', -8],
    ['0', 0],
    ['3.14', 3.14],
    ['  7  ', 7],
    ['1e3', 1000],
  ])('parses %j as %d', (input, expected) => {
    expect(parseNumber(input)).toBe(expected)
  })

  it.each(['', '   ', 'abc', '1,5', 'Infinity', '1e400'])(
    'returns null for %j',
    (input) => {
      expect(parseNumber(input)).toBeNull()
    },
  )
})

describe('formatNumber', () => {
  it.each([
    [5, '5'],
    [12.5, '12.5'],
    [-3, '-3'],
    [0.1 + 0.2, '0.3'],
    [-0, '0'],
    [1 / 3, '0.333333333333'],
    [1e21, '1e+21'],
  ])('formats %d as %j', (input, expected) => {
    expect(formatNumber(input)).toBe(expected)
  })
})
