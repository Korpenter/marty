package handlers

import (
	"errors"
	"fmt"
	"github.com/Mldlr/marty/internal/app"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/container"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/service"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			logging.Logger.Error("registration error :" + err.Error())
		}
	}()
	cred := r.Context().Value(constant.CredKey).(*models.Authorization)
	userService := container.Container.Get("userService").(service.UserService)
	err = userService.CreateUser(r.Context(), cred)
	switch {
	case errors.Is(app.ErrUserExists, err):
		http.Error(w, fmt.Sprintf("registration error: %s", err), http.StatusConflict)
		return
	case errors.Is(app.ErrDataValidation, err):
		http.Error(w, fmt.Sprintf("registration error: %s", err), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, fmt.Sprintf("registration error: %s", err), http.StatusInternalServerError)
		return
	}
	jwtCookie, err := userService.BakeJWTCookie(cred.Login)
	if err != nil {
		http.Error(w, fmt.Sprintf("registration error: %s", err), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, jwtCookie)
	w.WriteHeader(http.StatusOK)
}
