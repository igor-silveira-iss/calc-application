const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

/** Error with a message that can be shown to the user. */
export class CalculatorError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'CalculatorError'
  }
}

type ApiResponse = {
  result?: number
  err?: string
}

/** Sends the operands to an operation endpoint and returns the result. */
export async function calculate(
  endpoint: string,
  body: Record<string, number>,
): Promise<number> {
  let response: Response
  try {
    response = await fetch(`${API_URL}${endpoint}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
  } catch {
    throw new CalculatorError('Could not reach the server. Please try again.')
  }

  const data = (await response.json().catch(() => null)) as ApiResponse | null
  if (!response.ok || typeof data?.result !== 'number') {
    throw new CalculatorError(
      capitalize(data?.err ?? 'Something went wrong. Please try again.'),
    )
  }
  return data.result
}

function capitalize(text: string): string {
  return text.charAt(0).toUpperCase() + text.slice(1)
}
