package ffmpeg_converter

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
)

func New() *Converter {
	return &Converter{}
}

type Converter struct {
}

type MediaInfo struct {
	FrameRate   float64 `json:"fps"`
	FramesCount int     `json:"nframes"`
	Duration    float64 `json:"duration"`
}

func (c *Converter) CreateThumbnailFromImage(inputFileName, outputFileName string, size [2]int) (err error) {
	//nolint:lll
	cmd := exec.Command(
		"ffmpeg", "-y", "-i", inputFileName,
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease:force_divisible_by=2,pad=%d:%d:-1:-1:color=black", size[0], size[1], size[0], size[1]),
		"-c:v", "libwebp",
		outputFileName,
	)

	if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "ThumbnailFromImage converter service")
	}

	if err = cmd.Wait(); err != nil {
		return errors.Wrap(err, "ThumbnailFromImage converter service")
	}

	return
}

func (c *Converter) CreateThumbnailFromVideo(inputFileName, outputFileName string, size [2]int) (err error) {
	//nolint:lll
	cmd := exec.Command(
		"ffmpeg", "-y", "-i", inputFileName,
		"-vf", fmt.Sprintf("thumbnail,scale=%d:%d:force_original_aspect_ratio=decrease:force_divisible_by=2,pad=%d:%d:-1:-1:color=black", size[0], size[1], size[0], size[1]),
		"-frames:v", "1",
		"-c:v", "libwebp",
		outputFileName,
	)

	if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "ThumbnailFromVideo converter service")
	}

	if err = cmd.Wait(); err != nil {
		return errors.Wrap(err, "ThumbnailFromVideo converter service")
	}

	return
}

func (c *Converter) DownscaleVideoToFullHD(inputFileName, outputFileName string) (err error) {
	cmd := exec.Command(
		"ffmpeg", "-y", "-i", inputFileName,
		"-vf", "scale='min(1920, iw)':'min(1080, ih)':force_original_aspect_ratio=decrease",
		"-c:v", "libx264", "-crf", "30", "-preset", "slower", "-movflags", "faststart", "-threads", "10",
		outputFileName,
	)

	if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "DownscaleVideoToFullHD converter service")
	}

	if err = cmd.Wait(); err != nil {
		return errors.Wrap(err, "DownscaleVideoToFullHD converter service")
	}

	return
}

func parseFraction(fractionStr string) (result float64, err error) {
	parts := strings.Split(fractionStr, "/")
	if len(parts) != 2 { //nolint:gomnd
		return 0.0, errors.Wrap(
			errors.New("invalid fraction format"),
			"ParseFraction converter service",
		)
	}

	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0.0, errors.Wrap(err, "ParseFraction converter service")
	}

	denominator, err := strconv.Atoi(parts[1])
	if err != nil || denominator == 0 {
		return 0.0, errors.Wrap(err, "ParseFraction converter service")
	}

	result = float64(numerator) / float64(denominator)

	return
}

func (c *Converter) GetMediaInfo(inputFileName string) (mediaInfo MediaInfo, err error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-show_streams",
		"-select_streams", "v:0",
		"-of", "json",
		inputFileName,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return MediaInfo{}, errors.Wrap(err, "GetMediaInfo converter service")
	}

	var ffprobeOutput struct {
		Streams []struct {
			FrameRate   string `json:"r_frame_rate"`
			FramesCount string `json:"nb_frames"`
			Duration    string `json:"duration"`
		} `json:"streams"`
	}

	if err = json.Unmarshal(output, &ffprobeOutput); err != nil {
		return MediaInfo{}, err
	}

	if len(ffprobeOutput.Streams) == 0 {
		return MediaInfo{},
			errors.Wrap(
				errors.New("no video stream found in the input file"),
				"GetMediaInfo converter service",
			)
	}

	mediaInfo.FramesCount, err = strconv.Atoi(ffprobeOutput.Streams[0].FramesCount)
	if err != nil {
		return MediaInfo{}, errors.Wrap(err, "GetMediaInfo converter service")
	}

	mediaInfo.Duration, err = strconv.ParseFloat(ffprobeOutput.Streams[0].Duration, 64)
	if err != nil {
		return MediaInfo{}, errors.Wrap(err, "GetMediaInfo converter service")
	}

	mediaInfo.FrameRate, err = parseFraction(ffprobeOutput.Streams[0].FrameRate)
	if err != nil {
		return MediaInfo{}, errors.Wrap(err, "GetMediaInfo converter service")
	}

	return mediaInfo, nil
}
