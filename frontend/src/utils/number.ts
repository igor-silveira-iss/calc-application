/** Converts user input to a number, or returns null if it is not a valid finite number. */
export function parseNumber(value: string): number | null {
  const trimmed = value.trim()
  if (trimmed === '') {
    return null
  }
  const number = Number(trimmed)
  return Number.isFinite(number) ? number : null
}

/**
 * Formats a result for display. It rounds to 12 significant digits to hide
 * floating-point noise (0.1 + 0.2 shows 0.3) and shows -0 as 0.
 */
export function formatNumber(value: number): string {
  return String(Number(value.toPrecision(12)))
}
