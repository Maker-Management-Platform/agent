package entities

import (
	"archive/zip"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eduardooliveira/stLib/core/utils"
)

const ProjectModelType = "model"

var ModelExtensions = []string{".stl", ".3mf"}

type ProjectModel struct {
	ImageID string `json:"image_id"`
}

func (n *ProjectModel) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON string:", src))
	}
	return json.Unmarshal([]byte(str), &n)
}
func (n ProjectModel) Value() (driver.Value, error) {
	val, err := json.Marshal(n)
	return string(val), err
}

type renderJob struct {
	project    *Project
	renderName string
	renderPath string
	modelPath  string
}

func NewProjectModel(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectModel, []*ProjectAsset, error) {
	m := &ProjectModel{}

	return m, loadImage(m, asset, project), nil
}
func NewProjectModel2(asset *ProjectAsset, project *Project) (*ProjectModel, error) {
	m := &ProjectModel{}

	return m, nil
}

func loadImage(model *ProjectModel, parent *ProjectAsset, project *Project) []*ProjectAsset {
	if strings.ToLower(parent.Extension) == ".stl" {
		loadStlImage(model, parent, project)
		return []*ProjectAsset{}
	} else if strings.ToLower(parent.Extension) == ".3mf" {
		return load3MfImage(model, parent, project)
	}
	return nil
}

func loadStlImage(model *ProjectModel, parent *ProjectAsset, project *Project) {
	renderName := fmt.Sprintf("%s.render.png", parent.Name)
	renderPath := utils.ToLibPath(path.Join(project.FullPath(), renderName))

	if _, err := os.Stat(renderPath); err != nil {
		/*render.QueueJob(&renderJob{
			project:    project,
			renderName: renderName,
			renderPath: renderPath,
			modelPath:  utils.ToLibPath(path.Join(project.FullPath(), parent.Name)),
		})*/
	}
}

func load3MfImage(model *ProjectModel, parent *ProjectAsset, project *Project) []*ProjectAsset {
	rtn := make([]*ProjectAsset, 0)
	projectPath := utils.ToLibPath(project.FullPath())
	filePath := filepath.Join(projectPath, parent.Name)
	log.Println(filePath)

	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer os.RemoveAll(tmpDir)

	archive, err := zip.OpenReader(filePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer archive.Close()

	for _, f := range archive.File {
		// Only allow image files the platform supports
		if !slices.Contains(ImageExtensions, filepath.Ext(f.Name)) {
			continue
		}

		// Ignore thumbnail since we should have the original image already
		if strings.Contains(f.Name, ".thumbnails/") {
			continue
		}

		outputPath := filepath.Join(projectPath, filepath.Base(f.Name))
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			log.Println(err)
			continue
		}

		dstFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Println(err)
			continue
		}
		defer dstFile.Close()

		fileInArchive, err := f.Open()
		if err != nil {
			log.Println(err)
			continue
		}
		defer fileInArchive.Close()

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			log.Println(err)
			continue
		}

		asset, _, err := NewProjectAsset(filepath.Base(outputPath), project, dstFile)
		if err != nil {
			log.Println(err)
			return nil
		}

		rtn = append(rtn, asset)

		// Use first image as the default
		if model.ImageID == "" {
			model.ImageID = asset.ID
		}
	}

	return rtn
}

func (job *renderJob) ModelPath() string {
	return job.modelPath
}

func (job *renderJob) RenderPath() string {
	return job.renderPath
}

func (job *renderJob) OnComplete(err error) {
	f, err := os.Open(job.renderPath)
	if err != nil {
		log.Println(err)
	}

	asset, _, err := NewProjectAsset(job.renderName, job.project, f)
	if err != nil {
		log.Println(err)
	}
	log.Println("rendering complete : ", asset.Name)
}
