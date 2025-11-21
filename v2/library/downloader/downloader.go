package downloader

import (
	"context"
	"net/http"
	"strings"

	"github.com/eduardooliveira/stLib/v2/library/downloader/makerworld"
	"github.com/eduardooliveira/stLib/v2/library/downloader/thingiverse"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
)

type DownloadInput struct {
	Ctx        context.Context
	Parent     entities.Asset
	URL        string
	Repo       *repo.AssetRepo
	Processor  *process.Processor
	Cookies    []*http.Cookie // For MakerWorld authentication
	UserAgent  string         // User agent for requests
}

func Download(input DownloadInput) error {
	urls := strings.FieldsFunc(input.URL, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n'
	})

	// Default user agent if not provided
	if input.UserAgent == "" {
		input.UserAgent = "Mozilla/5.0 (compatible; MMP-Agent/2.0)"
	}

	for _, url := range urls {
		if strings.Contains(url, "thingiverse.com") || strings.Contains(url, "thing:") {
			td, err := thingiverse.New(input.Ctx, input.Repo, input.Processor)
			if err != nil {
				return err
			}
			err = td.Fetch(url, input.Parent)
			if err != nil {
				return err
			}
		} else if strings.Contains(url, "makerworld.com") {
			mw, err := makerworld.New(input.Ctx, input.Repo, input.Processor, input.Cookies, input.UserAgent)
			if err != nil {
				return err
			}
			err = mw.Fetch(url, input.Parent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
