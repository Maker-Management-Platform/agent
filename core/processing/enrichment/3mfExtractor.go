package enrichment

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

type mfExtractor struct{}

func New3MFExtractor() *mfExtractor {
	return &mfExtractor{}
}

func (me *mfExtractor) Extract(e Enrichable) ([]string, error) {
	rtn := make([]string, 0)
	baseName := fmt.Sprintf("%s.%s.e", e.Asset().ProjectUUID, e.Asset().ID)
	path := utils.ToLibPath(filepath.Join(e.Project().FullPath(), e.Asset().Name))

	archive, err := zip.OpenReader(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer archive.Close()

	for _, f := range archive.File {
		ext := filepath.Ext(f.Name)
		// Only allow image files the platform supports
		if !slices.Contains(state.AssetTypes["image"].Extensions, ext) {
			continue
		}

		// Ignore thumbnail since we should have the original image already
		if strings.Contains(f.Name, ".thumbnails/") {
			continue
		}
		dstName := fmt.Sprintf("%s%s", baseName, ext)
		dstFile, err := os.OpenFile(utils.ToGeneratedPath(dstName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Println(err)
			continue
		}
		defer dstFile.Close()

		fileInArchive, err := f.Open()
		if err != nil {
			log.Println(err)
			continue
		}
		defer fileInArchive.Close()

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			log.Println(err)
			continue
		}
		rtn = append(rtn, dstName)

	}

	return rtn, nil
}
