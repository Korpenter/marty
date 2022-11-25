package storage

import (
	"context"
	"fmt"
	"github.com/Mldlr/marty/internal/app"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/logging"
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

func New(cfg *config.Config) (Repository, error) {
	if cfg.PostgresURI != "" {
		r, err := NewPostgresRepo(cfg.PostgresURI)
		if err != nil {
			logging.Logger.Fatal(fmt.Sprintf("Error initiating postgres connection: %v", err))
		}
		err = r.Ping(context.Background())
		if err != nil {
			logging.Logger.Fatal(fmt.Sprintf("Error reaching db: %v", err))
		}
		err = r.NewTables()
		if err != nil {
			logging.Logger.Fatal(fmt.Sprintf("Error creating tables: %v", err))
		}
		return r, nil
	}
	return nil, app.ErrNoStorageSpecified
}
