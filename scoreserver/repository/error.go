package repository

/// Repositoryに値がなかったときのエラー
type NotFoundError struct {
	string
}

func NotFound(msg string) NotFoundError {
	return NotFoundError{
		msg,
	}
}

func (msg NotFoundError) Error() string {
	return msg.string
}

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
