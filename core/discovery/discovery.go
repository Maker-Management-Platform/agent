package discovery

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func Run(path string) {
	err := filepath.WalkDir(path, walker)
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
	log.Printf("walking the path %q\n", path)

	project := models.NewProjectFromPath(path)

	init, err := DiscoverProject(project)
	if err != nil {
		return err
	}

	if init {
		project.Initialized = true

		database.InsertProject(project)
		err := state.PersistProject(project)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func DiscoverProject(project *models.Project) (foundAssets bool, err error) {
	projectPath := utils.ToLibPath(project.FullPath())

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		log.Println("failed to read path", projectPath)
		return false, err
	}

	assets := make(map[string]*models.ProjectAsset, 0)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if e.Name() == ".project.stlib" {
			log.Println("found project", project.FullPath())
			err = loadProject(project)
			if err != nil {
				log.Printf("error loading the project %q: %v\n", project.Path, err)
				return false, err
			}
			continue
		}

		blacklisted := false
		for _, blacklist := range runtime.Cfg.FileBlacklist {
			if strings.HasSuffix(e.Name(), blacklist) {
				blacklisted = true
				break
			}
		}
		if blacklisted {
			continue
		}

		f, err := os.Open(utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), e.Name())))
		if err != nil {
			log.Println("failed to open file", err)
			continue
		}
		defer f.Close()
		asset, nestedAssets, err := models.NewProjectAsset(e.Name(), project, f)
		if err != nil {
			log.Println("failed create asset", err)
			continue
		}

		assets[asset.SHA1] = asset
		for _, a := range nestedAssets {
			assets[asset.SHA1] = a
		}

		foundAssets = true
	}

	if !project.Initialized {
		project.Tags = append(project.Tags, pathToTags(projectPath)...)
	}

	for _, a := range assets {
		if project.DefaultImagePath == "" && a.AssetType == "image" {
			project.DefaultImagePath = a.SHA1
		}

		err := database.InsertAsset(a)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return foundAssets, nil
}

func pathToTags(path string) []*models.Tag {

	path = strings.Trim(path, "/")
	tags := strings.Split(path, "/")
	tagSet := make(map[string]bool)
	for _, t := range tags {
		if t != "" {
			tagSet[t] = true
		}

	}
	rtn := make([]*models.Tag, len(tagSet))
	i := 0
	for k := range tagSet {
		rtn[i] = models.StringToTag(k)
		i++
	}

	return rtn
}

func loadProject(project *models.Project) error {
	_, err := toml.DecodeFile(utils.ToLibPath(fmt.Sprintf("%s/.project.stlib", project.FullPath())), &project)
	if err != nil {
		log.Printf("error decoding the project %q: %v\n", project.FullPath(), err)
		return err
	}

	return nil
}
