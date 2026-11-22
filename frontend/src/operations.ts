/** An input field of an operation. `name` is the key sent to the API. */
export type Field = {
  name: string
  label: string
}

/** A calculator operation and how the UI shows it. */
export type Operation = {
  id: string
  name: string
  symbol: string
  endpoint: string
  fields: Field[]
  /** Builds the text shown above the result, e.g. "2 + 3". */
  expression: (values: string[]) => string
  /** Text appended to the result, e.g. "%". */
  suffix?: string
}

const twoNumbers: Field[] = [
  { name: 'number1', label: 'First number' },
  { name: 'number2', label: 'Second number' },
]

export const OPERATIONS: Operation[] = [
  {
    id: 'addition',
    name: 'Add',
    symbol: '+',
    endpoint: '/addition',
    fields: twoNumbers,
    expression: ([a, b]) => `${a} + ${b}`,
  },
  {
    id: 'subtraction',
    name: 'Subtract',
    symbol: '−',
    endpoint: '/subtraction',
    fields: twoNumbers,
    expression: ([a, b]) => `${a} − ${b}`,
  },
  {
    id: 'multiplication',
    name: 'Multiply',
    symbol: '×',
    endpoint: '/multiplication',
    fields: twoNumbers,
    expression: ([a, b]) => `${a} × ${b}`,
  },
  {
    id: 'division',
    name: 'Divide',
    symbol: '÷',
    endpoint: '/division',
    fields: [
      { name: 'number1', label: 'Dividend' },
      { name: 'number2', label: 'Divisor' },
    ],
    expression: ([a, b]) => `${a} ÷ ${b}`,
  },
  {
    id: 'exponentiation',
    name: 'Power',
    symbol: 'xʸ',
    endpoint: '/exponentiation',
    fields: [
      { name: 'number1', label: 'Base' },
      { name: 'number2', label: 'Exponent' },
    ],
    expression: ([a, b]) => `${a} ^ ${b}`,
  },
  {
    id: 'square-root',
    name: 'Root',
    symbol: '√',
    endpoint: '/square-root',
    fields: [{ name: 'number', label: 'Number' }],
    expression: ([a]) => `√${a}`,
  },
  {
    id: 'percentage',
    name: 'Percent',
    symbol: '%',
    endpoint: '/percentage',
    fields: [
      { name: 'number1', label: 'Value' },
      { name: 'number2', label: 'Total' },
    ],
    expression: ([a, b]) => `${a} of ${b}`,
    suffix: '%',
  },
]
