package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const ProjectImageType = "image"

var ImageExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"}

type ProjectImage struct {
}

func (n *ProjectImage) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON string:", src))
	}
	return json.Unmarshal([]byte(str), &n)
}
func (n ProjectImage) Value() (driver.Value, error) {
	val, err := json.Marshal(n)
	return string(val), err
}

func NewProjectImage(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectImage, []*ProjectAsset, error) {
	return &ProjectImage{}, nil, nil
}
