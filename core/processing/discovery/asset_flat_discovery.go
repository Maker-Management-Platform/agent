package discovery

import (
	"errors"
	"log"
	"os"

	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/utils"
)

type FlatAssetDiscoverer struct {
}

func (FlatAssetDiscoverer) Discover(root string) ([]*types.ProcessableAsset, error) {

	projectPath := utils.ToLibPath(root)

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		log.Println("failed to read path", projectPath)
		return nil, err
	}
	dAssets := make([]*types.ProcessableAsset, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if shouldSkipFile(e.Name()) {
			continue
		}
		dAssets = append(dAssets, &types.ProcessableAsset{
			Name:   e.Name(),
			Origin: "fs",
		})
	}
	if len(dAssets) == 0 {
		return nil, errors.New("not a project folder")
	}

	return dAssets, nil
}
