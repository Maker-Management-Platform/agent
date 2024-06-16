package discovery

import (
	"strings"

	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/runtime"
)

type AssetDiscoverer interface {
	Discover(root string) ([]*types.ProcessableAsset, error)
}
type ProjectDiscoverer interface {
	Discover(root string) ([]types.ProcessableProject, error)
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
