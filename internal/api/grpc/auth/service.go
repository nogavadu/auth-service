package auth

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"github.com/nogavadu/auth-service/internal/service"
	authService "github.com/nogavadu/auth-service/internal/service/auth"
	"github.com/nogavadu/auth-service/internal/utils"
	authDesc "github.com/nogavadu/auth-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Implementation struct {
	authDesc.UnimplementedAuthV1Server
	serv service.AuthService
}

func New(authService service.AuthService) *Implementation {
	return &Implementation{
		serv: authService,
	}
}

func (i *Implementation) Register(ctx context.Context, req *authDesc.RegisterRequest) (*authDesc.RegisterResponse, error) {
	email := req.GetEmail()
	if err := validator.New().Var(email, "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	password := req.GetPassword()
	if err := validator.New().Var(password, "required"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	name := req.GetName()

	userId, err := i.serv.Register(ctx, &model.UserInfo{
		Name:  utils.ProtoStringToPtrString(name),
		Email: email,
	}, password)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, authService.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authDesc.RegisterResponse{
		UserId: uint64(userId),
	}, nil
}

func (i *Implementation) Login(ctx context.Context, req *authDesc.LoginRequest) (*authDesc.LoginResponse, error) {
	email := req.GetEmail()
	if err := validator.New().Var(email, "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	password := req.GetPassword()
	if err := validator.New().Var(password, "required"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	token, err := i.serv.Login(ctx, email, password)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authDesc.LoginResponse{
		RefreshToken: token,
	}, nil
}

func (i *Implementation) GetRefreshToken(ctx context.Context, req *authDesc.GetRefreshTokenRequest) (*authDesc.GetRefreshTokenResponse, error) {
	refreshToken := req.GetRefreshToken()
	if err := validator.New().Var(refreshToken, "required,jwt"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	newRefreshToken, err := i.serv.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.Aborted, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authDesc.GetRefreshTokenResponse{
		RefreshToken: newRefreshToken,
	}, nil
}

func (i *Implementation) GetAccessToken(ctx context.Context, req *authDesc.GetAccessTokenRequest) (*authDesc.GetAccessTokenResponse, error) {
	refreshToken := req.GetRefreshToken()
	if err := validator.New().Var(refreshToken, "required,jwt"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	accessToken, err := i.serv.GetAccessToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.Aborted, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authDesc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (i *Implementation) IsUser(ctx context.Context, req *authDesc.IsUserRequest) (*empty.Empty, error) {
	refreshToken := req.GetRefreshToken()
	if err := validator.New().Var(refreshToken, "required,jwt"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	userId := req.GetUserId()
	if err := validator.New().Var(userId, "required"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if err := i.serv.IsUser(ctx, int(userId), refreshToken); err != nil {
		return nil, status.Error(codes.PermissionDenied, "not a same user")
	}

	return nil, nil
}
