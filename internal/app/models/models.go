package models

import (
	"github.com/shopspring/decimal"
	"time"
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

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdrawal struct {
	OrderID   string           `json:"order"`
	Sum       *decimal.Decimal `json:"sum"`
	Processed time.Time        `json:"processed_at"`
}

type Balance struct {
	Current   string `json:"current"`
	Withdrawn string `json:"withdrawn"`
}
