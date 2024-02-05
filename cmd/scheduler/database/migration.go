package database

type DBMigration struct{}

func NewDBMigration() *DBMigration {
	return &DBMigration{}
}

type TaskDefintionModel struct {
	ID                  uint `gorm:"primaryKey"`
	Name                string
	Runner              string
	DockerImageURL      string
	DockerRegistryHost  string
	DockerAWSAccessCode string
	DockerAWSSecretCode string
	Timeout             int32
	JobInstanceModelId  *uint
}

func (db *DBMigration) InitDB() {}
