package runner

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/models"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func TestTaskRunJob(t *testing.T) {

	taskRunJob := NewTaskRun()

	err := taskRunJob.Run()

	if err != nil {
		t.Fail()
	}

}

func TestUpdateRunJob(t *testing.T) {
	var run models.TaskRunsModel

	db, err := database.GetDatabaseConnection()
	if err != nil {
		t.Fail()
	}

	db.Last(&run)

	run.JobInstanceModelId = 7

	db.Save(&run)

	fmt.Println(run)
}
