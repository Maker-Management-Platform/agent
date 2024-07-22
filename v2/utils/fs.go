package utils

import (
	"errors"
	"io"
	"os"

	cp "github.com/otiai10/copy"
)

func CreateFolder(name string) error {
	_, err := os.Stat(name)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.Mkdir(name, os.ModePerm); err != nil {
			if !errors.Is(err, os.ErrExist) {
				return err
			}
		}
	}

	return nil
}
func Move(src, dst string) error {
	if err := cp.Copy(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

func SaveFile(path string, in io.Reader) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
