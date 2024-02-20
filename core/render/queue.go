package render

import (
	"log"

	"github.com/eduardooliveira/stLib/core/entities"
)

type RenderJob interface {
	Name() string
	Project() *entities.Project
	Asset() *entities.ProjectAsset
	OnComplete(name string, err error)
}

var queue = make(chan RenderJob, 256)

func init() {
	go func() {
		for {
			job := <-queue
			if job.Asset().Extension == ".stl" {
				go job.OnComplete(renderStl(job))
			}
			log.Println("render queue size: ", len(queue), " - ", job.Name())
		}
	}()
}

func QueueJob(job RenderJob) {
	queue <- job
	log.Println("render queue size: ", len(queue), " + ", job.Name())
}
