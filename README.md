# PayFlow

PayFlow is a production-style payment processing foundation implemented in Go with PostgreSQL for persistence.

## Architecture

```text
HTTP Handler -> Service -> Repository -> PostgreSQL
```

## Current Features (Phase 1)

- Payment create and fetch APIs
- PostgreSQL-backed payment storage
- Payment state model and transition validation
- Idempotency key uniqueness at database level
- Unit tests, race test support, and CI workflow
- Docker and Docker Compose for local API/PostgreSQL setup

## Local Setup

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```
2. Start PostgreSQL:
   ```bash
   docker compose up -d postgres
   ```
3. Run migrations manually against your PostgreSQL instance from files in `/migrations`.
4. Run API:
   ```bash
   export $(cat .env | xargs)
   go run ./cmd/api
   ```

## Environment Variables

- `PORT` (default `8080`)
- `DATABASE_URL` (required)

## Running PostgreSQL

```bash
docker compose up -d postgres
```

## Running the API

```bash
go run ./cmd/api
```

## Running Tests

```bash
go test ./...
go test -race ./...
go vet ./...
```

## API Examples

Create payment:

```bash
curl -X POST http://localhost:8080/v1/payments \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: idem-123' \
  -d '{"customer_id":"cust_123","amount":1000,"currency":"INR"}'
```

Get payment:

```bash
curl http://localhost:8080/v1/payments/<payment_id>
```

## Current Limitations

Kafka, Redis, asynchronous workers, retries, DLQ, webhooks, refunds processing, reconciliation, and observability are intentionally out of scope for Phase 1.

## Planned Architecture (Future Phases)

Future phases will add Kafka eventing, Redis-based idempotency/caching, asynchronous workers, retry orchestration, DLQ handling, webhooks, refunds, reconciliation, and observability.
