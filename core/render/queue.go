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
type renderer interface {
	render(job RenderJob) (string, error)
}

var queue = make(chan RenderJob, 256)

func init() {
	go func() {
		for {
			job := <-queue
			if job.Asset().Extension == ".stl" {
				go job.OnComplete((&stlRenderer{}).render(job))
			} else if job.Asset().Extension == ".gcode" {
				go job.OnComplete((&gcodeRenderer{}).render(job))
			}
			log.Println("render queue size: ", len(queue), " - ", job.Name())
		}
	}()
}

func QueueJob(job RenderJob) {
	queue <- job
	log.Println("render queue size: ", len(queue), " + ", job.Name())
}
