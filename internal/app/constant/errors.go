package constant

import "errors"

var ErrUserExists = errors.New("login already exists")
var ErrDataValidation = errors.New("validation error")
var ErrWrongPassword = errors.New("wrong password")
var ErrNoStorageSpecified = errors.New("no storage specified")
var ErrUserNotFound = errors.New("login not found")
var ErrOrderAlreadyAdded = errors.New("order already added by another user")
var ErrOrderAlreadyAddedByUser = errors.New("order already added by user")
