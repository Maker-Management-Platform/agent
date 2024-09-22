package renderers

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type gCodeRenderer struct {
}
type tmpImg struct {
	height int
	width  int
	data   []byte
}

func (r *gCodeRenderer) Render(ctx context.Context, asset *entities.Asset) (*entities.Asset, error) {
	genFS, err := libfs.GetLibFS("generated")
	if err != nil {
		return nil, fmt.Errorf("render error getting fs: %w", err)
	}
	imgName := fmt.Sprintf("%s.r.png", asset.ID)

	if _, err := fs.Stat(genFS, imgName); err == nil {
		return entities.NewAsset(genFS.(entities.LibFS), imgName, false, asset), nil
	}

	slog.Info("Rendering", "asset", *asset.Path, "img", imgName, "asset", asset)

	objFs, err := libfs.GetAssetFS(ctx, *asset)
	if err != nil {
		return nil, fmt.Errorf("error getting fs: %w", err)
	}

	f, err := objFs.Open(utils.VoZ(asset.Path))
	if err != nil {
		return nil, err
	}
	img := &tmpImg{}

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
				if i.width > img.width || i.height > img.height {
					img = i
				}

			}

		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Join(err, errors.New("error reading gcode"))
	}

	if img.data == nil {
		return nil, errors.New("no thumbnail found")
	}
	h := sha1.New()
	_, err = h.Write(img.data)
	if err != nil {
		return nil, err
	}

	writer, err := genFS.Create(imgName)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	i, _, err := image.Decode(bytes.NewReader(img.data))
	if err != nil {
		return nil, err
	}
	if err := png.Encode(writer, i); err != nil {
		return nil, err
	}

	return entities.NewAsset(genFS.(entities.LibFS), imgName, false, asset), nil
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
