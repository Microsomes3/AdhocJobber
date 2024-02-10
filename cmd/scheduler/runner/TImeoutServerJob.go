package runner

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/models"
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

func (tsk *TimeoutServerJob) KillServer(tm models.TaskRunsModel) error {
	var jobInstance models.JobInstanceModel

	db, err := database.GetDatabaseConnection()

	if err != nil {
		return err
	}

	tx := db.Find(&jobInstance, "id=?", tm.JobInstanceModelId)

	if tx.Error != nil {
		return err
	}

	linode, err := servers.NewLinodeClient()

	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(jobInstance.ServerCloudProviderID)

	err = linode.DeleteServer(id)

	if err != nil {
		return err
	}

	fmt.Println("server deleted")

	return nil
}

func (tsj *TimeoutServerJob) Run() error {

	//loop over all runs and check if timeout has apprached if so kill the server and mark the run as timedout

	var taskRuns []models.TaskRunsModel

	db, err := database.GetDatabaseConnection()

	if err != nil {
		return err
	}

	db.Limit(1000).Find(&taskRuns)

	for _, tr := range taskRuns {

		if tr.Status == "TIMEOUT" {
			continue
		}

		shouldKill := tr.IsTimedOut()
		if shouldKill {
			fmt.Println("should kill")

			tr.Status = "TIMEOUT"
			tr.Ended = int32(time.Now().Unix())
			tsj.KillServer(tr)

			db.Save(&tr)

		} else {
			fmt.Println("should not kill")
		}
	}

	tsj.NextRun = int32(time.Now().Add(time.Minute * 4).Unix()) //should execute every 4 minutes

	return nil
}
