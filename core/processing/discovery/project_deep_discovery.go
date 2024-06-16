package discovery

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/processing/types"
)

type DeepProjectDiscoverer struct {
}

func (DeepProjectDiscoverer) Discover(root string) ([]types.ProcessableProject, error) {
	pp := make([]types.ProcessableProject, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			return nil
		}

		folder, _ := filepath.Rel(root, path)
		if folder == "." {
			return nil
		}

		pp = append(pp, types.ProcessableProject{
			Name:     folder, //TODO: extract project name
			Path:     folder,
			Root:     root,
			FullPath: path,
		})
		return nil

	})
	return pp, err
}
