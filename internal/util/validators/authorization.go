package validators

import (
	"fmt"
	"github.com/Mldlr/mart/marty/internal/app/models"
)

func ValidateAuthorization(user *models.Authorization) error {
	if user == nil {
		return fmt.Errorf("no data provided")
	}
	if user.Login == "" {
		return fmt.Errorf("empty login")
	}
	if user.Password == "" {
		return fmt.Errorf("empty password")
	}
	return nil
}
