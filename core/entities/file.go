package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const ProjectFileType = "file"

type ProjectFile struct {
}

func (n *ProjectFile) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON string:", src))
	}
	return json.Unmarshal([]byte(str), &n)
}
func (n ProjectFile) Value() (driver.Value, error) {
	val, err := json.Marshal(n)
	return string(val), err
}

func NewProjectFile(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectFile, []*ProjectAsset, error) {
	return &ProjectFile{}, nil, nil
}

func NewProjectFile2(asset *ProjectAsset, project *Project) (*ProjectFile, error) {
	return &ProjectFile{}, nil
}
