package models

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	StatusRegistered = "REGISTERED"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
)

type Order struct {
	Login   string
	Status  string           `json:"status"`
	Accrual *decimal.Decimal `json:"accrual,omitempty"`
	OrderID string           `json:"order"`
}

type OrderItem struct {
	OrderID  string           `json:"number"`
	Status   string           `json:"status"`
	Accrual  *decimal.Decimal `json:"accrual,omitempty"`
	Uploaded time.Time        `json:"uploaded_at"`
}
