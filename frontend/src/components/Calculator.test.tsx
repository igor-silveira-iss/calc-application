import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { calculate, CalculatorError } from '../api/calculator'
import Calculator from './Calculator'

vi.mock('../api/calculator', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../api/calculator')>()),
  calculate: vi.fn(),
}))

const calculateMock = vi.mocked(calculate)

beforeEach(() => {
  calculateMock.mockReset()
})

function setup() {
  const user = userEvent.setup()
  render(<Calculator />)
  return user
}

const operationButton = (name: string) => screen.getByRole('button', { name })
const submitButton = () => screen.getByRole('button', { name: /^calculat/i })

describe('Calculator', () => {
  it('shows every operation with addition selected', () => {
    setup()

    const buttons = screen.getByRole('group', { name: 'Operation' })
    expect(buttons.querySelectorAll('button')).toHaveLength(7)
    expect(operationButton('Add')).toHaveAttribute('aria-pressed', 'true')
    expect(operationButton('Divide')).toHaveAttribute('aria-pressed', 'false')
    expect(screen.getByLabelText('First number')).toBeInTheDocument()
    expect(screen.getByLabelText('Second number')).toBeInTheDocument()
  })

  it('changes the fields when another operation is selected', async () => {
    const user = setup()

    await user.click(operationButton('Divide'))
    expect(operationButton('Divide')).toHaveAttribute('aria-pressed', 'true')
    expect(operationButton('Add')).toHaveAttribute('aria-pressed', 'false')
    expect(screen.getByLabelText('Dividend')).toBeInTheDocument()
    expect(screen.getByLabelText('Divisor')).toBeInTheDocument()

    await user.click(operationButton('Root'))
    expect(screen.getByLabelText('Number')).toBeInTheDocument()
    expect(screen.getAllByRole('spinbutton')).toHaveLength(1)
  })

  it('calculates and shows the result', async () => {
    calculateMock.mockResolvedValue(5)
    const user = setup()

    await user.type(screen.getByLabelText('First number'), '2')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(submitButton())

    expect(calculateMock).toHaveBeenCalledWith('/addition', { number1: 2, number2: 3 })
    expect(await screen.findByText('2 + 3 =')).toBeInTheDocument()
    expect(screen.getByRole('status')).toHaveTextContent('5')
  })

  it('submits with the Enter key', async () => {
    calculateMock.mockResolvedValue(-6)
    const user = setup()

    await user.click(operationButton('Subtract'))
    await user.type(screen.getByLabelText('First number'), '-1')
    await user.type(screen.getByLabelText('Second number'), '5{Enter}')

    expect(calculateMock).toHaveBeenCalledWith('/subtraction', { number1: -1, number2: 5 })
    expect(await screen.findByText('-1 − 5 =')).toBeInTheDocument()
  })

  it('sends a single number for the square root', async () => {
    calculateMock.mockResolvedValue(4)
    const user = setup()

    await user.click(operationButton('Root'))
    await user.type(screen.getByLabelText('Number'), '16')
    await user.click(submitButton())

    expect(calculateMock).toHaveBeenCalledWith('/square-root', { number: 16 })
    expect(await screen.findByText('√16 =')).toBeInTheDocument()
  })

  it('shows percentage results with a percent sign', async () => {
    calculateMock.mockResolvedValue(12.5)
    const user = setup()

    await user.click(operationButton('Percent'))
    await user.type(screen.getByLabelText('Value'), '25')
    await user.type(screen.getByLabelText('Total'), '200')
    await user.click(submitButton())

    expect(await screen.findByRole('status')).toHaveTextContent('12.5%')
  })

  it('rounds floating-point noise in the result', async () => {
    calculateMock.mockResolvedValue(0.30000000000000004)
    const user = setup()

    await user.type(screen.getByLabelText('First number'), '0.1')
    await user.type(screen.getByLabelText('Second number'), '0.2')
    await user.click(submitButton())

    expect(await screen.findByRole('status')).toHaveTextContent(/^0\.3$/)
  })

  it('validates empty fields without calling the API', async () => {
    const user = setup()

    await user.click(submitButton())

    expect(screen.getAllByText('Enter a valid number')).toHaveLength(2)
    expect(screen.getByLabelText('First number')).toHaveAttribute('aria-invalid', 'true')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('clears a field error when the field changes', async () => {
    const user = setup()

    await user.click(submitButton())
    await user.type(screen.getByLabelText('First number'), '4')

    expect(screen.getAllByText('Enter a valid number')).toHaveLength(1)
    expect(screen.getByLabelText('First number')).toHaveAttribute('aria-invalid', 'false')
  })

  it('shows the API error message', async () => {
    calculateMock.mockRejectedValue(new CalculatorError('Division by zero is not allowed'))
    const user = setup()

    await user.click(operationButton('Divide'))
    await user.type(screen.getByLabelText('Dividend'), '1')
    await user.type(screen.getByLabelText('Divisor'), '0')
    await user.click(submitButton())

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Division by zero is not allowed',
    )
    expect(screen.queryByRole('status')).not.toBeInTheDocument()
  })

  it('shows a generic message for unexpected errors', async () => {
    calculateMock.mockRejectedValue(new Error('boom'))
    const user = setup()

    await user.type(screen.getByLabelText('First number'), '1')
    await user.type(screen.getByLabelText('Second number'), '2')
    await user.click(submitButton())

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Something went wrong. Please try again.',
    )
  })

  it('disables the form while calculating', async () => {
    let resolve: (value: number) => void = () => {}
    calculateMock.mockReturnValue(new Promise((r) => (resolve = r)))
    const user = setup()

    await user.type(screen.getByLabelText('First number'), '2')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(submitButton())

    expect(submitButton()).toHaveTextContent('Calculating…')
    expect(submitButton()).toBeDisabled()
    expect(operationButton('Divide')).toBeDisabled()

    resolve(5)
    expect(await screen.findByRole('status')).toHaveTextContent('5')
    expect(submitButton()).toHaveTextContent('Calculate')
    expect(submitButton()).toBeEnabled()
  })

  it('clears the result when another operation is selected', async () => {
    calculateMock.mockResolvedValue(5)
    const user = setup()

    await user.type(screen.getByLabelText('First number'), '2')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(submitButton())
    expect(await screen.findByRole('status')).toBeInTheDocument()

    await user.click(operationButton('Multiply'))

    expect(screen.queryByRole('status')).not.toBeInTheDocument()
    expect(screen.getByLabelText('First number')).toHaveValue(2)
  })
})
