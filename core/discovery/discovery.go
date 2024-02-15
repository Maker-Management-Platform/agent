package discovery

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/duke-git/lancet/v2/maputil"
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

	project := models.NewProjectFromPath(path)

	init, _, err := DiscoverProject(project)
	if err != nil {
		return err
	}

	if init {
		project.Initialized = true

		if err := database.InsertProject(project); err != nil {
			log.Println(err)
		}
		if err := state.PersistProject(project); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func DiscoverProject(project *models.Project) (foundAssets bool, a []*models.ProjectAsset, err error) {
	projectPath := utils.ToLibPath(project.FullPath())

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		log.Println("failed to read path", projectPath)
		return false, nil, err
	}

	assets := make(map[string]*models.ProjectAsset, 0)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if e.Name() == ".project.stlib" {
			log.Println("found project", project.FullPath())
			err := loadProject(project)
			if err != nil {
				log.Printf("error loading the project %q: %v\n", project.Path, err)
				return false, nil, err
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

		assets[asset.ID] = asset
		for _, a := range nestedAssets {
			assets[a.ID] = a
		}

		foundAssets = true
	}

	if !project.Initialized {
		project.Tags = append(project.Tags, pathToTags(project.Path)...)
	}

	for _, a := range assets {
		if project.DefaultImageID == "" && a.AssetType == "image" {
			project.DefaultImageID = a.ID
		}

		err := database.InsertAsset(a)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return foundAssets, maputil.Values(assets), nil
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
	p := models.NewProject()
	_, err := toml.DecodeFile(utils.ToLibPath(fmt.Sprintf("%s/.project.stlib", project.FullPath())), &p)
	if err != nil {
		log.Printf("error decoding the project %q: %v\n", project.FullPath(), err)
		return err
	}
	project.UUID = p.UUID
	project.Description = p.Description
	project.Tags = p.Tags
	project.Assets = p.Assets
	project.DefaultImageID = p.DefaultImageID
	project.Initialized = true

	return nil
}
