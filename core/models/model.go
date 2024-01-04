package models

import (
	"archive/zip"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eduardooliveira/stLib/core/render"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
)

const ProjectModelType = "model"

var ModelExtensions = []string{".stl", ".3mf"}

type ProjectModel struct {
	ImageSha1 string `json:"image_sha1"`
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

type cacheJob struct {
	renderName string
	parent     *ProjectAsset
	model      *ProjectModel
	project    *Project
	err        chan error
}

var cacheJobs chan *cacheJob

func init() {
	log.Println("Starting", runtime.Cfg.MaxRenderWorkers, "render workers")
	cacheJobs = make(chan *cacheJob, runtime.Cfg.MaxRenderWorkers)
	go renderWorker(cacheJobs)
}

func NewProjectModel(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectModel, []*ProjectAsset, error) {
	m := &ProjectModel{}

	return m, loadImage(m, asset, project), nil
}

func loadImage(model *ProjectModel, parent *ProjectAsset, project *Project) []*ProjectAsset {
	if strings.ToLower(parent.Extension) == ".stl" {
		return []*ProjectAsset{loadStlImage(model, parent, project)}
	} else if strings.ToLower(parent.Extension) == ".3mf" {
		return load3MfImage(model, parent, project)
	}
	return nil
}

func loadStlImage(model *ProjectModel, parent *ProjectAsset, project *Project) *ProjectAsset {
	renderName := fmt.Sprintf("%s.render.png", parent.Name)
	renderPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), renderName))

	if _, err := os.Stat(renderPath); err != nil {
		errChan := make(chan error, 1)
		cacheJobs <- &cacheJob{
			renderName: renderName,
			parent:     parent,
			model:      model,
			project:    project,
			err:        errChan,
		}
		log.Println("produced", renderName)
		if err := <-errChan; err != nil {
			log.Println(err)
		}
		log.Println("terminated", renderName)
	}
	f, err := os.Open(renderPath)
	if err != nil {
		log.Println(err)
		return nil
	}

	asset, _, err := NewProjectAsset(renderName, project, f)
	if err != nil {
		log.Println(err)
		return nil
	}

	project.Assets[asset.SHA1] = asset
	model.ImageSha1 = asset.SHA1
	return asset
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
		if model.ImageSha1 == "" {
			model.ImageSha1 = asset.SHA1
		}
	}

	return rtn
}

func renderWorker(jobs <-chan *cacheJob) {
	for job := range jobs {
		go func(job *cacheJob) {
			log.Println("rendering", job.renderName)
			err := render.RenderModel(job.renderName, job.parent.Name, job.project.FullPath())
			job.err <- err
			log.Println("rendered", job.renderName)
		}(job)
	}
}
