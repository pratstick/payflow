package payment

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidRequest          = errors.New("invalid request")
	ErrMissingIdempotencyKey   = errors.New("missing idempotency key")
	ErrDuplicateIdempotencyKey = errors.New("duplicate idempotency key")
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrInvalidStateTransition  = errors.New("invalid state transition")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePayment(ctx context.Context, req CreatePaymentRequest, idempotencyKey string) (Payment, error) {
	if strings.TrimSpace(idempotencyKey) == "" {
		return Payment{}, ErrMissingIdempotencyKey
	}
	if strings.TrimSpace(req.CustomerID) == "" || req.Amount <= 0 || len(strings.TrimSpace(req.Currency)) != 3 {
		return Payment{}, ErrInvalidRequest
	}

	payment := Payment{
		CustomerID:     req.CustomerID,
		Amount:         req.Amount,
		Currency:       strings.ToUpper(req.Currency),
		Status:         StatusCreated,
		IdempotencyKey: idempotencyKey,
	}

	return s.repo.CreatePayment(ctx, payment)
}

func (s *Service) GetPayment(ctx context.Context, id string) (Payment, error) {
	if strings.TrimSpace(id) == "" {
		return Payment{}, ErrInvalidRequest
	}
	return s.repo.GetPaymentByID(ctx, id)
}
