package runner

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/servers"
)

type TaskRunJob struct {
	NextRun    int32
	ErrorCount int32
	M          *sync.Mutex
}

func (trj *TaskRunJob) GetName() string {
	return "TaskRunJob"
}

func NewTaskRun() *TaskRunJob {
	return &TaskRunJob{
		NextRun:    int32(time.Now().Add(time.Minute).Unix()),
		ErrorCount: 0,
		M:          &sync.Mutex{},
	}
}

func (trj *TaskRunJob) ShouldExecute() bool {
	return trj.NextRun <= int32(time.Now().Unix())
}

func (trj *TaskRunJob) Run() error {

	//lets execute logic then set the next run time based on my choice

	db, err := servers.GetDatabaseConnection()

	if err != nil {
		trj.ErrorCount++
		return err
	}

	var taskRuns []*database.TaskRunsModel

	tx := db.Limit(100).Find(&taskRuns)

	if tx.Error != nil {
		trj.ErrorCount++
		return err
	}

	for _, tr := range taskRuns {

		if tr.Status == "pending" {
			trj.M.Lock() //locked to prevent dublicate server creation
			defer trj.M.Unlock()

			linode, err := servers.NewLinodeClient()

			if err != nil {
				trj.ErrorCount++
				return err
			}

			id := strconv.FormatUint(uint64(tr.ID), 10)

			_, err = linode.CreateServer(fmt.Sprintf("server_%s", id), tr.ID)
			if err != nil {
				trj.ErrorCount++
				return err
			}

		}
	}

	trj.NextRun = int32(time.Now().Add(time.Minute).Unix())

	return nil
}
