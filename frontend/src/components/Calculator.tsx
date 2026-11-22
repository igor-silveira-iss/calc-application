import { useState, type FormEvent } from 'react'
import { calculate, CalculatorError } from '../api/calculator'
import { OPERATIONS, type Operation } from '../operations'
import { formatNumber, parseNumber } from '../utils/number'

type Result = {
  expression: string
  value: string
}

const EMPTY_VALUES = ['', '']

function Calculator() {
  const [operation, setOperation] = useState<Operation>(OPERATIONS[0])
  const [values, setValues] = useState<string[]>(EMPTY_VALUES)
  const [fieldErrors, setFieldErrors] = useState<string[]>([])
  const [result, setResult] = useState<Result | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  function selectOperation(next: Operation) {
    setOperation(next)
    setFieldErrors([])
    setResult(null)
    setError(null)
  }

  function changeValue(index: number, value: string) {
    setValues((current) => current.map((v, i) => (i === index ? value : v)))
    setFieldErrors((current) => current.map((e, i) => (i === index ? '' : e)))
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setResult(null)
    setError(null)

    const fieldValues = values.slice(0, operation.fields.length)
    const numbers = fieldValues.map(parseNumber)
    const errors = numbers.map((n) => (n === null ? 'Enter a valid number' : ''))
    setFieldErrors(errors)
    if (errors.some(Boolean)) {
      return
    }

    const body = Object.fromEntries(
      operation.fields.map((field, i) => [field.name, numbers[i] as number]),
    )

    setLoading(true)
    try {
      const value = await calculate(operation.endpoint, body)
      setResult({
        expression: operation.expression(fieldValues.map((v) => v.trim())),
        value: formatNumber(value) + (operation.suffix ?? ''),
      })
    } catch (err) {
      setError(
        err instanceof CalculatorError
          ? err.message
          : 'Something went wrong. Please try again.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="calculator" aria-labelledby="calculator-title">
      <h1 id="calculator-title">Calculator</h1>

      <form onSubmit={handleSubmit} noValidate>
        <fieldset className="calculator__body" disabled={loading}>
          <div className="operations" role="group" aria-label="Operation">
            {OPERATIONS.map((op) => (
              <button
                key={op.id}
                type="button"
                className="operation"
                aria-pressed={op.id === operation.id}
                onClick={() => selectOperation(op)}
              >
                <span className="operation__symbol" aria-hidden="true">
                  {op.symbol}
                </span>
                <span className="operation__name">{op.name}</span>
              </button>
            ))}
          </div>

          <div className="fields">
            {operation.fields.map((field, i) => {
              const id = `field-${i}`
              const fieldError = fieldErrors[i]
              return (
                <div className="field" key={id}>
                  <label htmlFor={id}>{field.label}</label>
                  <input
                    id={id}
                    type="number"
                    step="any"
                    inputMode="decimal"
                    value={values[i]}
                    onChange={(e) => changeValue(i, e.target.value)}
                    aria-invalid={Boolean(fieldError)}
                    aria-describedby={fieldError ? `${id}-error` : undefined}
                  />
                  {fieldError && (
                    <span className="field__error" id={`${id}-error`}>
                      {fieldError}
                    </span>
                  )}
                </div>
              )
            })}
          </div>

          <button type="submit" className="submit">
            {loading ? 'Calculating…' : 'Calculate'}
          </button>
        </fieldset>
      </form>

      <div className="output" aria-live="polite">
        {result && (
          <div className="result">
            <span className="result__expression">{result.expression} =</span>
            <output className="result__value">{result.value}</output>
          </div>
        )}
        {error && (
          <p className="error" role="alert">
            {error}
          </p>
        )}
      </div>
    </section>
  )
}

export default Calculator
