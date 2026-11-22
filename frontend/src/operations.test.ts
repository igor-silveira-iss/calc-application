import { describe, expect, it } from 'vitest'
import { OPERATIONS } from './operations'

function find(id: string) {
  const operation = OPERATIONS.find((op) => op.id === id)
  if (!operation) {
    throw new Error(`operation ${id} not found`)
  }
  return operation
}

describe('OPERATIONS', () => {
  it('has one entry per API endpoint', () => {
    expect(OPERATIONS.map((op) => op.endpoint)).toEqual([
      '/addition',
      '/subtraction',
      '/multiplication',
      '/division',
      '/exponentiation',
      '/square-root',
      '/percentage',
    ])
  })

  it('uses the field names the API expects', () => {
    for (const op of OPERATIONS) {
      const names = op.fields.map((field) => field.name)
      const expected = op.id === 'square-root' ? ['number'] : ['number1', 'number2']
      expect(names, op.id).toEqual(expected)
    }
  })

  it.each([
    ['addition', ['2', '3'], '2 + 3'],
    ['subtraction', ['2', '3'], '2 − 3'],
    ['multiplication', ['2', '3'], '2 × 3'],
    ['division', ['2', '3'], '2 ÷ 3'],
    ['exponentiation', ['2', '3'], '2 ^ 3'],
    ['square-root', ['16'], '√16'],
    ['percentage', ['25', '200'], '25 of 200'],
  ])('builds the %s expression', (id, values, expected) => {
    expect(find(id).expression(values)).toBe(expected)
  })

  it('adds a percent sign only to percentage results', () => {
    const withSuffix = OPERATIONS.filter((op) => op.suffix).map((op) => op.id)
    expect(withSuffix).toEqual(['percentage'])
    expect(find('percentage').suffix).toBe('%')
  })
})
