package user

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service/user"
	"github.com/Mldlr/marty/internal/app/util/validators"
	"github.com/go-chi/render"
	"github.com/samber/do"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	userService user.UserService
	logger      *zap.Logger
}

func NewUserController(i *do.Injector) *UserController {
	userService := do.MustInvoke[user.UserService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

func (c *UserController) HandleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	c.logger.Error("request error",
		zap.String("controller", "user"),
		zap.String("url", r.URL.String()),
		zap.Error(err),
	)
	http.Error(w, err.Error(), code)
}

func (c *UserController) Balance(w http.ResponseWriter, r *http.Request) {
	balance, err := c.userService.GetBalance(r.Context())
	if err != nil {
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, balance)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred := r.Context().Value(models.CredKey{}).(*models.Authorization)
	err := c.userService.LogInUser(ctx, cred)
	if err != nil {
		switch err {
		case models.ErrWrongPassword:
			c.HandleError(w, r, err, http.StatusUnauthorized)
			return
		case models.ErrUserNotFound:
			c.HandleError(w, r, err, http.StatusUnauthorized)
			return
		case models.ErrDataValidation:
			c.HandleError(w, r, err, http.StatusBadRequest)
			return
		default:
			c.HandleError(w, r, err, http.StatusInternalServerError)
			return
		}
	}
	jwtCookie, err := c.userService.BakeJWTCookie(cred.Login)
	if err != nil {
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, jwtCookie)
	w.WriteHeader(http.StatusOK)
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred := r.Context().Value(models.CredKey{}).(*models.Authorization)
	err := c.userService.CreateUser(ctx, cred)
	switch err {
	case nil:
		break
	case models.ErrUserExists:
		c.HandleError(w, r, err, http.StatusConflict)
		return
	case models.ErrDataValidation:
		c.HandleError(w, r, err, http.StatusBadRequest)
		return
	default:
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	jwtCookie, err := c.userService.BakeJWTCookie(cred.Login)
	if err != nil {
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, jwtCookie)
	w.WriteHeader(http.StatusOK)
}

func (c *UserController) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	withdrawalsItems, err := c.userService.GetWithdrawals(r.Context())
	if err != nil {
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	if len(withdrawalsItems) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, withdrawalsItems)
}

func (c *UserController) Withdraw(w http.ResponseWriter, r *http.Request) {
	var withdrawal *models.Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		c.HandleError(w, r, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	withdrawal.Login = r.Context().Value(models.LoginKey{}).(string)
	if !validators.Luhn(withdrawal.OrderID) {
		c.HandleError(w, r, fmt.Errorf("invalid order ID"), http.StatusUnprocessableEntity)
		return
	}
	err := c.userService.Withdraw(r.Context(), withdrawal)
	switch err {
	case nil:
		break
	case models.ErrInsufficientBalance:
		c.HandleError(w, r, err, http.StatusBadRequest)
		return
	default:
		c.HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
