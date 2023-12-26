package models

import (
	"os"
)

const ProjectFileType = "file"

type ProjectFile struct {
	*ProjectAsset
}

func NewProjectFile(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectFile, error) {
	return &ProjectFile{}, nil
}
