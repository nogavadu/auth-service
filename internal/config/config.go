package config

import "time"

type JWTConfig interface {
	RefreshTokenSecret() string
	RefreshTokenExp() time.Duration
	AccessTokenSecret() string
	AccessTokenExp() time.Duration
}

type PGConfig interface {
	DSN() string
}

type GRPCServerConfig interface {
	Port() int
	Address() string
}

type HTTPServerConfig interface {
	Port() int
	Address() string
}
