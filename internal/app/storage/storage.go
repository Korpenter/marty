package storage

import (
	"context"
	"fmt"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/models"
	"log"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.Authorization) error
	GetHashedPasswordByLogin(ctx context.Context, login string) (string, error)
	AddOrder(ctx context.Context, order *models.Order) error
	GetOrdersByUser(ctx context.Context, userID string) ([]models.OrderItem, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	Ping(ctx context.Context) error
	DeleteRepo(ctx context.Context) error
}

func New(c *config.Config) (Repository, error) {
	if c.PostgresURI != "" {
		r, err := NewPostgresRepo(c.PostgresURI)
		if err != nil {
			log.Fatal(fmt.Errorf("error initiating postgres connection : %v", err))
		}
		err = r.Ping(context.Background())
		if err != nil {
			log.Fatal(fmt.Errorf("error reaching db : %v", err))
		}
		err = r.NewTables()
		if err != nil {
			log.Fatal(fmt.Errorf("error creating tables : %v", err))
		}
		return r, nil
	}
	return nil, constant.ErrNoStorageSpecified
}
