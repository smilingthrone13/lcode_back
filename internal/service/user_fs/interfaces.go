package user_fs

import (
	"context"
	"lcode/internal/domain"
)

type UserFS interface {
	MakeUserDir(ctx context.Context, user domain.User) error
	CreateAvatar(ctx context.Context, dto domain.UploadUserAvatarDTO) (origPath, thumbnailPath string, err error)
	DeleteAvatar(ctx context.Context, userID string) error
	AvatarPath(ctx context.Context, userID string) (origPath string, err error)
	AvatarThumbnailPath(ctx context.Context, userID string) (thumbnailPath string, err error)
}

type ThumbnailsService interface {
	CreateThumbnail(ctx context.Context, item domain.CreateThumbnailData) (thumbNailFilePath string, err error)
}
