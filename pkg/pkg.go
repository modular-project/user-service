package pkg

// var (
// 	ErrNoRowsAffected error = errors.New("no rows affected")
// 	ErrNullValue            = errors.New("null value")
// )

const (
	USER UserType = iota + 1
	KITCHEN
)

type UserType int

type UnauthorizedErr string

type BadErr string

type ForbiddenErr string

type NotFoundErr string

func (f ForbiddenErr) Error() string {
	return string(f)
}

func (f ForbiddenErr) IsForbidden() {}

func (u UnauthorizedErr) Error() string {
	return string(u)
}

func (u UnauthorizedErr) IsUnauthorized() {}

func (b BadErr) Error() string {
	return string(b)
}

func (b BadErr) IsBad() {}

func (n NotFoundErr) Error() string {
	return string(n)
}

func (n NotFoundErr) IsNotFound() {}
