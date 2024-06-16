package enrichment

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/system"
	"github.com/eduardooliveira/stLib/core/utils"
)

type gCodeParser struct {
}

func NewGCodeParser() *gCodeParser {
	return &gCodeParser{}
}

func (p *gCodeParser) Parse(e types.ProcessableAsset) error {
	path := utils.ToLibPath(path.Join(e.Project.FullPath(), e.Asset.Name))
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	system.Publish("parser", e.Asset.Name)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), ";") {
			line := strings.Trim(scanner.Text(), " ;")

			if !strings.HasPrefix(line, "thumbnail begin") {
				parseComment(e.Asset, line)
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Join(err, errors.New("error reading gcode"))
	}
	return nil
}

func parseComment(a *entities.ProjectAsset, line string) {

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
