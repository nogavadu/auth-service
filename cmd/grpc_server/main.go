package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	accessAPI "github.com/nogavadu/auth-service/internal/api/grpc/access"
	authAPI "github.com/nogavadu/auth-service/internal/api/grpc/auth"
	envConfig "github.com/nogavadu/auth-service/internal/config/env"
	roleRepo "github.com/nogavadu/auth-service/internal/repository/role"
	userRepo "github.com/nogavadu/auth-service/internal/repository/user"
	accessService "github.com/nogavadu/auth-service/internal/service/access"
	authService "github.com/nogavadu/auth-service/internal/service/auth"
	descAccess "github.com/nogavadu/auth-service/pkg/access_v1"
	descAuth "github.com/nogavadu/auth-service/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"strconv"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	jwtConfig, err := envConfig.NewJWTConfig()
	if err != nil {
		log.Error("failed to load JWT config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pgConfig, err := envConfig.NewPGConfig()
	if err != nil {
		log.Error("failed to load PG config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	grpcServerConfig, err := envConfig.NewGRPCServerConfig()
	if err != nil {
		log.Error("failed to load GRPC server config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Error("failed to connect to DB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Error("failed to ping DB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", grpcServerConfig.Address())
	if err != nil {
		os.Exit(1)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	descAuth.RegisterAuthV1Server(
		s, authAPI.New(
			authService.New(
				log,
				jwtConfig.RefreshTokenSecret(),
				jwtConfig.RefreshTokenExp(),
				jwtConfig.AccessTokenSecret(),
				jwtConfig.AccessTokenExp(),
				userRepo.New(db),
			),
		),
	)
	descAccess.RegisterAccessV1Server(
		s, accessAPI.New(
			accessService.New(
				log,
				jwtConfig.RefreshTokenSecret(),
				jwtConfig.RefreshTokenExp(),
				jwtConfig.AccessTokenSecret(),
				jwtConfig.AccessTokenExp(),
				userRepo.New(db),
				roleRepo.New(db),
			),
		),
	)

	log.Info("Starting gRPC Server", slog.String("port", strconv.Itoa(grpcServerConfig.Port())))
	if err = s.Serve(lis); err != nil {
		os.Exit(1)
	}
}
