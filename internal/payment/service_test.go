package payment

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepo struct {
	createFn func(ctx context.Context, payment Payment) (Payment, error)
	getFn    func(ctx context.Context, id string) (Payment, error)
}

func (f fakeRepo) CreatePayment(ctx context.Context, payment Payment) (Payment, error) {
	return f.createFn(ctx, payment)
}

func (f fakeRepo) GetPaymentByID(ctx context.Context, id string) (Payment, error) {
	return f.getFn(ctx, id)
}

func TestServiceCreatePayment(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	tests := []struct {
		name           string
		req            CreatePaymentRequest
		idempotencyKey string
		repoErr        error
		wantErr        error
	}{
		{
			name:           "valid payment creation",
			req:            CreatePaymentRequest{CustomerID: "cust_123", Amount: 1000, Currency: "inr"},
			idempotencyKey: "idem-1",
		},
		{
			name:           "invalid payment input",
			req:            CreatePaymentRequest{CustomerID: "", Amount: 1000, Currency: "INR"},
			idempotencyKey: "idem-2",
			wantErr:        ErrInvalidRequest,
		},
		{
			name:    "missing idempotency key",
			req:     CreatePaymentRequest{CustomerID: "cust_123", Amount: 1000, Currency: "INR"},
			wantErr: ErrMissingIdempotencyKey,
		},
		{
			name:           "duplicate idempotency behavior",
			req:            CreatePaymentRequest{CustomerID: "cust_123", Amount: 1000, Currency: "INR"},
			idempotencyKey: "idem-dup",
			repoErr:        ErrDuplicateIdempotencyKey,
			wantErr:        ErrDuplicateIdempotencyKey,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := NewService(fakeRepo{
				createFn: func(_ context.Context, payment Payment) (Payment, error) {
					if tc.repoErr != nil {
						return Payment{}, tc.repoErr
					}
					payment.ID = "pay_1"
					payment.CreatedAt = now
					payment.UpdatedAt = now
					return payment, nil
				},
				getFn: func(_ context.Context, _ string) (Payment, error) {
					return Payment{}, nil
				},
			})

			got, err := service.CreatePayment(context.Background(), tc.req, tc.idempotencyKey)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Status != StatusCreated {
				t.Fatalf("expected CREATED status, got %s", got.Status)
			}
			if got.Currency != "INR" {
				t.Fatalf("expected uppercase currency, got %s", got.Currency)
			}
		})
	}
}
