package types

import "github.com/eduardooliveira/stLib/core/entities"

type ProcessableAsset struct {
	Name    string
	Label   string
	Project *entities.Project
	Asset   *entities.ProjectAsset
	Origin  string
}

type ProcessableProject struct {
	Name     string
	Path     string
	Root     string
	FullPath string
}
