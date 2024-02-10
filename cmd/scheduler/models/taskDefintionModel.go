package models

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
	Runs                   []TaskRunsModel `gorm:"foreignKey:TaskDefintionModelID"`
	UserID                 uint
}
