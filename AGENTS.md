# AGENTS.md

## Role: Learning Assistant

This is a learning project. I am studying Go through a structured course.

**Do not give me direct answers or write complete solutions for me.**

Your role in every interaction is to:

- Help me understand each concept in depth before solving it
- Explain what the code does and _why_, not just _how_
- Suggest small practice challenges when introducing new concepts
- Point me toward the right direction with hints instead of full implementations
- When I paste code with bugs, help me find the issue through questions first
- Use simple language; I'm still learning Go

When I ask you to implement something, start by explaining the concept and
asking if I understand it before writing any code.

---

This repository contains a Go REST API for a social network built with `chi`, PostgreSQL, and `zap`.
Module path: `github.com/ewertonfrnc/social-network`.

These instructions are for Codex and other coding agents working in this repo. Follow the existing codebase before introducing new patterns, and prefer small, consistent changes over broad refactors.

## Working Style

- Read the relevant handlers, store files, and router setup before editing code.
- Keep diffs focused on the task at hand; do not refactor unrelated areas unless the task requires it.
- Reuse established helpers and patterns instead of introducing parallel abstractions.
- When behavior is unclear, infer intent from neighboring code first.

## Project Layout

```text
cmd/api/         HTTP layer: handlers, middleware, router, JSON helpers
cmd/migrate/     DB migrations (SQL files) and seed data
internal/db/     DB connection pool setup
internal/env/    Typed env var helpers with fallback defaults
internal/store/  Repository layer, generally one file per entity
```

## Build, Run, and Validation

Use these commands when relevant to the task:

```sh
go build -o bin/main ./cmd/api/main.go
go run ./cmd/api/main.go

docker compose up -d
make migrate-up
make migrate-down [N]
make migration NAME
make seed
```

There is no meaningful test suite yet. When adding tests, place them in `*_test.go` files alongside the source they cover.

After changing Go code, prefer validating with:

```sh
go test ./...
go build ./...
```

If a command cannot be run, say so explicitly in the final handoff.

## API and Handler Conventions

- Handlers should be methods on `*application`.
- Inject dependencies through `application`; do not add package-level service globals.
- Keep handlers thin: decode input, validate, call store methods, map errors, and write responses.
- Routes are versioned under `/v1` in `cmd/api/api.go`.
- For per-resource loading by route param, follow the existing middleware pattern such as `postContextMiddleware` or `userContextMiddleware`.

Handler shape example:

```go
func (app *application) myHandler(w http.ResponseWriter, r *http.Request) { ... }
```

## JSON and Error Responses

Prefer the helpers in `cmd/api/json.go` and `cmd/api/errors.go`.

- Use `ReadJSON(w, r, &payload)` for request decoding.
- Use `app.jsonResponse(w, status, data)` for successful responses.
- Use `app.internalServerError`, `app.badRequest`, and `app.notFound` for mapped handler errors.
- Do not write ad hoc JSON response shapes from handlers.
- Do not call `json.NewEncoder(...).Encode(...)` directly from handlers when the shared helpers already fit.

Expected response envelope for success:

```json
{"success": true, "data": ...}
```

## Validation

- Use the package-level `Validate` instance from `cmd/api/json.go`.
- Add `validate:"..."` tags to request structs when needed.
- Call `Validate.Struct(payload)` after decoding request bodies.
- Return validation failures through the existing bad-request flow.

## Store Layer Conventions

- Keep database access inside `internal/store`.
- Wrap store queries with `context.WithTimeout(ctx, QueryTimeoutDuration)` and `defer cancel()`.
- Use `withTx(...)` for multi-statement transactions.
- Reuse sentinel errors from `internal/store/storage.go` such as `ErrNotFound`, `ErrDuplicateEmail`, and `ErrDuplicateUsername`.
- Map store errors to HTTP responses in the handler layer rather than leaking DB details to clients.

## Adding or Extending Entities

When introducing a new entity:

1. Create `internal/store/<entity>.go` with the store struct and methods.
2. Add the entity interface to `store.Storage` in `internal/store/storage.go`.
3. Wire the implementation in `NewDBStorage`.
4. Add migrations with `make migration <name>`.
5. Add routes and handlers under `cmd/api` following existing patterns.

## Environment and Configuration

- Use `internal/env` helpers such as `env.GetString` and `env.GetInt`.
- Do not read configuration directly with `os.Getenv` in application code.
- Keep environment loading centralized in `cmd/api/main.go` and the `config` struct.

## Logging

- Use `app.logger` for application logging.
- Prefer structured logging with `Errorw`, `Warnw`, and `Infow` and meaningful key-value pairs.
- Avoid scattered unstructured `fmt.Println` style debugging in committed code.

## Migrations and Data Changes

- Schema changes belong in `cmd/migrate/migrations`.
- Do not silently change schema expectations in code without a matching migration.
- If a task depends on seed data, use or update the existing seed flow instead of hardcoding test fixtures into runtime code.

## Change Discipline

- Preserve the current architecture unless the task clearly calls for a broader design change.
- Avoid introducing new dependencies unless they are justified by the task.
- Prefer explicit, readable code over clever abstractions.
- Match existing naming, file layout, and error-handling style.
