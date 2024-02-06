package runner

import (
	"fmt"
	"time"
)

type Heartbeat struct {
	StaticJobs    []JobRunner
	MaxJobAllowed int32
	MaxWorkers    int
	TaskQueue     chan JobRunner
	Lockstep      chan bool
}

func (h *Heartbeat) Do(task JobRunner) {
	task.Run()
	<-h.Lockstep
}

func (h *Heartbeat) StartWorkers() {
	for i := 1; i < h.MaxWorkers; i++ {
		fmt.Println("workings awaiting task")
		go func() {
			for task := range h.TaskQueue {
				h.Lockstep <- true
				h.Do(task)
			}
		}()
	}
}

func NewHeartbeat() *Heartbeat {
	r := NewTaskRun() //task runner job

	var jobs []JobRunner

	jobs = append(jobs, r)

	return &Heartbeat{
		StaticJobs:    jobs,
		MaxJobAllowed: 10,
		MaxWorkers:    10,
		TaskQueue:     make(chan JobRunner, 1000),
		Lockstep:      make(chan bool, 10),
	}
}

func (hb *Heartbeat) AddJob(job JobRunner) error {
	fmt.Println(job.GetName() + "has been scheduled")
	hb.TaskQueue <- job
	return nil
}

func (hb *Heartbeat) GetName() string {
	return "Heartbeat"
}

func (hb *Heartbeat) Beat() {

	//find all jobs that need to run and execute them in parallel goroutines

	for _, j := range hb.StaticJobs {
		if j.ShouldExecute() {
			fmt.Println("should execute")

			hb.AddJob(j)

		} else {
			fmt.Println("not yet")
		}
	}

}

type JobRunner interface {
	Run() error
	ShouldExecute() bool
	GetName() string
}

func (hb *Heartbeat) StartTaskScheduler() {

	for {
		time.Sleep(time.Second)
		hb.Beat()
	}

}
