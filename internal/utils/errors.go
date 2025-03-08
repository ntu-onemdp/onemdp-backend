package utils

type ErrUnauthorized struct {
}

func (e ErrUnauthorized) Error() string {
	return "Unauthorized"
}
