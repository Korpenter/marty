package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type PostgresRepo struct {
	conn *pgxpool.Pool
	sync.Mutex
}

func NewPostgresRepo(connString string) (*PostgresRepo, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{conn: conn}, nil
}

func (r *PostgresRepo) NewTables() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Exec(ctx, createUsers)
	if err != nil {
		return fmt.Errorf("users : %s", err)
	}
	_, err = r.conn.Exec(ctx, createOrders)
	if err != nil {
		return fmt.Errorf("orders : %s", err)
	}
	_, err = r.conn.Exec(ctx, createWithdrawals)
	if err != nil {
		return fmt.Errorf("orders : %s", err)
	}
	return nil
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user *models.Authorization) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	commandTag, err := tx.Exec(ctx, createUser, user.Login, user.Password)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return constant.ErrUserExists
	}
	return nil
}
func (r *PostgresRepo) GetHashedPasswordByLogin(ctx context.Context, login string) (string, error) {
	var hash string
	err := r.conn.QueryRow(ctx, getHashByLogin, login).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", constant.ErrUserNotFound
		}
		return "", err
	}
	return hash, nil
}

func (r *PostgresRepo) AddOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	commandTag, err := tx.Exec(ctx, addOrder, order.OrderID, order.Login)
	if err != nil {
		return err
	}
	var returnLogin string
	if commandTag.RowsAffected() != 1 {
		err = tx.QueryRow(ctx, getOrderUserid, order.OrderID).Scan(&returnLogin)
		if err != nil {
			return err
		}
		if returnLogin != order.Login {
			return constant.ErrOrderAlreadyAdded
		}
		return constant.ErrOrderAlreadyAddedByUser
	}
	return nil
}

func (r *PostgresRepo) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return nil, nil
}
func (r *PostgresRepo) GetOrdersByUser(ctx context.Context, userID string) ([]*models.Order, error) {
	return nil, nil
}

func (r *PostgresRepo) UpdateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	if order.Status == "PROCESSED" {
		_, err = tx.Exec(ctx, updateProcessedOrder, order.Status, order.Accrual, order.OrderID)
	} else {
		_, err = tx.Exec(ctx, updateOrder, order.Status, order.OrderID)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

func (r *PostgresRepo) DeleteRepo(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, dropTables)
	if err != nil {
		return err
	}
	return nil
}
