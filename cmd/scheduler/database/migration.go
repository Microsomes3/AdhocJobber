package database

import "gorm.io/gorm"

type DBMigration struct{}

func NewDBMigration() *DBMigration {
	return &DBMigration{}
}

type TaskDefintionModel struct {
	gorm.Model
	ID                  uint `gorm:"primaryKey"`
	Name                string
	Runner              string
	DockerImageURL      string
	DockerRegistryHost  string
	DockerAWSAccessCode string
	DockerAWSSecretCode string
	Timeout             int32
	Created             int64 `gorm:"autoCreateTime"` // Use unix seconds as creating time
	JobInstanceId       *uint
}

func (db *DBMigration) InitDB() {}
