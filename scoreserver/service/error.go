package service

type ErrorMessage struct {
	string
}

func NewErrorMessage(msg string) ErrorMessage {
	return ErrorMessage{
		msg,
	}
}

func (msg ErrorMessage) Error() string {
	return msg.string
}
