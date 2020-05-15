package service

type errorMessage string

func (err errorMessage) Error() string {
	return string(err)
}

func IsErrorMessage(err error) bool {
	_, ok := err.(errorMessage)
	return ok
}

func ErrorMessage(m string) error {
	return errorMessage(m)
}
