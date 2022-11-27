package storage

import (
	"context"
	models2 "github.com/Mldlr/marty/internal/app/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models2.Authorization) error
	GetHashedPasswordByLogin(ctx context.Context, login string) (string, error)
	AddOrder(ctx context.Context, order *models2.Order) error
	GetOrdersByUser(ctx context.Context) ([]models2.OrderItem, error)
	UpdateOrder(ctx context.Context, order *models2.Order) error
	GetBalance(ctx context.Context) (*models2.Balance, error)
	GetWithdrawals(ctx context.Context) ([]models2.Withdrawal, error)
	Withdraw(ctx context.Context, withdrawal *models2.Withdrawal) error
	Ping(ctx context.Context) error
	DeleteRepo(ctx context.Context) error
}
