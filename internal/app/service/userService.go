package service

import (
	"context"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/storage"
	"github.com/Mldlr/marty/internal/util/validators"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.Authorization) error
	LogInUser(ctx context.Context, user *models.Authorization) (bool, error)
	MakeToken(login string) string
	hashPassword(password string) string
	checkPasswordHash(pass string, hash string) bool
}

type UserServiceImpl struct {
	tokenAuth *jwtauth.JWTAuth
	repo      storage.Repository
	cfg       *config.Config
}

func NewUserService(c *config.Config, repo storage.Repository) (UserService, error) {
	return &UserServiceImpl{
		tokenAuth: jwtauth.New("HS256", []byte(c.SecretKey), nil),
		repo:      repo,
		cfg:       c,
	}, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *models.Authorization) error {
	if err := validators.ValidateAuthorization(user); err != nil {
		return constant.ErrDataValidation
	}
	user.Password = s.hashPassword(user.Password)
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) LogInUser(ctx context.Context, user *models.Authorization) (bool, error) {
	if err := validators.ValidateAuthorization(user); err != nil {
		return false, constant.ErrDataValidation
	}
	passwordHash, err := s.repo.GetHashedPasswordByLogin(ctx, user.Login)
	if err != nil {
		return false, err
	}
	if !s.checkPasswordHash(user.Password, passwordHash) {
		return false, constant.ErrWrongPassword
	}
	return true, nil
}

func (s *UserServiceImpl) MakeToken(login string) string {
	_, tokenString, _ := s.tokenAuth.Encode(map[string]interface{}{"login": login})
	return tokenString
}

func (s *UserServiceImpl) hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func (s *UserServiceImpl) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
