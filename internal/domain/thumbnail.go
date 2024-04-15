package domain

const (
	PictureMedia = "picture"
	VideoMedia   = "video"
)

const (
	DefaultPreviewName      = "preview"
	DefaultPreviewExtension = "webp"
	DefaultOriginalName     = "original"
)

type CreateThumbnailData struct {
	ThumbnailSize     [2]int
	MediaType         string
	SrcFilePath       string
	DestPath          string
	ThumbnailFileName string
}
