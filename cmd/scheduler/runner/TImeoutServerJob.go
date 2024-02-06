package runner

import (
	"sync"
	"time"

	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/servers"
)

type TimeoutServerJob struct {
	NextRun int32
	M       *sync.Mutex
}

func NewTimeoutServerJob() *TimeoutServerJob {
	return &TimeoutServerJob{
		NextRun: int32(time.Now().Add(time.Minute).Unix()),
	}
}

func (tsj *TimeoutServerJob) GetName() string {
	return "TimeoutServerJob"
}

func (tsj *TimeoutServerJob) ShouldExecute() bool {
	return tsj.NextRun <= int32(time.Now().Unix())
}

func (tsj *TimeoutServerJob) Run() error {

	//loop over all runs and check if timeout has apprached if so kill the server and mark the run as timedout

	var taskRuns []database.TaskRunsModel

	db, err := servers.GetDatabaseConnection()

	if err != nil {
		return err
	}

	db.Limit(1000).Find(&taskRuns)

	for _, tr := range taskRuns {
		tr.IsTimedOut()
	}

	return nil
}
