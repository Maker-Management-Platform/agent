package downloader

import (
	"strings"

	"github.com/eduardooliveira/stLib/v2/library/downloader/thingiverse"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
)

type DownloadInput struct {
	Parent    entities.Asset
	URL       string
	Repo      *repo.AssetRepo
	Processor *process.Processor
}

func Download(input DownloadInput) error {
	urls := strings.FieldsFunc(input.URL, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n'
	})

	for _, url := range urls {
		if strings.Contains(url, "thingiverse.com") || strings.Contains(url, "thing:") {
			td, err := thingiverse.New(input.Repo, input.Processor)
			if err != nil {
				return err
			}
			err = td.Fetch(url, input.Parent) //TODO: fix only shows one thing after download
			if err != nil {
				return err
			}
		} else if strings.Contains(url, "makerworld.com") {
			panic("not implemented")
		}
	}
	return nil
}
