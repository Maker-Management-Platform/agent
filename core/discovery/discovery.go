package discovery

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
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
	"golang.org/x/exp/slices"
)

func Run(path string) {
	err := filepath.WalkDir(path, walker)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
	j, _ := json.Marshal(state.Projects)
	log.Println(string(j))
}

func walker(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	log.Println(path)
	if !d.IsDir() {
		return nil
	}
	log.Printf("walking the path %q\n", path)

	project := models.NewProjectFromPath(path)

	init, err := DiscoverProjectAssets2(project)
	if err != nil {
		return err
	}

	if init {
		project.Initialized = true
		state.Projects[project.UUID] = project
		database.InsertProject(project)
		err := state.PersistProject(project)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func DiscoverProjectAssets2(project *models.Project) (foundAssets bool, err error) {
	projectPath := utils.ToLibPath(project.FullPath())

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		log.Println("failed to read path", projectPath)
		return false, err
	}
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

		if err = database.InsertAsset(asset); err != nil {
			log.Println(err)
			continue
		}

		for _, a := range nestedAssets {
			if err = database.InsertAsset(a); err != nil {
				log.Println(err)
				continue
			}
		}
		foundAssets = true
	}

	if !project.Initialized {
		project.Tags = pathToTags(project.Path)
	}

	return foundAssets, nil
}

func DiscoverProjectAssets(project *models.Project) error {
	libPath := utils.ToLibPath(project.FullPath())
	files, err := ioutil.ReadDir(libPath)
	if err != nil {
		return err
	}
	fNames, err := getDirFileSlice(files)
	if err != nil {
		log.Printf("error reading the directory %q: %v\n", libPath, err)
		return err
	}

	if slices.Contains(fNames, ".project.stlib") {
		log.Println("found project", project.FullPath())
		err = loadProject(project)
		if err != nil {
			log.Printf("error loading the project %q: %v\n", project.Path, err)
			return err
		}
	}

	if !project.Initialized {
		project.Tags = pathToTags(project.Path)
	}

	err = initProjectAssets(project, files)
	if err != nil {
		log.Printf("error loading the project %q: %v\n", project.FullPath(), err)
		return err
	}

	if project.DefaultImagePath == "" {
		for _, asset := range project.Assets {
			if asset.AssetType == models.ProjectImageType {
				project.DefaultImagePath = asset.SHA1
				break
			}
		}
	}

	return nil
}

func pathToTags(path string) []string {
	log.Println("pathToTags", path)
	path = strings.Trim(path, "/")
	tags := strings.Split(path, "/")
	tagSet := make(map[string]bool)
	for _, t := range tags {
		if t != "" {
			tagSet[t] = true
		}

	}
	rtn := make([]string, len(tagSet))
	i := 0
	for k := range tagSet {
		rtn[i] = k
		i++
	}
	log.Println("pathToTags", rtn)
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

func initProjectAssets(project *models.Project, files []fs.FileInfo) error {
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		blacklisted := false
		for _, blacklist := range runtime.Cfg.FileBlacklist {
			if strings.HasSuffix(file.Name(), blacklist) {
				blacklisted = true
				break
			}
		}
		if blacklisted {
			continue
		}
		f, err := os.Open(utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), file.Name())))
		if err != nil {
			return err
		}
		defer f.Close()
		asset, _, err := models.NewProjectAsset(file.Name(), project, f)

		if err != nil {
			return err
		}

		project.Assets[asset.SHA1] = asset
		state.Assets[asset.SHA1] = asset
		database.InsertAsset(asset)

	}

	return nil
}

func getDirFileSlice(files []fs.FileInfo) ([]string, error) {

	fNames := make([]string, 0)
	for _, file := range files {
		fNames = append(fNames, file.Name())
	}

	return fNames, nil
}
