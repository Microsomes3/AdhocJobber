package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/models"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func TestCanPullDefs(t *testing.T) {
	db, err := database.GetDatabaseConnection()
	if err != nil {
		t.Fail()
	}

	var tasks []models.TaskDefintionModel

	db.Limit(10).Find(&tasks)

	fmt.Println(tasks)

}
