package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Withdrawal struct {
	Login     string
	OrderID   string           `json:"order"`
	Sum       *decimal.Decimal `json:"sum"`
	Processed time.Time        `json:"processed_at"`
}
