import { afterEach, describe, expect, it, vi } from 'vitest'
import { calculate, CalculatorError } from './calculator'

/** Replaces fetch with a mock that answers with the given status and body. */
function mockFetch(status: number, body: string) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(body, { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('calculate', () => {
  it('posts the operands as JSON and returns the result', async () => {
    const fetchMock = mockFetch(200, '{"result":5}')

    await expect(calculate('/addition', { number1: 2, number2: 3 })).resolves.toBe(5)

    expect(fetchMock).toHaveBeenCalledOnce()
    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toMatch(/\/addition$/)
    expect(init).toEqual({
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: '{"number1":2,"number2":3}',
    })
  })

  it('returns a zero result', async () => {
    mockFetch(200, '{"result":0}')
    await expect(calculate('/addition', { number1: 0, number2: 0 })).resolves.toBe(0)
  })

  it('throws the API error message, capitalized', async () => {
    mockFetch(422, '{"err":"division by zero is not allowed"}')

    const promise = calculate('/division', { number1: 1, number2: 0 })

    await expect(promise).rejects.toBeInstanceOf(CalculatorError)
    await expect(promise).rejects.toThrow('Division by zero is not allowed')
  })

  it('throws a connection error when the server cannot be reached', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    await expect(calculate('/addition', { number1: 1, number2: 2 })).rejects.toThrow(
      'Could not reach the server. Please try again.',
    )
  })

  it.each([
    ['a response that is not JSON', 502, '<html>Bad Gateway</html>'],
    ['an error without a message', 500, '{}'],
    ['a success without a result', 200, '{}'],
  ])('throws a generic error for %s', async (_, status, body) => {
    mockFetch(status, body)

    await expect(calculate('/addition', { number1: 1, number2: 2 })).rejects.toThrow(
      'Something went wrong. Please try again.',
    )
  })
})
