package model

type ErrorMessage string

func (msg ErrorMessage) Error() string {
	return string(msg)
}
