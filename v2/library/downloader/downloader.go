package downloader

import (
	"strings"

	"github.com/eduardooliveira/stLib/v2/library/downloader/thingiverse"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
)

type DownloadInput struct {
	Parent entities.Asset
	URL    string
	R      *repo.AssetRepo
	P      *process.Processor
}

func Download(in DownloadInput) error {
	urls := strings.Split(in.URL, ",")

	for _, url := range urls {
		if strings.Contains(url, "thingiverse.com") || strings.Contains(url, "thing:") {
			td, err := thingiverse.New(in.R, in.P)
			if err != nil {
				return err
			}
			err = td.Fetch(in.URL, &in.Parent)
			if err != nil {
				return err
			}
		} else if strings.Contains(url, "makerworld.com") {
			panic("not implemented")
		}
	}
	return nil
}
