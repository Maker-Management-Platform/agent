package models

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io"
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
	*ProjectAsset
	ImageSha1 string `json:"image_sha1"`
}

type cacheJob struct {
	renderName string
	model      *ProjectModel
	project    *Project
	err        chan error
}

var cacheJobs chan *cacheJob

type marshalProjectModel struct {
	ImageSha1 string `json:"image_sha1"`
}

func init() {
	log.Println("Starting", runtime.Cfg.MaxRenderWorkers, "render workers")
	cacheJobs = make(chan *cacheJob, runtime.Cfg.MaxRenderWorkers)
	go renderWorker(cacheJobs)
}

func NewProjectModel(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectModel, error) {
	m := &ProjectModel{
		ProjectAsset: asset,
	}

	loadImage(m, project)

	return m, nil
}

func loadImage(model *ProjectModel, project *Project) {
	log.Println(model.ProjectAsset.Extension)

	if strings.ToLower(model.ProjectAsset.Extension) == ".stl" {
		loadStlImage(model, project)
	} else if strings.ToLower(model.ProjectAsset.Extension) == ".3mf" {
		load3MfImage(model, project)
	}

}
func loadStlImage(model *ProjectModel, project *Project) {
	renderName := fmt.Sprintf("%s.render.png", model.Name)
	renderPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), renderName))

	if _, err := os.Stat(renderPath); err != nil {
		errChan := make(chan error, 1)
		cacheJobs <- &cacheJob{
			renderName: renderName,
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
		return
	}

	asset, err := NewProjectAsset(renderName, project, f)
	if err != nil {
		log.Println(err)
		return
	}

	project.Assets[asset.SHA1] = asset
	model.ImageSha1 = asset.SHA1

}

func load3MfImage(model *ProjectModel, project *Project) {
	projectPath := utils.ToLibPath(project.FullPath())
	filePath := filepath.Join(projectPath, model.Name)
	log.Println(filePath)

	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.RemoveAll(tmpDir)

	archive, err := zip.OpenReader(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer archive.Close()

	for _, f := range archive.File {
		// Only allow image files the platform supports
		if !slices.Contains(ImageExtensions, filepath.Ext(f.Name)) {
			continue
		}

		// Ignore thumbnail since we should have the original image already
		// Also ignore metadata since it has a lot of half-render images
		if strings.Contains(f.Name, ".thumbnails/") || strings.Contains(f.Name, "Metadata/") {
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

		asset, err := NewProjectAsset(filepath.Base(outputPath), project, dstFile)
		if err != nil {
			log.Println(err)
			return
		}

		project.Assets[asset.SHA1] = asset
		// Use first image as the default
		if model.ImageSha1 == "" {
			model.ImageSha1 = asset.SHA1
		}
	}
}

func renderWorker(jobs <-chan *cacheJob) {
	for job := range jobs {
		go func(job *cacheJob) {
			log.Println("rendering", job.renderName)
			err := render.RenderModel(job.renderName, job.model.Name, job.project.FullPath())
			log.Println(err)
			job.err <- err
			log.Println("rendered", job.renderName)
		}(job)
	}
}

func (p ProjectModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalProjectModel{
		ImageSha1: p.ImageSha1,
	})
}
