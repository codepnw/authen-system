package errs

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("auth: invalid email or password")
	ErrGenerateToken          = errors.New("auth: generate token failed")
	ErrInvalidToken           = errors.New("auth: invalid token")
	ErrSaveToken              = errors.New("auth: save token failed")
)
