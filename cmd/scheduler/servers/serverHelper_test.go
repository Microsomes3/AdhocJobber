package servers

import (
	"bytes"
	"fmt"
	"log"
	"os"
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

func TestGenerateSSHKeyPair(t *testing.T) {
	priv, pub, err := generateSSHKeyPair("example")

	if err != nil {
		t.Fail()
	}

	fmt.Println(priv)
	fmt.Println(pub)

}

func TestCanGetDBConnection(t *testing.T) {
	_, err := database.GetDatabaseConnection()

	if err != nil {
		t.Fail()
	}

}

func TestCanAddServer(t *testing.T) {
	priv, pub, err := generateSSHKeyPair("example123")

	if err != nil {
		t.Fail()
	}

	jobInstance := &models.JobInstanceModel{
		ServerCloudProviderID: "jbjb",
		InstanceStatus:        "created",
		CloudProvider:         "linode",
		SSHPrivate:            []byte(priv),
		SSHPublic:             []byte(pub),
	}

	db, _ := database.GetDatabaseConnection()

	db.Create(jobInstance)

}

func TestEstabilishSSHConnection(t *testing.T) {

	db, _ := database.GetDatabaseConnection()

	var JobInstanceModel models.JobInstanceModel

	result := db.Last(&JobInstanceModel)

	if result.Error != nil {
		t.Fail()
	}

	jobInstance := NewJobInstance(JobInstanceModel)

	_, err := jobInstance.SSHConnection()

	if err != nil {
		t.Fail()
	}

}

func TestExecuteCommands(t *testing.T) {

	db, _ := database.GetDatabaseConnection()

	commands := []string{
		// "apt-get update -y",
		// "pwd",
		"cat main.go",
	}

	var JobInstanceModel models.JobInstanceModel

	result := db.Last(&JobInstanceModel)

	if result.Error != nil {
		t.Fail()
	}
	JobInstance := NewJobInstance(JobInstanceModel)

	err := JobInstance.ExecuteCommands(commands)

	if err != nil {
		t.Fail()
	}
}

func TestUploadToJobInstance(t *testing.T) {

	db, _ := database.GetDatabaseConnection()

	var JobInstanceModel models.JobInstanceModel

	result := db.Last(&JobInstanceModel)

	if result.Error != nil {
		t.Fail()
	}

	JobInstance := NewJobInstance(JobInstanceModel)

	fi, err := os.ReadFile("../main.go")

	if err != nil {
		t.Fail()
	}

	z := bytes.NewBuffer(fi)

	// var bytes bytes.Buffer

	// fmt.Print/ln(b)

	err = JobInstance.UploadFile(*z, "main.go")

	if err != nil {
		t.Fail()
	}
}

func TestCreateTaskRun(t *testing.T) {
	db, err := database.GetDatabaseConnection()
	if err != nil {
		t.Fail()
	}

	// updateData := map[string]interface{}{
	// 	"status":               "running",
	// 	"TaskDefintionModelID": 1,
	// 	"JobInstanceModelId":   nil,
	// }

	taskRun := &models.TaskRunsModel{
		Status:               "pending nigger",
		TaskDefintionModelID: 1,
		// JobInstanceModelId:  ,
	}

	// tx := db.Model(&database.TaskRunsModel{}).Create(updateData)

	tx := db.Create(taskRun)

	if tx.Error != nil {
		t.Fail()
	}

	// fmt.Println(tx)
}
