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

	_, err := servers.GetDatabaseConnection()

	if err != nil {
		return nil, err
	}

	var tasks []*scheduler.TaskDefintion

	// db.Limit(10).Find(&tasks)

	return &scheduler.GetTasksResponse{
		Tasks: tasks,
	}, nil

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