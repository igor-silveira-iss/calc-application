# Calculator API

REST API written in Go that powers the calculator application. It exposes
basic and advanced arithmetic operations as JSON endpoints, validates every
request and reports errors in a consistent JSON format.

## Contents

- [Requirements](#requirements)
- [Setup](#setup)
- [Running](#running)
- [API usage](#api-usage)
- [Tests and coverage](#tests-and-coverage)
- [Project structure](#project-structure)
- [Design decisions](#design-decisions)
- [Assumptions and limitations](#assumptions-and-limitations)
- [Use of AI](#use-of-ai)

## Requirements

- [Go](https://go.dev/dl/) 1.27 or later
- `make` (optional, for the shortcuts below)

## Setup

```bash
cd api
cp .env.example .env
go mod download
```

The `.env` file holds the settings below. It is optional: without it, the
values come from the environment or fall back to the defaults.

| Variable         | Description                                 | Default                 |
| ---------------- | ------------------------------------------- | ----------------------- |
| `API_PORT`       | Port the server listens on                  | `8080`                  |
| `ALLOWED_ORIGIN` | Frontend origin allowed by CORS             | `http://localhost:5173` |

Variables already set in the environment take precedence over the `.env` file.

## Running

```bash
make run      # or: go run .
```

The API is then available at `http://localhost:8080`.

To run the API together with the frontend in Docker, see the
[Docker section of the main README](../README.md#quick-start-with-docker).

## API usage

Every operation is a `POST` request with a JSON body and returns a JSON body.

### Endpoints

| Endpoint                | Body                               | Result                      |
| ----------------------- | ---------------------------------- | --------------------------- |
| `POST /addition`        | `{"number1": a, "number2": b}`     | `a + b`                     |
| `POST /subtraction`     | `{"number1": a, "number2": b}`     | `a - b`                     |
| `POST /multiplication`  | `{"number1": a, "number2": b}`     | `a * b`                     |
| `POST /division`        | `{"number1": a, "number2": b}`     | `a / b`                     |
| `POST /exponentiation`  | `{"number1": a, "number2": b}`     | `a` raised to `b`           |
| `POST /percentage`      | `{"number1": a, "number2": b}`     | what percent `a` is of `b`  |
| `POST /square-root`     | `{"number": a}`                    | square root of `a`          |

A successful response has status `200` and the body `{"result": <number>}`.
An error response has the body `{"err": "<message>"}`.

### Examples

```bash
curl -X POST http://localhost:8080/addition \
  -H "Content-Type: application/json" \
  -d '{"number1": 2, "number2": 3}'
# 200 {"result":5}

curl -X POST http://localhost:8080/percentage \
  -H "Content-Type: application/json" \
  -d '{"number1": 25, "number2": 200}'
# 200 {"result":12.5}

curl -X POST http://localhost:8080/square-root \
  -H "Content-Type: application/json" \
  -d '{"number": 16}'
# 200 {"result":4}

curl -X POST http://localhost:8080/division \
  -H "Content-Type: application/json" \
  -d '{"number1": 1, "number2": 0}'
# 422 {"err":"division by zero is not allowed"}

curl -X POST http://localhost:8080/addition \
  -H "Content-Type: application/json" \
  -d '{"number1": "abc", "number2": 1}'
# 400 {"err":"number1 must be a valid number"}
```

### Status codes

| Status | When                                                      | Example message                                           |
| ------ | --------------------------------------------------------- | --------------------------------------------------------- |
| `200`  | The operation succeeded                                   | —                                                         |
| `400`  | The request body is invalid                               | `number1 is required`, `request body must be valid JSON`  |
| `404`  | The route does not exist                                  | `route not found`                                         |
| `405`  | The route exists but the method is not `POST`             | `method not allowed`                                      |
| `413`  | The request body is larger than 1024 bytes                | `request body must not exceed 1024 bytes`                 |
| `422`  | The input is valid but the result cannot be calculated    | see the table below                                       |

### Edge cases

| Case                                         | Status | Message                                                  |
| -------------------------------------------- | ------ | -------------------------------------------------------- |
| Division or percentage by zero (including `0/0`) | `422` | `division by zero is not allowed`                        |
| Square root of a negative number             | `422`  | `square root of a negative number is not allowed`        |
| Zero raised to a negative power              | `422`  | `zero cannot be raised to a negative power`              |
| Negative number raised to a fractional power | `422`  | `negative number cannot be raised to a fractional power` |
| Result larger than `float64` can hold        | `422`  | `result exceeds the supported numeric range`             |
| Missing or `null` field                      | `400`  | `number1 is required`                                    |
| Unknown field, e.g. `{"num1": 5}`            | `400`  | `unknown field "num1"`                                   |
| Field that is not a number                   | `400`  | `number1 must be a valid number`                         |
| Empty body                                   | `400`  | `request body is required`                               |
| Array instead of an object                   | `400`  | `request body must be a JSON object`                     |
| More than one JSON object                    | `400`  | `request body must contain a single JSON object`         |

## Tests and coverage

```bash
make test       # or: go test ./...
make coverage   # coverage per function + coverage.html
```

`make coverage` prints the coverage of each function and writes
`coverage.html`, which highlights covered lines in the browser.

The suite has 99 test cases, written as table-driven tests, and covers
**96.8%** of the statements:

| Package                                         | Coverage |
| ----------------------------------------------- | -------- |
| `controllers` (operations, handlers, decoding)  | 100%     |
| `middleware`, `response`, `router`, `routes`    | 100%     |
| `config`                                        | 100%     |
| `main`                                          | 0%       |

The only statements not covered are in `main`, which just starts the server.

## Project structure

```
api/
├── main.go                     # Entry point: loads config, starts the server
├── Makefile                    # run, test and coverage shortcuts
└── src/
    ├── config/                 # Reads settings from .env
    ├── controllers/
    │   ├── mathOperations.go   # Operations and their rules (no HTTP)
    │   ├── handlers.go         # HTTP handlers, one per endpoint
    │   └── request.go          # JSON body decoding and validation
    ├── middleware/             # CORS and request logging
    ├── models/                 # Operation input and output types
    ├── response/               # JSON and error responses
    └── router/
        ├── router.go           # Router, 404 and 405 responses
        └── routes/             # Route table
```

A request goes through these steps:

```
CORS → router → Logger → handler → decodeBody → operation → JSON response
```

## Design decisions

- **Operations are separate from HTTP.** The functions in
  `mathOperations.go` receive numbers and return a result or an error. They
  know nothing about requests, so they are easy to test and reuse. The
  handlers only decode the request, call an operation and write the response.
- **One shared handler flow.** All two-number endpoints use the same
  `handleOperands` function, so validation and error handling live in one
  place instead of being repeated in every handler.
- **Small set of dependencies.** Only
  [gorilla/mux](https://github.com/gorilla/mux) for routing and
  [godotenv](https://github.com/joho/godotenv) for the `.env` file. Everything
  else, including CORS, uses the standard library.
- **`POST` with a JSON body for every operation.** It keeps all endpoints
  consistent and lets the request body be validated the same way everywhere.
- **Required fields.** The request types use `*float64`, so a missing or
  `null` field is rejected instead of being silently treated as `0`.
- **Strict request bodies.** Unknown fields, extra JSON values and bodies
  larger than 1024 bytes are rejected. Error messages describe the problem
  without exposing Go internals.
- **`400` versus `422`.** `400` means the request itself is wrong. `422`
  means the request is well formed but the math is not possible, such as a
  division by zero.
- **Consistent error format.** Every error, including `404` and `405`, is
  returned as `{"err": "<message>"}`, so the frontend handles all of them the
  same way.
- **CORS for one origin.** Only the frontend origin set in `ALLOWED_ORIGIN`
  may call the API from the browser, instead of allowing any site with `*`.

## Assumptions and limitations

- Numbers are `float64`, so results follow floating-point rules:
  `0.1 + 0.2` returns `0.30000000000000004`, and `-1 * 0` returns `-0`.
  Rounding for display is left to the frontend.
- `0` raised to `0` returns `1`, following the usual mathematical convention.
- Percentage answers "what percent of `number2` is `number1`", e.g. `25` of
  `200` is `12.5`.
- The API has no authentication or rate limiting, since it only performs
  stateless calculations.

## Use of AI

This project was built with the help of Claude Code. What it was used for,
and the prompts, are described in the
[main README](../README.md#use-of-ai).
