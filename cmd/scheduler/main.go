package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/servers"
	scheduler "microsomes.com/scheduler/pkg/bufs"
)

type ScheduleService struct {
	scheduler.UnimplementedSchedulerServiceServer
}

func (s *ScheduleService) RunTask(ctx context.Context, in *scheduler.RunTaskRequest) (*scheduler.RunTaskResponse, error) {
	fmt.Println("scheduling task")
	time.Sleep(time.Second)
	return &scheduler.RunTaskResponse{
		StatusCode: 201,
	}, nil
}

func (s *ScheduleService) CreateTask(ctx context.Context, in *scheduler.TaskDefintion) (*scheduler.CreateTaskResponse, error) {

	db, err := servers.GetDatabaseConnection()

	if err != nil {
		return nil, err
	}

	td := &database.TaskDefintionModel{
		Name:                in.Name,
		Runner:              in.Runner,
		DockerImageURL:      in.DockerImageUrl,
		Timeout:             in.Timeout,
		DockerRegistryHost:  in.DockerRegistryHost,
		DockerAWSAccessCode: in.DockerAWSAccessCode,
		DockerAWSSecretCode: in.DockerAWSSecretCode,
	}

	tx := db.Create(td)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &scheduler.CreateTaskResponse{
		CreateId: int32(td.ID),
	}, nil
}

func (s *ScheduleService) GetTasks(ctx context.Context, in *scheduler.VoidNo) (*scheduler.GetTasksResponse, error) {

	db, err := servers.GetDatabaseConnection()

	if err != nil {
		return nil, err
	}

	var tasks []*scheduler.TaskDefintion

	var taskItems []*database.TaskDefintionModel

	db.Limit(1000).Find(&taskItems)

	for _, i := range taskItems {
		tasks = append(tasks, &scheduler.TaskDefintion{
			Name:                i.Name,
			Runner:              i.Runner,
			DockerImageUrl:      i.DockerImageURL,
			Timeout:             i.Timeout,
			DockerRegistryHost:  i.DockerRegistryHost,
			DockerAWSAccessCode: i.DockerAWSAccessCode,
			DockerAWSSecretCode: i.DockerAWSSecretCode,
		})

	}

	return &scheduler.GetTasksResponse{
		Tasks: tasks,
	}, nil

}

func (s *ScheduleService) GetTask(ctx context.Context, in *scheduler.IdNo) (*scheduler.TaskDefintion, error) {
	var taskDef *database.TaskDefintionModel

	db, err := servers.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	tx := db.First(&taskDef, "id=?", in.Id)

	if tx.Error != nil {
		return nil, err
	}

	var taskDefR scheduler.TaskDefintion

	taskDefR.Name = taskDef.Name
	taskDefR.Runner = taskDef.Runner
	taskDefR.Timeout = taskDef.Timeout
	taskDefR.DockerImageUrl = taskDef.DockerImageURL
	taskDefR.DockerRegistryHost = taskDef.DockerRegistryHost
	taskDefR.DockerAWSAccessCode = taskDef.DockerAWSAccessCode
	taskDefR.DockerAWSSecretCode = taskDef.DockerAWSSecretCode
	taskDefR.DockerRegistryProvider = taskDef.DockerRegistryProvider

	return &taskDefR, nil
}

func InitDB() {

	fmt.Println("hi")

	db, err := servers.GetDatabaseConnection()

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&servers.JobInstanceModel{})
	db.AutoMigrate(database.TaskDefintionModel{})
	fmt.Println("auto migrate completed")
}
func main() {

	//load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	InitDB()

	fmt.Println("Scheduler recording service reporting for duty")
	lis, err := net.Listen("tcp", "localhost:4000")

	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{}

	server := grpc.NewServer(opts...)

	// server.RegisterService(S)

	sc := ScheduleService{}

	scheduler.RegisterSchedulerServiceServer(server, &sc)

	server.Serve(lis)

}
