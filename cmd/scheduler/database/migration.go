package database

type DBMigration struct{}

func NewDBMigration() *DBMigration {
	return &DBMigration{}
}

type TaskDefintionModel struct {
	ID                     uint `gorm:"primaryKey"`
	Name                   string
	Runner                 string
	DockerImageURL         string
	DockerRegistryHost     string
	DockerAWSAccessCode    string
	DockerAWSSecretCode    string
	DockerRegistryProvider string
	Timeout                int32
	TaskRunsModels         []TaskRunsModel
}

type TaskRunsModel struct {
	ID                   uint   `gorm:"primaryKey"`
	Status               string `default:"pending"`
	TaskDefintionModelID uint
	Started              int32
	Ended                int32
	JobInstanceModelId   uint
}
