# Simple Bank

A minimal REST API for a banking-style application built in Go. This project demonstrates core backend skills: database design, SQL code generation, HTTP APIs, transactions, and testing.

---

## What I Learned & Implemented

### Database layer

- **Schema design** – Normalized tables: `accounts`, `entries`, `transfers`, and `users` with foreign keys, indexes, and constraints (`owner_currency_key` for one account per currency per owner).
- **Migrations** – Versioned up/down migrations with [golang-migrate](https://github.com/golang-migrate/migrate) (e.g. `000001_init_schema`, `000002_add_users`). Rollback a single step with `down 1`.
- **SQL-first codegen** – [sqlc](https://sqlc.dev/) to generate type-safe Go from SQL (pgx/v5), with `emit_empty_slice` and type overrides for `timestamptz` → `time.Time`.
- **Connection handling** – Single connection pool via `pgxpool`; config loaded from env (e.g. `app.env`) with Viper.

### Business logic & transactions

- **Transfer transaction** – `TransferTx` in a single DB transaction: create transfer record, two entries (debit/credit), and update both account balances. Uses a helper `addMoney` to keep logic clear.
- **Deadlock avoidance** – Consistent lock order by account ID (always update lower ID first) so concurrent A→B and B→A transfers cannot deadlock.
- **Store abstraction** – `Store` embeds sqlc `Queries` and holds the pool; `execTx` runs a callback inside a transaction and commits or rolls back with error wrapping.

### HTTP API

- **Gin server** – REST endpoints with [Gin](https://github.com/gin-gonic/gin): create account (POST), get account by ID (GET), list accounts with pagination (GET with `page_id` / `page_size`).
- **Validation** – Request validation via struct tags (`binding:"required"`, `oneof=USD EUR`, `min=0`, etc.) and `ShouldBindJSON` / `ShouldBindQuery`.
- **Structured errors** – Central `errorResponse(err)` returning JSON `{"error": "..."}` and appropriate status codes (400, 500).

### Testing

- **Table-driven CRUD tests** – Tests for all generated operations: accounts (Create, Get, GetForUpdate, List, Update, AddBalance, Delete), entries (Create, Get, List), transfers (Create, Get, List, Update, Delete). Helpers like `createAccountInTx` and `runTestWithTransaction` keep tests isolated and rolled back.
- **Concurrent transfer test** – `TestTransferTx` runs multiple transfers concurrently and asserts final balances. `TestTransferTxDeadlock` alternates direction (A→B, B→A) to stress-test lock ordering.
- **Test setup** – `TestMain` loads config and creates a shared `testDB` pool; tests use it directly (e.g. for `Store`) or via `runTestWithTransaction` for per-test rollback.

### Tooling & workflow

- **Makefile** – Targets for Postgres (`postgres`, `createdb`, `dropdb`), migrations (`migrateup`, `migratedown`, `migratedown1`), sqlc (`sqlc`), tests (`test`), and running the server (`server`). Env vars (e.g. from `env.sh`) for DB URL and credentials.
- **Config** – `util.LoadConfig(".")` reads `app.env` (or env) and fills `DB_DRIVER`, `DB_SOURCE`, `SERVER_ADDRESS` for main and tests.

---

## Tech stack

| Area        | Choice                |
| ----------- | --------------------- |
| Language    | Go 1.23+              |
| DB          | PostgreSQL            |
| DB driver   | pgx v5 (pool)         |
| SQL codegen | sqlc                  |
| Migrations  | golang-migrate        |
| HTTP        | Gin                   |
| Config      | Viper                 |
| Tests       | testing + testify     |

---

## Quick start

1. **Environment** – Copy `env.example` to `app.env` (or use `env.sh`) and set `DB_SOURCE`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `SERVER_ADDRESS`.
2. **Postgres** – Start container and create DB:
   ```bash
   source env.sh   # or export vars
   make postgres
   make createdb
   ```
3. **Migrations** – Apply migrations:
   ```bash
   make migrateup
   ```
4. **Run** – Start the API server:
   ```bash
   make server
   ```
   Or: `go run main.go` (uses config from current directory).

**Run tests**

```bash
make test
```

**Roll back only the last migration (e.g. users)**

```bash
make migratedown1
```

---

## Project structure

```
simple_bank/
├── api/              # HTTP handlers and server setup
│   ├── server.go     # Gin engine, routes, Start()
│   └── account.go    # createAccount, getAccount, listAccounts
├── db/
│   ├── migration/    # Up/down SQL migrations
│   ├── query/        # SQL queries for sqlc (account, entry, transfer)
│   └── sqlc/        # Generated code + Store and TransferTx
├── util/             # Config loading, random helpers for tests
├── main.go           # Load config, connect DB, create store & server
├── Makefile          # postgres, migrate, sqlc, test, server
├── sqlc.yaml         # sqlc config (pgx, emit_empty_slice, overrides)
└── app.env / env.example
```

---

## API overview

| Method | Path             | Description                    |
| ------ | ----------------- | ------------------------------ |
| POST   | /accounts         | Create account (owner, balance, currency) |
| GET    | /accounts/:id     | Get account by ID              |
| GET    | /accounts         | List accounts (query: page_id, page_size) |

---

## License

MIT (or your choice).
