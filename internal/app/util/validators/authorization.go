package validators

import (
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/pkg/errors"
)

func ValidateAuthorization(user *models.Authorization) error {
	if user.Login == "" {
		return errors.Wrap(models.ErrDataValidation, "empty login")
	}
	if user.Password == "" {
		return errors.Wrap(models.ErrDataValidation, "empty password")
	}
	return nil
}
