package runner

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	heartbeat := NewHeartbeat()

	var wg sync.WaitGroup
	wg.Add(1)

	heartbeat.StartTaskScheduler()

	wg.Wait()
}

func TestHeartBeatTaskQueue(t *testing.T) {
	heartbeat := NewHeartbeat()

	go heartbeat.StartWorkers()

	heartbeatJob := NewHeartbeatJob()

	go func() {
		time.Sleep(time.Second)

		heartbeat.AddJob(heartbeatJob)

		fmt.Println("add job")

	}()

	select {}
}
