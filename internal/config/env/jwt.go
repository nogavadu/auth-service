package env

import (
	"fmt"
	"github.com/nogavadu/auth-service/internal/config"
	"os"
	"time"
)

const (
	refreshTokenSecretEnv = "REFRESH_TOKEN_SECRET"
	refreshTokenExpEnv    = "REFRESH_TOKEN_EXP"
	accessTokenSecretEnv  = "ACCESS_TOKEN_SECRET"
	accessTokenExpEnv     = "ACCESS_TOKEN_EXP"
)

type jwtConfig struct {
	refreshTokenSecret string
	refreshTokenExp    time.Duration
	accessTokenSecret  string
	accessTokenExp     time.Duration
}

func NewJWTConfig() (config.JWTConfig, error) {
	const op = "config.NewJWTConfig"

	refreshTokenSecret := os.Getenv(refreshTokenSecretEnv)
	if refreshTokenSecret == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, refreshTokenSecretEnv)
	}

	refreshTokenExp := os.Getenv(refreshTokenExpEnv)
	if refreshTokenExp == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, refreshTokenExpEnv)
	}
	refreshTokenExpTime, err := time.ParseDuration(refreshTokenExp)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, refreshTokenExpEnv, err)
	}

	accessTokenSecret := os.Getenv(accessTokenSecretEnv)
	if accessTokenSecret == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, accessTokenSecretEnv)
	}

	accessTokenExp := os.Getenv(accessTokenExpEnv)
	if accessTokenExp == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, accessTokenExpEnv)
	}
	accessTokenExpTime, err := time.ParseDuration(refreshTokenExp)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, accessTokenExpEnv, err)
	}

	return &jwtConfig{
		refreshTokenSecret: refreshTokenSecret,
		refreshTokenExp:    refreshTokenExpTime,
		accessTokenSecret:  accessTokenSecret,
		accessTokenExp:     accessTokenExpTime,
	}, nil
}

func (j *jwtConfig) RefreshTokenSecret() string {
	return j.refreshTokenSecret
}

func (j *jwtConfig) AccessTokenSecret() string {
	return j.accessTokenSecret
}

func (j *jwtConfig) RefreshTokenExp() time.Duration {
	return j.refreshTokenExp
}

func (j *jwtConfig) AccessTokenExp() time.Duration {
	return j.accessTokenExp
}
