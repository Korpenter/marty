package service

import (
	"context"
	"github.com/Mldlr/marty/internal/app"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/storage"
	"github.com/Mldlr/marty/internal/util/validators"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
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

type UserServiceImpl struct {
	tokenAuth *jwtauth.JWTAuth
	repo      storage.Repository
	cfg       *config.Config
}

func NewUserService(cfg *config.Config, repo storage.Repository) (UserService, error) {
	return &UserServiceImpl{
		tokenAuth: jwtauth.New("HS256", []byte(cfg.SecretKey), nil),
		repo:      repo,
		cfg:       cfg,
	}, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *models.Authorization) error {
	if err := validators.ValidateAuthorization(user); err != nil {
		return app.ErrDataValidation
	}
	user.Password = s.hashPassword(user.Password)
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) LogInUser(ctx context.Context, user *models.Authorization) error {
	if err := validators.ValidateAuthorization(user); err != nil {
		return app.ErrDataValidation
	}
	passwordHash, err := s.repo.GetHashedPasswordByLogin(ctx, user.Login)
	if err != nil {
		return err
	}
	if !s.checkPasswordHash(user.Password, passwordHash) {
		return app.ErrWrongPassword
	}
	return nil
}

func (s *UserServiceImpl) GetBalance(ctx context.Context) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return balance, nil
}

func (s *UserServiceImpl) GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error) {
	withdrawals, err := s.repo.GetWithdrawals(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return withdrawals, nil
}

func (s *UserServiceImpl) Withdraw(ctx context.Context, withdrawal *models.Withdrawal) error {
	err := s.repo.Withdraw(ctx, withdrawal)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *UserServiceImpl) BakeJWTCookie(login string) (*http.Cookie, error) {
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{"login": login})
	if err != nil {
		logging.Logger.Error("constant making token", zap.String("login", login))
		return nil, err
	}
	jwt := &http.Cookie{
		Path:    "/",
		Name:    "jwt",
		Expires: time.Now().Add(7 * 24 * time.Hour),
		Value:   tokenString,
	}
	return jwt, nil
}

func (s *UserServiceImpl) hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func (s *UserServiceImpl) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
