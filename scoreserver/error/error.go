package error

/// UniqueKeyConstraintに違反したときのエラー
type DuplicatedError struct {
	string
}

func Duplicated(msg string) DuplicatedError {
	return DuplicatedError{
		msg,
	}
}

func (msg DuplicatedError) Error() string {
	return msg.string
}
