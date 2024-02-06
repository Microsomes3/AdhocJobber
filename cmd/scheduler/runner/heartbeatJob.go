package runner

import (
	"fmt"
	"sync"
	"time"
)

type HeartbeatJob struct {
	M       *sync.Mutex
	NextRun int32
}

func NewHeartbeatJob() *HeartbeatJob {
	return &HeartbeatJob{
		NextRun: int32(time.Now().Unix()),
		M:       &sync.Mutex{},
	}
}
func (etj *HeartbeatJob) ShouldExecute() bool {
	return etj.NextRun <= int32(time.Now().Unix())
}

func (etj *HeartbeatJob) GetName() string {
	return "HeartbeatJob"
}

func (etj *HeartbeatJob) Run() error {
	etj.M.Lock()
	fmt.Println("heartbeat boop")
	time.Sleep(time.Second * 5)
	fmt.Println("heartbeat done")

	etj.M.Unlock()

	etj.NextRun = int32(time.Now().Add(time.Second * 20).Unix())

	return nil
}
