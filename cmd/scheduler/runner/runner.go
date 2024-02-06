package runner

import (
	"fmt"
	"sync"
	"time"
)

type Heartbeat struct {
	Jobs []JobRunner
}

func NewHeartbeat() *Heartbeat {
	r := NewTaskRun() //task runner job

	var jobs []JobRunner

	jobs = append(jobs, r)

	return &Heartbeat{
		Jobs: jobs,
	}
}

func (hb *Heartbeat) Beat() {

	fmt.Println("beat")

	//find all jobs that need to run and execute them in parallel goroutines

	for _, j := range hb.Jobs {
		if j.ShouldExecute() {
			fmt.Println("should execute")

		} else {
			fmt.Println("not yet")
		}
	}

}

type JobRunner interface {
	Run() error
	ShouldExecute() bool
}

func (hb *Heartbeat) StartTaskPollers(wg *sync.WaitGroup) {

	for {
		time.Sleep(time.Minute)
		hb.Beat()
	}

}
