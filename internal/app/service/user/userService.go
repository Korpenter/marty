package user

import (
	"context"
	models2 "github.com/Mldlr/marty/internal/app/models"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models2.Authorization) error
	LogInUser(ctx context.Context, user *models2.Authorization) error
	GetBalance(ctx context.Context) (*models2.Balance, error)
	GetWithdrawals(ctx context.Context) ([]models2.Withdrawal, error)
	Withdraw(ctx context.Context, withdrawal *models2.Withdrawal) error
	BakeJWTCookie(login string) (*http.Cookie, error)
	hashPassword(password string) string
	checkPasswordHash(pass string, hash string) bool
}
