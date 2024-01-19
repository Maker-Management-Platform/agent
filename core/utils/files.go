package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/eduardooliveira/stLib/core/runtime"
)

func GetFileSha1(path string) (string, error) {
	f, err := os.Open(ToLibPath(path))
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func ToLibPath(path string) string {
	if strings.HasPrefix(path, runtime.Cfg.Library.Path) {
		return path
	}
	return fmt.Sprintf("%s/%s", runtime.Cfg.Library.Path, path)
}

func Move(src, dst string) error {
	log.Print(path.Dir(dst))
	if err := os.MkdirAll(path.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	return os.Rename(ToLibPath(src), dst)
}

func CreateFolder(name string) error {
	_, err := os.Stat(name)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.Mkdir("data", os.ModePerm); err != nil {
			if !errors.Is(err, os.ErrExist) {
				return err
			}
		}
	}

	return nil
}
