package error

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordInCorrect    = errors.New("password incorrect")
	ErrUsernameExist        = errors.New("username already exist")
	ErrPasswordDoesNotMatch = errors.New("password does not match")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordInCorrect,
	ErrUsernameExist,
	ErrPasswordDoesNotMatch,
}
