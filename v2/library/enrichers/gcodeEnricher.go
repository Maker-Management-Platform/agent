package enrichers

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type gCodeEnricher struct {
	l *slog.Logger
}

func (g *gCodeEnricher) Enrich(asset *entities.Asset) func() error {
	g.l = slog.With("module", "gcodeEnricher").With("asset", asset.ID)
	return func() error {
		if asset.Properties == nil {
			asset.Properties = make(entities.Properties)
		}

		f, err := os.Open(filepath.Join(*asset.Root, *asset.Path))
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
