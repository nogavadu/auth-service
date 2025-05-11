package env

import (
	"fmt"
	"github.com/nogavadu/auth-service/internal/config"
	"net"
	"os"
	"strconv"
)

const (
	httpHostEnv = "GRPC_SERVER_HOST"
	httpPortEnv = "GRPC_SERVER_PORT"
)

type httpServerConfig struct {
	host string
	port int
}

func NewHTTPServerConfig() (config.GRPCServerConfig, error) {
	const op = "config.NewHTTPServerConfig"

	host := os.Getenv(httpHostEnv)
	if host == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, httpHostEnv)
	}

	portStr := os.Getenv(httpPortEnv)
	if portStr == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, httpPortEnv)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: invalid env variable", op, httpPortEnv)
	}

	return &httpServerConfig{
		host: host,
		port: port,
	}, nil
}

func (c *httpServerConfig) Port() int {
	return c.port
}

func (c *httpServerConfig) Address() string {
	return net.JoinHostPort(c.host, strconv.Itoa(c.port))
}
