package commands

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [flags] filePath [...filePath]",
	Short: "Import files into a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Provide at least the path to one file")
		}

		projectName, _ := cmd.Flags().GetString("project")
		description, _ := cmd.Flags().GetString("description")
		rawTags, _ := cmd.Flags().GetStringSlice("tags")
		defaultImageName, _ := cmd.Flags().GetString("mainImage")

		tags := models.StringsToTags(rawTags)

		fileMap := make(map[string]io.ReadCloser)
		for _, filePath := range args {
			fileName := filepath.Base(filePath)
			file, err := os.Open(filePath)

			if err != nil {
				return err
			}
			defer file.Close()

			fileMap[fileName] = file
		}

		err := database.InitDatabase()
		if err != nil {
			return err
		}

		command := projects.NewCreateProjectCommand(
			projectName,
			"/",
			description,
			tags,
			fileMap,
			defaultImageName,
		)

		_, err = projects.CreateProject(command)

		return err
	},
}

func InitImport() *cobra.Command {
	importCmd.Flags().StringP("project", "p", "", "Project name")
	importCmd.MarkFlagRequired("project")

	importCmd.Flags().StringP("description", "d", "", "Description of the project")
	importCmd.Flags().StringSliceP("tags", "t", []string{}, "List of tags separated by a comma (tag1,tag2,...)")
	importCmd.Flags().StringP("mainImage", "", "", "Name of the image file to use as the project cover")

	return importCmd
}
