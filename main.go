package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	scheduler "microsomes.com/scheduler/bufs"
	"microsomes.com/scheduler/servers"
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

func (s *ScheduleService) CreateTask(ctx context.Context, in *scheduler.TaskDefintion) (*scheduler.CreateSuccess, error) {
	time.Sleep(time.Second * 4)
	return &scheduler.CreateSuccess{
		CreateId: 1,
		Name:     "hi",
	}, nil
}

func main() {

	//load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	servers.InitDB()

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
