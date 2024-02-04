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

func TestCreateLinodeSSHKey(t *testing.T) {
	linode, _ := NewLinodeClient()

	key, err := linode.CreateSSH("testing123")
	if err != nil {
		t.Fail()
	}

	t.Log(key)
}

func TestLinodeCreateServer(t *testing.T) {

	linode, _ := NewLinodeClient()

	server, err := linode.CreateServer("scheduler")

	if err != nil {
		t.Fail()
	}

	fmt.Println(server)

}

func TestLinodeDeleteServer(t *testing.T) {
	exampleInstanceId := 54606663

	linode, _ := NewLinodeClient()

	err := linode.DeleteServer(exampleInstanceId)

	if err != nil {
		t.Fail()
	}

}

func TestGetServer(t *testing.T) {
	lini, _ := NewLinodeClient()
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
