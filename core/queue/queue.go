package queue

import "log"

type Job interface {
	JobName() string
	JobAction()
}

var queue = make(chan Job, 100)

func init() {
	go func() {
		for {
			job := <-queue
			job.JobAction()
			log.Println("job queue size: ", len(queue), " - ", job.JobName())
		}
	}()
}

func Enqueue(job Job) {
	queue <- job
	log.Println("job queue size: ", len(queue), " + ", job.JobName())
}
