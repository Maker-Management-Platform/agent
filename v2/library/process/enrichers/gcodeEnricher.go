package enrichers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type gCodeEnricher struct {
	l *slog.Logger
}

func (g *gCodeEnricher) Enrich(ctx context.Context, asset *entities.Asset) error {
	g.l = slog.With("module", "gcodeEnricher").With("asset", asset.ID)
	if asset.Properties == nil {
		asset.Properties = make(entities.Properties)
	}

	fs, err := libfs.GetAssetFS(ctx, *asset)
	if err != nil {
		return fmt.Errorf("error getting fs: %w", err)
	}

	f, err := fs.Open(utils.VoZ(asset.Path))
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), ";") {
			line := strings.Trim(scanner.Text(), " ;")

			if !strings.HasPrefix(line, "thumbnail begin") {
				parseComment(asset, line)
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Join(err, errors.New("error reading gcode"))
	}
	return nil
}

func parseComment(a *entities.Asset, line string) {

	if strings.HasPrefix(line, "SuperSlicer_config") {
		a.Properties["slicer"] = "SuperSlicer"
		return
	}

	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return
	}

	if v, err := strconv.Atoi(params[1]); err != nil {
		a.Properties[params[0]] = v
		return
	}
	if v, err := strconv.ParseFloat(params[1], 64); err != nil {
		a.Properties[params[0]] = v
		return
	}
	a.Properties[params[0]] = params[1]

}
