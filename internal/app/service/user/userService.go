package user

import (
	"context"
	"github.com/Mldlr/marty/internal/app/models"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.Authorization) error
	LogInUser(ctx context.Context, user *models.Authorization) error
	GetBalance(ctx context.Context) (*models.Balance, error)
	GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error)
	Withdraw(ctx context.Context, withdrawal *models.Withdrawal) error
	BakeJWTCookie(login string) (*http.Cookie, error)
	hashPassword(password string) string
	checkPasswordHash(pass string, hash string) bool
}
