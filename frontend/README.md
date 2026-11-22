# Calculator Frontend

React + TypeScript app for the calculator. The user picks an operation, types
the numbers and sees the result calculated by the [Go API](../api/README.md).

For running the whole project, see the [main README](../README.md).

## Contents

- [Requirements](#requirements)
- [Setup and running](#setup-and-running)
- [Scripts](#scripts)
- [Tests and coverage](#tests-and-coverage)
- [Project structure](#project-structure)
- [Design decisions](#design-decisions)

## Requirements

- [Node.js](https://nodejs.org/) 22.13+ (or 20.19+)
- The API running. See [`api/README.md`](../api/README.md).

## Setup and running

```bash
cd frontend
npm install
npm run dev
```

Open http://localhost:5173.

The app calls the API at `http://localhost:8080` by default. To use another
address, create `frontend/.env.local`:

```bash
VITE_API_URL=http://localhost:5000
```

Restart `npm run dev` after changing it. If the frontend runs on a port other
than `5173`, also set `ALLOWED_ORIGIN` in `api/.env`, or the browser will
block the requests (CORS).

In Docker, `VITE_API_URL` is set to `/api` and nginx forwards those requests
to the API.

## Scripts

| Command              | What it does                                  |
| -------------------- | --------------------------------------------- |
| `npm run dev`        | Starts the development server                 |
| `npm run build`      | Type-checks and builds to `dist/`             |
| `npm run preview`    | Serves the built app                          |
| `npm run lint`       | Runs ESLint                                   |
| `npm test`           | Runs the tests once                           |
| `npm run test:watch` | Runs the tests on every change                |
| `npm run coverage`   | Runs the tests with a coverage report         |

## Tests and coverage

The tests use [Vitest](https://vitest.dev/) with
[Testing Library](https://testing-library.com/), which drives the components
the way a user would: clicking buttons and typing into labeled fields.

```bash
npm run coverage
```

It prints the table below and writes an HTML report to
`coverage/index.html`.

| File                        | Tests | Lines | Branches |
| --------------------------- | ----- | ----- | -------- |
| `utils/number.ts`           | 19    | 100%  | 100%     |
| `operations.ts`             | 10    | 100%  | 100%     |
| `api/calculator.ts`         | 7     | 100%  | 87.5%    |
| `components/Calculator.tsx` | 13    | 100%  | 100%     |
| `App.tsx`                   | 1     | 100%  | 100%     |
| **Total**                   | **50**| **100%** | **97%** |

The only branch not covered is the default API address in
`api/calculator.ts`, which is used or not depending on whether a
`.env.local` file exists on the machine.

What the tests check:

- **`number.test.ts`:** parsing input (empty, spaces, text, numbers too
  large) and formatting results (`0.1 + 0.2` shows `0.3`, `-0` shows `0`).
- **`operations.test.ts`:** each operation's endpoint, field names and
  expression match the API.
- **`calculator.test.ts`:** the request sent to the API, and how successes,
  API errors, network failures and invalid responses are handled. `fetch`
  is mocked.
- **`Calculator.test.tsx`:** choosing an operation, the field labels,
  submitting with the button or Enter, validation messages, API errors, the
  loading state and clearing the result. The API client is mocked.

## Project structure

```
frontend/
├── Dockerfile            # Builds the app and serves it with nginx
├── nginx.conf            # Serves the app, forwards /api to the API
├── vite.config.ts        # Vite and Vitest configuration
└── src/
    ├── main.tsx          # Entry point
    ├── App.tsx           # Page layout
    ├── App.css           # Calculator styles
    ├── index.css         # Colors (light and dark) and base styles
    ├── operations.ts     # The operations shown in the UI
    ├── api/
    │   └── calculator.ts # API client
    ├── components/
    │   └── Calculator.tsx
    ├── utils/
    │   └── number.ts     # Input parsing and result formatting
    └── test/
        └── setup.ts      # Test setup
```

## Design decisions

- **Operations are data.** `operations.ts` describes each operation: its
  endpoint, symbol, field labels and how to show the expression. The
  component renders whatever is in that list, so adding an operation is a
  new entry there and no new UI code.
- **The API client is isolated.** `api/calculator.ts` is the only file that
  knows about HTTP. It returns the result or throws a `CalculatorError` with
  a message ready to show, so the component never deals with status codes.
- **Validation before sending.** Empty or invalid fields show an error under
  the field and no request is made. The API still validates everything on
  its own.
- **Readable results.** Results are rounded to 12 significant digits, which
  hides floating-point noise (`0.1 + 0.2` shows `0.3`) and shows `-0` as `0`.
- **Simple state.** The component uses `useState` only. The app has a single
  screen and no shared state, so a state library would add complexity
  without benefit.
- **Responsive without breakpoints for the buttons.** The operation buttons
  use a CSS grid with `auto-fit`, so they fill one row on wide screens and
  wrap on narrow ones, whatever the font or screen size.
- **Accessible.** Every field has a label, the selected operation uses
  `aria-pressed`, invalid fields use `aria-invalid`, and results and errors
  are announced to screen readers.
- **Light and dark themes.** Colors are CSS variables that follow the
  system theme.
