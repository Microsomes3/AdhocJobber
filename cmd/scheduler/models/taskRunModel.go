package models

import (
	"time"

	"microsomes.com/scheduler/cmd/scheduler/database"
)

type TaskRunsModel struct {
	ID                   uint   `gorm:"primaryKey"`
	Status               string `default:"pending"`
	Started              int32
	Ended                int32
	JobInstanceModelId   uint
	TaskDefintionModelID uint
}

func (trm *TaskRunsModel) IsTimedOut() bool {
	//WIP

	taskDef := TaskDefintionModel{}

	db, err := database.GetDatabaseConnection()

	if err != nil {
		return true // if error assume a timeout for now WIP
	}

	tx := db.Find(&taskDef, "id=?", trm.TaskDefintionModelID)

	if tx.Error != nil {
		return true //since we cannot find a task defintion it means it must have been deleted, in this case we should kill any servers executing on this old defintion
	}

	endBy := int64(trm.Started) + int64(taskDef.Timeout)

	return time.Now().Unix() >= endBy
}
