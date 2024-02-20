package queue

import "log"

type Job interface {
	Name() string
	Run()
}

var queue = make(chan Job, 100)

func init() {
	go func() {
		for {
			job := <-queue
			job.Run()
			log.Println("job queue size: ", len(queue), " - ", job.Name())
		}
	}()
}

func Enqueue(job Job) {
	queue <- job
	log.Println("job queue size: ", len(queue), " + ", job.Name())
}
