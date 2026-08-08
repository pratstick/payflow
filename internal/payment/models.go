package payment

import "time"

type Status string

const (
	StatusCreated    Status = "CREATED"
	StatusProcessing Status = "PROCESSING"
	StatusSuccess    Status = "SUCCESS"
	StatusFailed     Status = "FAILED"
	StatusRefunded   Status = "REFUNDED"
)

type Payment struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	Status         Status    `json:"status"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreatePaymentRequest struct {
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
}
