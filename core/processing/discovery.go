package processing

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/system"
	"github.com/eduardooliveira/stLib/core/utils"
)

func Run(path string) {
	tempPath := filepath.Clean(filepath.Join(runtime.GetDataPath(), "assets"))
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err := os.MkdirAll(tempPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}
	time.Sleep(5 * time.Second)
	system.Publish("discovery.scan", map[string]any{"state": "started"})
	err := filepath.WalkDir(path, walker)
	system.Publish("discovery.scan", map[string]any{"state": "finished"})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
}

func walker(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	if !d.IsDir() {
		return nil
	}

	folder, _ := filepath.Rel(runtime.Cfg.Library.Path, path)
	if folder == "." {
		return nil
	}
	_, err = HandlePath(folder)
	return err
}

func HandlePath(folder string) (*entities.Project, error) {
	newProject := true
	project := entities.NewProjectFromPath(folder)
	if p, err := database.GetProjectByPathAndName(project.Path, project.Name); err == nil {
		project = p
		newProject = false
	}

	dAssets, err := discoverAssets(project)
	if err != nil {
		return nil, err
	}

	if newProject {
		project.Tags = append(project.Tags, pathToTags(project.Path)...)
	}

	if len(dAssets) > 0 {
		if newProject {
			if err := utils.CreateAssetsFolder(project.UUID); err != nil {
				log.Println(err)
				return nil, err
			}
			if err := database.InsertProject(project); err != nil {
				log.Println(err)
			}
		}

		for _, dAsset := range dAssets {
			EnqueueInitJob(dAsset)
		}
	}
	return project, nil
}

func discoverAssets(project *entities.Project) (assets []*ProcessableAsset, err error) {
	projectPath := utils.ToLibPath(project.FullPath())

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		log.Println("failed to read path", projectPath)
		return nil, err
	}
	dAssets := make([]*ProcessableAsset, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if shouldSkipFile(e.Name()) {
			continue
		}
		dAssets = append(dAssets, &ProcessableAsset{
			Name:    e.Name(),
			Project: project,
			Origin:  "fs",
		})

	}

	return dAssets, nil
}

func shouldSkipFile(name string) bool {

	if strings.HasPrefix(name, ".") {
		if runtime.Cfg.Library.IgnoreDotFiles {
			return true
		}
	}

	for _, blacklist := range runtime.Cfg.Library.Blacklist {
		if strings.HasSuffix(name, blacklist) {
			return true
		}
	}

	return false
}

func pathToTags(path string) []*entities.Tag {

	path = strings.Trim(path, "/")
	tags := strings.Split(path, "/")
	tagSet := make(map[string]bool)
	for _, t := range tags {
		if t != "" {
			tagSet[t] = true
		}

	}
	rtn := make([]*entities.Tag, len(tagSet))
	i := 0
	for k := range tagSet {
		rtn[i] = entities.StringToTag(k)
		i++
	}

	return rtn
}
