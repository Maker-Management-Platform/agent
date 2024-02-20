package entities

import (
	"bufio"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/core/utils"
)

const ProjectSliceType = "slice"

var SliceExtensions = []string{".gcode"}

type tmpImg struct {
	Height int
	Width  int
	Data   []byte
}

type ProjectSlice struct {
	ImageID    string    `json:"image_id"`
	Slicer     string    `json:"slicer" toml:"slicer" form:"slicer" query:"slicer"`
	Filament   *Filament `json:"filament" toml:"filament" form:"filament" query:"filament"`
	Cost       float64   `json:"cost" toml:"cost" form:"cost" query:"cost"`
	LayerCount int       `json:"layer_count" toml:"layer_count" form:"layer_count" query:"layer_count"`
	Duration   string    `json:"duration" toml:"duration" form:"duration" query:"duration"`
}

func (n *ProjectSlice) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON string:", src))
	}
	return json.Unmarshal([]byte(str), &n)
}
func (n ProjectSlice) Value() (driver.Value, error) {
	val, err := json.Marshal(n)
	return string(val), err
}

type Filament struct {
	Length float64 `json:"length" toml:"length" form:"length" query:"length"`
	Mass   float64 `json:"mass" toml:"mass" form:"mass" query:"mass"`
	Weight float64 `json:"weight" toml:"weight" form:"weight" query:"weight"`
}

func NewProjectSlice(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectSlice, []*ProjectAsset, error) {
	s := &ProjectSlice{
		Filament: &Filament{},
	}
	err := parseGcode(s, asset, project)
	if err != nil {
		log.Println("Error parsing gecode for ", fileName, err)
		return s, nil, nil
	}

	return s, []*ProjectAsset{}, nil
}

func NewProjectSlice2(asset *ProjectAsset, project *Project) (*ProjectSlice, error) {
	asset.AssetType = ProjectSliceType
	s := &ProjectSlice{
		Filament: &Filament{},
	}
	err := parseGcode(s, asset, project)
	if err != nil {
		log.Println("Error parsing gecode for ", asset.Name, err)
		return s, nil
	}

	return s, nil
}

func parseGcode(s *ProjectSlice, parent *ProjectAsset, project *Project) error {
	path := utils.ToLibPath(path.Join(project.FullPath(), parent.Name))
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), ";") {
			line := strings.Trim(scanner.Text(), " ;")

			if !strings.HasPrefix(line, "thumbnail begin") {
				parseComment(s, line)
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Join(err, errors.New("error reading gcode"))
	}
	return nil
}

func parseComment(s *ProjectSlice, line string) {
	if strings.HasPrefix(line, "SuperSlicer_config") {
		s.Slicer = "SuperSlicer"
	} else if strings.HasPrefix(line, "filament used [mm]") {
		s.Filament.Length = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "filament used [cm3]") {
		s.Filament.Mass = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "filament used [g]") {
		s.Filament.Weight = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "filament cost") {
		s.Cost = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "total layers count") {
		s.LayerCount = parseGcodeParamInt(line)
	} else if strings.HasPrefix(line, "estimated printing time (normal mode)") {
		//https://stackoverflow.com/a/66053163/768516
		//((?P<day>\d*)d\s)?((?P<hour>\d*)h\s)?((?P<min>\d*)m\s)?((?P<sec>\d*)s)?
		s.Duration = parseGcodeParamString(line)

	}

}

func parseGcodeParamString(line string) string {
	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return ""
	}

	return params[1]
}
func parseGcodeParamInt(line string) int {
	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return 0
	}

	i, err := strconv.Atoi(params[1])
	if err != nil {
		return 0
	}

	return i
}
func parseGcodeParamFloat(line string) float64 {
	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return 0
	}

	f, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0
	}

	return f
}
