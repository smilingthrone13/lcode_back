//go:build windows

package filesystem

import (
	"github.com/pkg/errors"
	"syscall"
)

func (f *FileSystem) SetHidden(path string) error { // only for windows
	filenameW, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return errors.Wrap(err, "SetHidden filesystem:")
	}

	err = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return errors.Wrap(err, "SetHidden filesystem:")
	}

	return nil
}
