package user

import (
	"context"
	"github.com/Mldlr/marty/internal/app/config"
	models2 "github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/storage"
	"github.com/Mldlr/marty/internal/app/util/validators"
	"github.com/go-chi/jwtauth/v5"
	"github.com/samber/do"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserServiceImpl struct {
	tokenAuth *jwtauth.JWTAuth
	repo      storage.Repository
	cfg       *config.Config
	log       *zap.Logger
}

func NewUserService(i *do.Injector) UserService {
	cfg := do.MustInvoke[*config.Config](i)
	repo := do.MustInvoke[storage.Repository](i)
	return &UserServiceImpl{
		tokenAuth: jwtauth.New("HS256", []byte(cfg.SecretKey), nil),
		repo:      repo,
		cfg:       cfg,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *models2.Authorization) error {
	if err := validators.ValidateAuthorization(user); err != nil {
		return models2.ErrDataValidation
	}
	user.Password = s.hashPassword(user.Password)
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) LogInUser(ctx context.Context, user *models2.Authorization) error {
	if err := validators.ValidateAuthorization(user); err != nil {
		return models2.ErrDataValidation
	}
	passwordHash, err := s.repo.GetHashedPasswordByLogin(ctx, user.Login)
	if err != nil {
		return err
	}
	if !s.checkPasswordHash(user.Password, passwordHash) {
		return models2.ErrWrongPassword
	}
	return nil
}

func (s *UserServiceImpl) GetBalance(ctx context.Context) (*models2.Balance, error) {
	balance, err := s.repo.GetBalance(ctx)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (s *UserServiceImpl) GetWithdrawals(ctx context.Context) ([]models2.Withdrawal, error) {
	withdrawals, err := s.repo.GetWithdrawals(ctx)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (s *UserServiceImpl) Withdraw(ctx context.Context, withdrawal *models2.Withdrawal) error {
	err := s.repo.Withdraw(ctx, withdrawal)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) BakeJWTCookie(login string) (*http.Cookie, error) {
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{"login": login})
	if err != nil {
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
