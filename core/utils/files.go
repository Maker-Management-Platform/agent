package utils

import (
	"crypto/sha1"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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

func GetFileSha512(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func GetStringSha512(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ToLibPath(p string) string {
	if strings.HasPrefix(p, runtime.Cfg.Library.Path) {
		return p
	}
	return path.Clean(path.Join(runtime.Cfg.Library.Path, p))
}

func ToAssetsPath(p string, f string) string {
	return path.Clean(path.Join(runtime.GetDataPath(), "assets", p, f))
}

func CreateAssetsFolder(p string) error {
	assetsPath := filepath.Join(runtime.GetDataPath(), "assets")
	return CreateFolder(filepath.Join(assetsPath, p))
}

func DeleteAssetsFolder(p string) error {
	assetsPath := filepath.Join(runtime.GetDataPath(), "assets")
	return os.RemoveAll(filepath.Join(assetsPath, p))
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

func PathToTags(path string) []string {

	path = strings.Trim(path, "/")
	tags := strings.Split(path, "/")
	tagSet := make(map[string]bool)
	for _, t := range tags {
		if t != "" {
			tagSet[t] = true
		}

	}
	rtn := make([]string, len(tagSet))
	i := 0
	for k := range tagSet {
		rtn[i] = k
		i++
	}

	return rtn
}
