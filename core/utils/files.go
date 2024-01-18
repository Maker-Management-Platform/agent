package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/core/runtime"
	cp "github.com/otiai10/copy"
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
	if strings.HasPrefix(path, runtime.Cfg.LibraryPath) {
		return path
	}
	return fmt.Sprintf("%s/%s", runtime.Cfg.LibraryPath, path)
}

func Move(src, dst string, toLibPath bool) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if toLibPath {
		src = ToLibPath(src)
		dst = ToLibPath(dst)
	}
	log.Print(dst)
	if err := cp.Copy(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

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
