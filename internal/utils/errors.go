package utils

type ErrUnauthorized struct {
}

func (e ErrUnauthorized) Error() string {
	return "Unauthorized"
}

// Constructor function for ErrUnauthorized
func NewErrUnauthorized() ErrUnauthorized {
	return ErrUnauthorized{}
}
