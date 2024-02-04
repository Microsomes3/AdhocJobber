package servers

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestLinodeCreateServer(t *testing.T) {

	linode := Linode{}

	server, err := linode.CreateServer("scheduler")

	if err != nil {
		t.Fail()
	}

	fmt.Println(server)

}

func TestLinodeDeleteServer(t *testing.T) {
	exampleInstanceId := 54589548

	linode := NewLinodeClient()

	err := linode.DeleteServer(exampleInstanceId)

	if err != nil {
		t.Fail()
	}

}

func TestGetServer(t *testing.T) {
	lini := NewLinodeClient()
	exampleServerId, err := lini.CreateServer("example")
	if err != nil {
		t.Fail()
	}

	inst, err := lini.GetServer(exampleServerId)
	if err != nil {
		t.Fail()
	}

	fmt.Println(inst)
}
