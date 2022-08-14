package agent

import "github.com/theoremoon/kosenctfx/scoreserver/client"

type Agent interface {
	Client() *Client
}

type agent struct {
	client *Client
}

func New(client *client.Client) (Agent, error) {
	hostname, err := GetHostname()
	if err != nil {
		return nil, err
	}
	hostaddr, err := GetHostAddress()
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client:   client,
		Hostname: hostname,
		Hostaddr: hostaddr,
	}
	return &agent{
		client: c,
	}, nil
}

func (a *agent) Client() *Client {
	return a.client
}
