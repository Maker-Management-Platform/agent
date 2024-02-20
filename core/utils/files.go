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

func ToLibPath(p string) string {
	if strings.HasPrefix(p, runtime.Cfg.Library.Path) {
		return p
	}
	return path.Clean(path.Join(runtime.Cfg.Library.Path, p))
}

func Move(src, dst string, toLibPath bool) error {
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
