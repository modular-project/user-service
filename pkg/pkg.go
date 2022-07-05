package pkg

import "errors"

var (
	ErrNoRowsAffected error = errors.New("no rows affected")
	ErrNullValue            = errors.New("null value")
)

const (
	USER UserType = iota + 1
	KITCHEN
)

type UserType int
