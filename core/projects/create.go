package projects

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
)

type CreateProjectCommand struct {
	Name             string
	Path             string
	Description      string
	Tags             []*models.Tag
	Files            map[string]io.ReadCloser
	DefaultImageName string
}

func NewCreateProjectCommand(
	name string,
	path string,
	description string,
	tags []*models.Tag,
	files map[string]io.ReadCloser, // map[FileName]File
	defaultImageName string,
) *CreateProjectCommand {
	return &CreateProjectCommand{
		Name:             name,
		Path:             path,
		Description:      description,
		Tags:             tags,
		Files:            files,
		DefaultImageName: defaultImageName,
	}
}

func CreateProject(command *CreateProjectCommand) (*models.Project, error) {
	project := models.NewProject()
	project.Name = command.Name
	project.Path = command.Path
	project.Description = command.Description
	project.Tags = command.Tags

	path := fmt.Sprintf("%s%s", runtime.Cfg.LibraryPath, project.FullPath())
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return nil, err
	}

	if err := createFiles(command, path); err != nil {
		return nil, err
	}

	if err := discoverAssets(command, project); err != nil {
		return nil, err
	}

	if err := persistProject(project); err != nil {
		return nil, err
	}

	return project, nil
}

func createFiles(command *CreateProjectCommand, path string) error {
	for fileName, src := range command.Files {
		defer src.Close()

		// Destination
		dst, err := os.Create(fmt.Sprintf("%s/%s", path, fileName))
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}

	return nil
}

func discoverAssets(command *CreateProjectCommand, project *models.Project) error {
	ok, assets, err := discovery.DiscoverProject(project)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to find assets")
	}

	if command.DefaultImageName != "" {
		for _, a := range assets {
			if a.Name == command.DefaultImageName {
				project.DefaultImageID = a.ID
			}
		}
	}

	return nil
}

func persistProject(project *models.Project) error {
	err := database.InsertProject(project)
	if err != nil {
		return err
	}

	err = state.PersistProject(project)
	if err != nil {
		return err
	}

	return nil
}
