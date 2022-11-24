package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mldlr/mart/marty/internal/app/constant"
	"github.com/Mldlr/mart/marty/internal/app/container"
	"github.com/Mldlr/mart/marty/internal/app/models"
	"github.com/Mldlr/mart/marty/internal/app/service"
	"github.com/Mldlr/mart/marty/internal/util/validators"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var cred *models.Authorization
	var err error
	if err = json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err = validators.ValidateAuthorization(cred); err != nil {
		http.Error(w, "data validation error", http.StatusBadRequest)
		return
	}
	userService := container.Container.Get("userService").(service.UserService)
	_, err = userService.LogInUser(r.Context(), cred)
	switch {
	case errors.Is(constant.ErrWrongPassword, err) || errors.Is(constant.ErrUserNotFound, err):
		http.Error(w, fmt.Sprintf("cant login: %s", err), http.StatusUnauthorized)
		return
	case errors.Is(constant.ErrDataValidation, err):
		http.Error(w, fmt.Sprintf("cant login: %s", err), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, fmt.Sprintf("cant login: %s", err), http.StatusInternalServerError)
		return
	}
	jwt := &http.Cookie{
		Path:    "/",
		Name:    "jwt",
		Expires: time.Now().Add(7 * 24 * time.Hour),
		Value:   userService.MakeToken(cred.Login),
	}
	http.SetCookie(w, jwt)
	w.WriteHeader(http.StatusOK)
}
