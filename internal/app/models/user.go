package models

import "github.com/shopspring/decimal"

type LoginKey struct{}
type CredKey struct{}

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Balance struct {
	Current   *decimal.Decimal `json:"current"`
	Withdrawn *decimal.Decimal `json:"withdrawn"`
}
