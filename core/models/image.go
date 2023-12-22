package models

import (
	"os"
)

const ProjectImageType = "image"

var ImageExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"}

type ProjectImage struct {
	*ProjectAsset
}

func NewProjectImage(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectImage, error) {
	return &ProjectImage{}, nil
}
