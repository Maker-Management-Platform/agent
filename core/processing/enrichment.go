package processing

import (
	"log"

	"github.com/eduardooliveira/stLib/core/processing/enrichment"
)

var enrichmentQueue = make(chan *processableAsset, 256)

func init() {
	renderers := make(map[string]enrichment.Renderer, 0)
	renderers[".stl"] = enrichment.NewSTLRenderer()
	renderers[".gcode"] = enrichment.NewGCodeRenderer()
	parsers := make(map[string]enrichment.Parser, 0)
	parsers[".gcode"] = enrichment.NewGCodeParser()
	extractors := make(map[string]enrichment.Extractor, 0)
	extractors[".3mf"] = enrichment.New3MFExtractor()

	go enrichementRoutine(renderers, extractors, parsers)
}

func QueueEnrichmentJob(job *processableAsset) {
	enrichmentQueue <- job
	log.Println("enrichment queue size: ", len(enrichmentQueue), " + ", job.Name())
}

func enrichementRoutine(renderers map[string]enrichment.Renderer, extractors map[string]enrichment.Extractor, parsers map[string]enrichment.Parser) {
	for {
		job := <-enrichmentQueue
		if renderer, ok := renderers[job.asset.Extension]; ok {
			if err := render(job, renderer); err != nil {
				log.Println(err)
			}
		}
		if extractor, ok := extractors[job.asset.Extension]; ok {
			if err := extract(job, extractor); err != nil {
				log.Println(err)
			}
		}
		if parser, ok := parsers[job.asset.Extension]; ok {
			if err := parser.Parse(job); err != nil {
				log.Println(err)
			}
		}

		job.OnEnrichmentComplete(nil)
		log.Println("enrichment queue size: ", len(enrichmentQueue), " - ", job.Name())
	}
}

func extract(p *processableAsset, extractor enrichment.Extractor) error {
	excracted, err := extractor.Extract(p)
	if err != nil {
		return err
	}
	for _, e := range excracted {
		EnqueueInitJob(&processableAsset{
			name:    e.File,
			label:   e.Label,
			parent:  p.asset,
			project: p.project,
			origin:  "extract",
		})
	}
	return nil
}

func render(p *processableAsset, renderer enrichment.Renderer) error {
	file, err := renderer.Render(p)
	if err != nil {
		return err
	}
	EnqueueInitJob(&processableAsset{
		name:    file,
		parent:  p.asset,
		project: p.project,
		origin:  "render",
	})
	return nil
}
