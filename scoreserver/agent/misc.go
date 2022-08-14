package agent

import (
	"fmt"
	"net"
	"os"

	"github.com/pkg/errors"
)

func GetHostname() (string, error) {
	n, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return n, nil
}

func GetHostAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.Wrap(err, "InterfaceAddrs")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no good ip address")
}
