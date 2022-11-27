package models

import "errors"

var (
	ErrUserExists              = errors.New("login already exists")
	ErrWrongPassword           = errors.New("wrong password")
	ErrNoStorageSpecified      = errors.New("no storage specified")
	ErrUserNotFound            = errors.New("login not found")
	ErrOrderAlreadyAdded       = errors.New("order already added by another user")
	ErrOrderAlreadyAddedByUser = errors.New("order already added by user")
	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrDataValidation          = errors.New("validation constant")
	ErrAcrrualServerError      = errors.New("accrual server error")
	ErrTooManyRequests         = errors.New("too many requests")
	ErrNoContent               = errors.New("no content")
	ErrDecodingJSON            = errors.New("error decoding json")
)
