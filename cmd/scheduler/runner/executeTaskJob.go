package runner

import (
	"sync"
	"time"

	"microsomes.com/scheduler/cmd/scheduler/models"
)

type ExecuteTaskJob struct {
	TaskRunModel models.TaskRunsModel
	M            *sync.Mutex
	NextRun      int32
}

func NewExecuteTask(taskRunModel models.TaskRunsModel) *ExecuteTaskJob {
	return &ExecuteTaskJob{
		TaskRunModel: taskRunModel,
		M:            &sync.Mutex{},
		NextRun:      0,
	}
}

func (etj *ExecuteTaskJob) ShouldExecute() bool {
	return etj.NextRun <= int32(time.Now().Unix())
}

func (etj *ExecuteTaskJob) Run() error {

	return nil
}
