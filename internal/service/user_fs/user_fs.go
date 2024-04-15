package user_fs

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/filesystem"
	"lcode/pkg/struct_errors"
	"log/slog"
	"os"
	"path"
	"strings"
)

var (
	avatarThumbnailSize = [2]int{360, 360} // todo:  config value
)

type Services struct {
	Thumbnails ThumbnailsService
}

type Service struct {
	cfg        *config.Config
	logger     *slog.Logger
	services   *Services
	fileSystem *filesystem.FileSystem
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	services *Services,
) *Service {
	return &Service{
		cfg:        cfg,
		logger:     logger,
		services:   services,
		fileSystem: &filesystem.FileSystem{},
	}
}

func (s *Service) MakeUserDir(ctx context.Context, user domain.User) error {
	mainDir := path.Join(s.cfg.Files.MainFolder, domain.UsersFolder, user.ID)
	avatarDir := path.Join(mainDir, domain.AvatarFolder)
	avatarThumbnailDir := path.Join(avatarDir, domain.ThumbnailFolder)

	err := os.Mkdir(mainDir, os.ModeDir)
	if err != nil {
		return errors.Wrap(err, "MakeUserDir user_fs service")
	}

	err = os.Mkdir(avatarDir, os.ModeDir)
	if err != nil {
		return errors.Wrap(err, "MakeUserDir user_fs service")
	}

	err = os.Mkdir(avatarThumbnailDir, os.ModeDir)
	if err != nil {
		return errors.Wrap(err, "MakeUserDir user_fs service")
	}

	return nil
}

func (s *Service) AvatarPath(ctx context.Context, userID string) (origPath string, err error) {
	avatarDir, _, err := s.getAvatarAndThumbnailDirs(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "AvatarPath user_fs service")
	}

	dirEntries, err := os.ReadDir(avatarDir)
	for i := range dirEntries {
		if strings.Contains(dirEntries[i].Name(), domain.DefaultOriginalName) {
			origPath = path.Join(avatarDir, dirEntries[i].Name())
			return origPath, nil
		}
	}

	return "", struct_errors.NewErrNotFound("avatar not found", nil)
}

func (s *Service) AvatarThumbnailPath(ctx context.Context, userID string) (thumbnailPath string, err error) {
	_, thumbnailDir, err := s.getAvatarAndThumbnailDirs(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "AvatarThumbnailPath user_fs service")
	}

	thumbnailPath = path.Join(
		thumbnailDir,
		fmt.Sprintf("%s.%s", domain.DefaultPreviewName, domain.DefaultPreviewExtension),
	)

	return thumbnailPath, nil
}

func (s *Service) getAvatarAndThumbnailDirs(
	ctx context.Context,
	userID string,
) (avatarDir string, avatarThumbnailDir string, err error) {
	mainDir := path.Join(s.cfg.Files.MainFolder, domain.UsersFolder, userID)
	avatarDir = path.Join(mainDir, domain.AvatarFolder)
	avatarThumbnailDir = path.Join(avatarDir, domain.ThumbnailFolder)

	return avatarDir, avatarThumbnailDir, nil
}

func (s *Service) CreateAvatar(
	ctx context.Context,
	dto domain.UploadUserAvatarDTO,
) (origPath, thumbnailPath string, err error) {
	avatarDir, thumbnailDir, err := s.getAvatarAndThumbnailDirs(ctx, dto.User.ID)
	if err != nil {
		return "", "", errors.Wrap(err, "CreateAvatar user_fs service")
	}

	avatarFullPath := path.Join(avatarDir, fmt.Sprintf("%s.%s", domain.DefaultOriginalName, dto.Extension))

	err = s.fileSystem.CreateFileFromReader(dto.Media, avatarFullPath)
	if err != nil {
		return "", "", errors.Wrap(err, "CreateAvatar user_fs service")
	}

	thumbnailPath, err = s.services.Thumbnails.CreateThumbnail(
		ctx, domain.CreateThumbnailData{
			ThumbnailSize:     avatarThumbnailSize,
			MediaType:         dto.MediaType,
			SrcFilePath:       avatarFullPath,
			DestPath:          thumbnailDir,
			ThumbnailFileName: domain.DefaultPreviewName,
		},
	)
	if err != nil {
		return "", "", errors.Wrap(err, "CreateAvatar user_fs service")
	}

	return avatarFullPath, thumbnailPath, nil
}

func (s *Service) DeleteAvatar(ctx context.Context, userID string) error {
	avatarDir, thumbnailDir, err := s.getAvatarAndThumbnailDirs(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "DeleteAvatar user_fs service")
	}

	dirEntries, err := os.ReadDir(avatarDir)
	if err != nil {
		return errors.Wrap(err, "DeleteAvatar user_fs service")
	}

	for i := range dirEntries {
		if !strings.Contains(dirEntries[i].Name(), domain.DefaultOriginalName) {
			continue
		}

		err = os.Remove(path.Join(avatarDir, dirEntries[i].Name()))
		if err != nil && !os.IsNotExist(err) {
			s.logger.Error("cannot remove avatar", slog.String("err", err.Error()))

			return errors.Wrap(err, "DeleteAvatar user_fs service")
		}
	}

	dirEntries, err = os.ReadDir(thumbnailDir)
	if err != nil {
		return errors.Wrap(err, "DeleteAvatar user_fs service")
	}

	for i := range dirEntries {
		if !strings.Contains(dirEntries[i].Name(), domain.DefaultPreviewName) {
			continue
		}

		err = os.Remove(path.Join(thumbnailDir, dirEntries[i].Name()))
		if err != nil && !os.IsNotExist(err) {
			s.logger.Error("cannot remove avatar thumbnail", slog.String("err", err.Error()))

			return errors.Wrap(err, "DeleteAvatar user_fs service")
		}
	}

	return nil
}
