package payment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment Payment) (Payment, error)
	GetPaymentByID(ctx context.Context, id string) (Payment, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreatePayment(ctx context.Context, payment Payment) (Payment, error) {
	const q = `
		INSERT INTO payments (customer_id, amount, currency, status, idempotency_key)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, customer_id, amount, currency, status, idempotency_key, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, q, payment.CustomerID, payment.Amount, payment.Currency, payment.Status, payment.IdempotencyKey)

	var created Payment
	if err := row.Scan(
		&created.ID,
		&created.CustomerID,
		&created.Amount,
		&created.Currency,
		&created.Status,
		&created.IdempotencyKey,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return Payment{}, ErrDuplicateIdempotencyKey
		}
		return Payment{}, err
	}
	return created, nil
}

func (r *PostgresRepository) GetPaymentByID(ctx context.Context, id string) (Payment, error) {
	const q = `
		SELECT id, customer_id, amount, currency, status, idempotency_key, created_at, updated_at
		FROM payments WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, q, id)
	var payment Payment
	if err := row.Scan(
		&payment.ID,
		&payment.CustomerID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.IdempotencyKey,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Payment{}, ErrPaymentNotFound
		}
		return Payment{}, err
	}
	return payment, nil
}
