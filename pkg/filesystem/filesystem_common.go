package filesystem

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type FileSystem struct {
}

type ByteCountWriter struct {
	Writer io.Writer
	Count  int64
}

func (w *ByteCountWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	if err != nil {
		err = errors.Wrap(err, "Write filesystem:")

		return n, err
	}

	w.Count += int64(n) // todo: if count > 2gb то прервать

	return n, nil
}

func (f *FileSystem) CreateFileFromReader(reader io.Reader, fullPath string) error {
	file, err := os.Create(fullPath)
	if err != nil {
		return errors.Wrap(err, "CreateFileFromReader filesystem:")
	}

	defer file.Close() //nolint:errcheck

	buf := make([]byte, 1024*1024) // todo: get buffer from sync pool

	_, err = io.CopyBuffer(&ByteCountWriter{Writer: file}, reader, buf)
	if err != nil {
		return errors.Wrap(err, "CreateFileFromReader filesystem:")
	}

	return nil
}

func (f *FileSystem) DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return errors.Wrap(err, "DeleteFile filesystem:")
	}

	return nil
}

func (f *FileSystem) ParseFileName(fileName string) (name string, ext string, err error) {
	ind := strings.LastIndex(fileName, ".")

	if ind == -1 {
		return "", "", errors.New("cannot find `.` in file_name")
	}

	ext = fileName[ind+1:]
	name = fileName[:ind]

	if ext == "" {
		return "", "", errors.New("empty extension")
	}

	if name == "" {
		return "", "", errors.New("empty name")
	}

	return name, ext, nil
}
