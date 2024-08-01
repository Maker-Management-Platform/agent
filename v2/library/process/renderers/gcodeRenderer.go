package renderers

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type gCodeRenderer struct {
}
type tmpImg struct {
	height int
	width  int
	data   []byte
}

func (r *gCodeRenderer) Render(asset *entities.Asset) (*entities.Asset, error) {
	imgName := fmt.Sprintf("%s.r.png", asset.ID)
	imgPath := filepath.Join(*asset.ParentID, imgName)
	imgRoot := filepath.Join(config.Cfg.Core.DataFolder, "img")
	parentDir := filepath.Join(imgRoot, *asset.ParentID)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if err := os.Mkdir(parentDir, 0755); err != nil {
			return nil, err
		}
	}
	fullPath := filepath.Join(parentDir, imgName)
	if _, err := os.Stat(fullPath); err == nil {
		return entities.NewAssetFromRootPath(imgRoot, imgPath, false, asset), nil
	}

	var fPath string
	if utils.VoZ(asset.NodeKind) == entities.NodeKindBundled {
		fPath = filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID), path.Base(*asset.Path))
	} else {
		fPath = filepath.Join(*asset.Root, *asset.Path)
	}
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	image := &tmpImg{}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), ";") {
			line := strings.Trim(scanner.Text(), " ;")

			if strings.HasPrefix(line, "thumbnail begin") {

				header := strings.Split(line, " ")
				length, err := strconv.Atoi(header[3])
				if err != nil {
					return nil, err
				}
				i, err := r.parseThumbnail(scanner, header[2], length)
				if err != nil {
					return nil, err
				}
				if i.width > image.width || i.height > image.height {
					image = i
				}

			}

		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Join(err, errors.New("error reading gcode"))
	}

	if image.data == nil {
		return nil, errors.New("no thumbnail found")
	}
	h := sha1.New()
	_, err = h.Write(image.data)
	if err != nil {
		return nil, err
	}

	f, err = r.storeImage(image, fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return entities.NewAssetFromRootPath(imgRoot, imgPath, false, asset), nil
}

func (r *gCodeRenderer) parseThumbnail(scanner *bufio.Scanner, size string, length int) (*tmpImg, error) {
	sb := strings.Builder{}
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ;")
		if strings.HasPrefix(line, "thumbnail end") {
			break
		}
		sb.WriteString(line)

	}
	if sb.Len() != length {
		return nil, errors.New("thumbnail length mismatch")
	}

	b, err := base64.StdEncoding.DecodeString(sb.String())
	if err != nil {
		return nil, err
	}

	dimensions := strings.Split(size, "x")

	img := &tmpImg{
		data: b,
	}
	img.height, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, err
	}

	img.width, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (r *gCodeRenderer) storeImage(img *tmpImg, path string) (*os.File, error) {
	i, _, err := image.Decode(bytes.NewReader(img.data))
	if err != nil {
		return nil, err
	}
	out, _ := os.Create(path)

	err = png.Encode(out, i)

	if err != nil {
		return nil, err
	}
	return out, nil
}
