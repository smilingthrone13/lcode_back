//go:build !windows

package filesystem

func (f *FileSystem) SetHidden(path string) error { // only for linux
	return nil
}
