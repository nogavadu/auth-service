package user

import (
	"context"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"github.com/nogavadu/auth-service/internal/service"
	"github.com/nogavadu/auth-service/internal/utils"
	userDesc "github.com/nogavadu/auth-service/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Implementation struct {
	userDesc.UnimplementedUserV1Server
	serv service.UserService
}

func New(userService service.UserService) *Implementation {
	return &Implementation{
		serv: userService,
	}
}

func (i *Implementation) GetById(ctx context.Context, request *userDesc.GetByIdRequest) (*userDesc.GetByIdResponse, error) {
	userId := request.GetId()
	user, err := i.serv.GetById(ctx, int(userId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &userDesc.GetByIdResponse{
		User: &userDesc.User{
			Int: int64(user.Id),
			Info: &userDesc.UserInfo{
				Name:   utils.StringPtrToProtoString(user.Name),
				Email:  user.Email,
				Avatar: utils.StringPtrToProtoString(user.Avatar),
				Role:   int32(user.RoleId),
			},
		},
	}, nil
}

func (i *Implementation) Update(ctx context.Context, request *userDesc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.serv.Update(ctx, int(request.GetId()), &model.UserUpdateInput{
		Name:   utils.ProtoStringToPtrString(request.GetUpdateInput().GetName()),
		Email:  utils.ProtoStringToPtrString(request.GetUpdateInput().GetEmail()),
		Avatar: utils.ProtoStringToPtrString(request.GetUpdateInput().GetAvatar()),
		RoleId: utils.ProtoInt32ToPtrInt(request.GetUpdateInput().GetRole()),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (i *Implementation) Delete(ctx context.Context, request *userDesc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.serv.Delete(ctx, int(request.GetId()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
