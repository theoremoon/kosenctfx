package agent

type Agent interface {
}

type agent struct {
}

func New() (Agent, error) {

	return &agent{}, nil
}
