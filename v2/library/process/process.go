package process

import (
	"log/slog"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process/enrichers"
	"github.com/eduardooliveira/stLib/v2/library/process/extractors"
	"github.com/eduardooliveira/stLib/v2/library/process/renderers"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/utils"
	"golang.org/x/sync/errgroup"
)

type Processor struct {
	eg *errgroup.Group
	r  *repo.AssetRepo
}

func New(r *repo.AssetRepo) (*Processor, error) {
	eg := &errgroup.Group{}
	eg.SetLimit(10)
	renderers.Init()
	return &Processor{
		eg: eg,
		r:  r,
	}, nil
}

func (p *Processor) Process(asset *entities.Asset) *Process {
	proc := &Process{
		p:     p,
		Asset: asset,
		done:  make(chan struct{}),
	}
	if r, ok := renderers.Get(asset); ok {
		proc.renderer = r
		proc.renderState = "pending"
	} else {
		proc.renderState = "skipped"
	}
	if e, ok := extractors.Get(asset); ok {
		proc.extractor = e
		proc.extractState = "pending"
	} else {
		proc.extractState = "skipped"
	}
	if e, ok := enrichers.Get(asset); ok {
		proc.enricher = e
		proc.enrichState = "pending"
	} else {
		proc.enrichState = "skipped"
	}
	p.eg.Go(proc.Run)

	return proc
}

type Process struct {
	p            *Processor
	done         chan struct{}
	Asset        *entities.Asset
	renderer     renderers.Renderer
	renderState  string
	renderError  error
	extractor    extractors.Extractor
	extractState string
	extractError error
	enricher     enrichers.Enricher
	enrichState  string
	enrichError  error
}

func (p *Process) Wait() {
	<-p.done
}

func (p *Process) Run() error {
	defer close(p.done)
	l := slog.With("module", "process").With("asset", *p.Asset.Label)

	if p.renderer != nil {
		if img, err := p.renderer.Render(p.Asset); err != nil {
			p.renderError = err
			p.renderState = "failed"
			l.Error("failed to render asset", "error", err)
		} else {
			p.renderState = "done"
			if err := p.p.r.SaveAsset(*img); err != nil {
				l.Error("failed to save render", "error", err)
			}
			p.Asset.Thumbnail = utils.Ptr(img.ID)
		}
	}

	if p.extractor != nil {
		if _, err := p.extractor.Extract(p.Asset); err != nil {
			p.extractError = err
			p.extractState = "failed"
			l.Error("failed to extract asset", "error", err)
		} else {
			p.extractState = "done"
			p.Asset.Kind = utils.Ptr(entities.NodeKindBundle)
		}
	}

	if p.enricher != nil {
		if err := p.enricher.Enrich(p.Asset); err != nil {
			p.enrichError = err
			p.enrichState = "failed"
			l.Error("failed to enrich asset", "error", err)
		} else {
			p.enrichState = "done"
		}
	}
	if p.renderState == "done" || p.extractState == "done" || p.enrichState == "done" {
		if err := p.p.r.SaveAsset(*p.Asset); err != nil {
			l.Error("failed to save asset", "error", err)
			return err
		}
	}

	return nil
}
