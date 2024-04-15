package user_manager

import (
	"context"
	"lcode/internal/domain"
	"lcode/pkg/simple_auth"
)

type (
	UserManager interface {
		Register(ctx context.Context, dto domain.CreateUserDTO) (user domain.User, err error)
		Login(ctx context.Context, dto domain.LoginDTO) (tokens simple_auth.Tokens, err error)
		UserByID(ctx context.Context, id string) (user domain.User, err error)
		Users(ctx context.Context) ([]domain.User, error)
		UpdateUser(ctx context.Context, dto domain.UpdateUserDTO) (user domain.User, err error)

		UploadUserAvatar(ctx context.Context, dto domain.UploadUserAvatarDTO) (thumbnailPath string, err error)
		DeleteUserAvatar(ctx context.Context, dto domain.DeleteUserAvatarDTO) error
		AvatarPath(ctx context.Context, userID string) (p string, err error)
		AvatarThumbnailPath(ctx context.Context, userID string) (p string, err error)
	}
)
