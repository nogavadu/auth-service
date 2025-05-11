package env

import (
	"fmt"
	"github.com/nogavadu/auth-service/internal/config"
	"net"
	"os"
	"strconv"
)

const (
	grpcHostEnv = "GRPC_SERVER_HOST"
	grpcPortEnv = "GRPC_SERVER_PORT"
)

type grpcServerConfig struct {
	host string
	port int
}

func NewGRPCServerConfig() (config.GRPCServerConfig, error) {
	const op = "config.NewGRPCServerConfig"

	host := os.Getenv(grpcHostEnv)
	if host == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, grpcHostEnv)
	}

	portStr := os.Getenv(grpcPortEnv)
	if portStr == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, grpcPortEnv)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: invalid env variable", op, grpcPortEnv)
	}

	return &grpcServerConfig{
		host: host,
		port: port,
	}, nil
}

func (c *grpcServerConfig) Port() int {
	return c.port
}

func (c *grpcServerConfig) Address() string {
	return net.JoinHostPort(c.host, strconv.Itoa(c.port))
}
