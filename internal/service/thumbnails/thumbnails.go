package thumbnails

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/ffmpeg_converter"
	"log/slog"
	"path"
)

type Service struct {
	cfg       *config.Config
	logger    *slog.Logger
	converter *ffmpeg_converter.Converter
}

func New(cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		cfg:       cfg,
		logger:    logger,
		converter: ffmpeg_converter.New(),
	}
}

func (s *Service) CreateThumbnail(
	ctx context.Context,
	item domain.CreateThumbnailData,
) (thumbNailFilePath string, err error) {
	previewFilePath := path.Join(item.DestPath, fmt.Sprintf("%s.webp", item.ThumbnailFileName))

	switch item.MediaType {
	case domain.VideoMedia:
		err = s.converter.CreateThumbnailFromVideo(item.SrcFilePath, previewFilePath, item.ThumbnailSize)
	case domain.PictureMedia:
		err = s.converter.CreateThumbnailFromImage(item.SrcFilePath, previewFilePath, item.ThumbnailSize)
	}

	if err != nil {
		return "", errors.Wrap(err, "CreateThumbnail thumbnails service")
	}

	thumbNailFilePath = previewFilePath

	return thumbNailFilePath, nil
}
