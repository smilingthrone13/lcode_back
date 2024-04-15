package thumbnails

import (
	"context"
	"lcode/internal/domain"
)

type Thumbnails interface {
	CreateThumbnail(
		ctx context.Context,
		item domain.CreateThumbnailData,
	) (thumbNailFilePath string, err error)
}
