package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	scheduler "microsomes.com/scheduler/pkg/bufs"
)

func main() {

	//load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("hi")

	conn, err := grpc.Dial("localhost:4001", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	scc := scheduler.NewSchedulerServiceClient(conn)

	// taskDef, err := scc.GetTask(context.Background(), &scheduler.IdNo{
	// 	Id: 1,
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(taskDef)

	runTaskResponse, err := scc.RunTask(context.Background(), &scheduler.RunTaskRequest{
		TaskId: 1,
		Envs: []*scheduler.EnvironmentOBject{
			&scheduler.EnvironmentOBject{
				Key:   "AWS_REGION",
				Value: "eu-west-1",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(runTaskResponse.Status)

	// delTaskResponse, err := scc.DeleteTask(context.Background(), &scheduler.IdNo{
	// 	Id: 10,
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(delTaskResponse.Status)

	// // taskResponse, err := scc.GetTasks(context.Background(), &scheduler.VoidNo{})

	// // if err != nil {
	// // 	panic(err)
	// // }

	// // fmt.Println(len(taskResponse.Tasks))

	// created, err := scc.CreateTask(context.Background(), &scheduler.TaskDefintion{
	// 	Name:                "scheduler",
	// 	Runner:              "docker",
	// 	DockerImageUrl:      "public.ecr.aws/m8l7i2c5/govideocapturev8:latest",
	// 	Timeout:             60,
	// 	DockerRegistryHost:  "public.ecr.aws/m8l7i2c5",
	// 	DockerAWSAccessCode: os.Getenv("ACCESS_TOKEN"),
	// 	DockerAWSSecretCode: os.Getenv("SECRET_TOKEN"),
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(created)

}
