package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlerCreatePayment_MissingIdempotency(t *testing.T) {
	t.Parallel()
	svc := NewService(fakeRepo{
		createFn: func(_ context.Context, _ Payment) (Payment, error) { return Payment{}, nil },
		getFn:    func(_ context.Context, _ string) (Payment, error) { return Payment{}, nil },
	})
	h := NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := []byte(`{"customer_id":"cust_1","amount":1000,"currency":"INR"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/payments", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHandlerGetPayment_NotFound(t *testing.T) {
	t.Parallel()
	svc := NewService(fakeRepo{
		createFn: func(_ context.Context, _ Payment) (Payment, error) { return Payment{}, nil },
		getFn: func(_ context.Context, _ string) (Payment, error) {
			return Payment{}, ErrPaymentNotFound
		},
	})
	h := NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/v1/payments/pay_1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandlerCreatePayment_Success(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	svc := NewService(fakeRepo{
		createFn: func(_ context.Context, p Payment) (Payment, error) {
			if p.IdempotencyKey == "dup" {
				return Payment{}, ErrDuplicateIdempotencyKey
			}
			p.ID = "pay_1"
			p.CreatedAt = now
			p.UpdatedAt = now
			return p, nil
		},
		getFn: func(_ context.Context, _ string) (Payment, error) { return Payment{}, errors.New("unused") },
	})
	h := NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	payload := CreatePaymentRequest{CustomerID: "cust_1", Amount: 1000, Currency: "INR"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/payments", bytes.NewReader(b))
	req.Header.Set("Idempotency-Key", "ok")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}
