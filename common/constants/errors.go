package constants

import "errors"

var (
	ErrCantParseBody       = errors.New("cant parse body bad request")
	ErrPasswordMismatch    = errors.New("username or password not found")
	ErrOldPasswordNotMatch = errors.New("old password is incorrect")
	ErrUserDeleted         = "user was deleted at %s"
)
