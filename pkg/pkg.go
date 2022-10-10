package pkg

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	USER UserType = iota + 1
	KITCHEN
)

var (
	ErrNoRowsAffected = errors.New("no rows affected")
)

type UserType int

type AppError struct {
	MSG  string
	Err  error
	Code int
}

func (ae AppError) Error() string {
	if ae.Err == nil {
		return ae.MSG
	}
	return fmt.Sprintf("%v: %v", ae.MSG, ae.Err)
}

func (ae AppError) Unwrap() error { return ae.Err }

func NewAppError(msg string, err error, code int) error {
	return &AppError{MSG: msg, Err: err, Code: code}
}

// FindError give an error try to find an http status code and message inside it
// If an *AppError is not wraped so returns internal server status code and an empty message
func FindError(err error) (int, string) {
	if err == nil {
		return 0, ""
	}
	temp := err
	for temp != nil {
		if ae, ok := temp.(*AppError); ok {
			return ae.Code, ae.MSG
		}
		temp = errors.Unwrap(temp)
	}
	code := http.StatusInternalServerError
	return code, ""
}
