package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (*gorm.DB, error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/scheduler2", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_IP"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

type User struct {
	ID                 uint `gorm:"primaryKey"`
	Email              string
	Password           string
	TaskDefintionModel []TaskDefintionModel `gorm:"foreignKey:UserID"`
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
	Runs                   []TaskRunsModel `gorm:"foreignKey:TaskDefintionModelID"`
	UserID                 uint
}

type TaskRunsModel struct {
	ID                   uint   `gorm:"primaryKey"`
	Status               string `default:"pending"`
	Started              int32
	Ended                int32
	JobInstanceModelId   uint
	TaskDefintionModelID uint
}
type JobInstanceModel struct {
	ID                    uint `gorm:"primaryKey"`
	ServerCloudProviderID string
	InstanceStatus        string
	CloudProvider         string
	SSHPublic             []byte
	SSHPrivate            []byte
	RootPassword          string
	IPV4Address           string
	Created               int64         `gorm:"autoCreateTime"` // Use unix seconds as creating time
	JobInstance           TaskRunsModel `gorm:"foreignKey:JobInstanceModelId"`
}

func main() {

	godotenv.Load(".env")

	db, err := GetDatabaseConnection()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&TaskDefintionModel{})
	db.AutoMigrate(&TaskRunsModel{})
	db.AutoMigrate(&JobInstanceModel{})

	fmt.Println("playing with gorm")
}
