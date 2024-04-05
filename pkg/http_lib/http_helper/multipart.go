package http_helper

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"mime"
	"mime/multipart"
)

type MultipartReader struct {
	reader   *multipart.Reader
	lastPart *multipart.Part
}

func NewMultipartReader(contentType string, body io.Reader) (*MultipartReader, error) {
	_, p, err := mime.ParseMediaType(contentType)
	if err != nil {
		err = errors.Wrap(err, "NewMultipartReader http_helper pkg")

		return nil, err
	}

	boundary := p["boundary"]
	reader := multipart.NewReader(body, boundary)

	return &MultipartReader{reader: reader}, nil
}

func (f *MultipartReader) NextPart() (*multipart.Part, error) {
	var err error
	f.lastPart, err = f.reader.NextPart()
	if err != nil {
		err = errors.Wrap(err, "NextPart http_helper pkg")

		return nil, err
	}

	return f.lastPart, nil
}

func (f *MultipartReader) DecodeLast(v any) error {
	decoder := json.NewDecoder(f.lastPart)

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

type FileParams struct {
	ID        string
	FullName  string
	Name      string
	Extension string
}
