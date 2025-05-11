package access

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/auth-service/internal/service"
	accessDesc "github.com/nogavadu/auth-service/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

const authPrefix = "Bearer "

type Implementation struct {
	accessDesc.UnimplementedAccessV1Server
	serv service.AccessService
}

func New(accessService service.AccessService) *Implementation {
	return &Implementation{
		serv: accessService,
	}
}

func (i *Implementation) Check(ctx context.Context, req *accessDesc.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	requiredLvl := req.GetRequiredLvl()
	if err := validator.New().Var(requiredLvl, "required"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := i.serv.Check(ctx, accessToken, int(requiredLvl)); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &emptypb.Empty{}, nil
}
