package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
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
	if strings.HasPrefix(path, filepath.Clean(runtime.Cfg.LibraryPath)) {
		return path
	}
	return filepath.Clean(fmt.Sprintf("%s/%s", runtime.Cfg.LibraryPath, path))
}

func Move(src, dst string) error {
	dst = ToLibPath(dst)
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
		if err := os.Mkdir("data", 0666); err != nil {
			if !errors.Is(err, os.ErrExist) {
				return err
			}
		}
	}

	return nil
}
