package storage

import (
	"context"
	"github.com/Mldlr/marty/internal/app/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.Authorization) error
	GetHashedPasswordByLogin(ctx context.Context, login string) (string, error)
	AddOrder(ctx context.Context, order *models.Order) error
	GetOrdersByUser(ctx context.Context) ([]models.OrderItem, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	GetBalance(ctx context.Context) (*models.Balance, error)
	GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error)
	Withdraw(ctx context.Context, withdrawal *models.Withdrawal) error
	Ping(ctx context.Context) error
	DeleteRepo(ctx context.Context) error
}
