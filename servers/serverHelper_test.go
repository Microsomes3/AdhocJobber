package servers

import (
	"fmt"
	"testing"
)

func TestGenerateSSHKeyPair(t *testing.T) {
	priv, pub, err := generateSSHKeyPair("example")

	if err != nil {
		t.Fail()
	}

	fmt.Println(priv)
	fmt.Println(pub)

}

func TestCanGetDBConnection(t *testing.T) {
	_, err := GetDatabaseConnection()

	if err != nil {
		t.Fail()
	}

}

func TestCanAddServer(t *testing.T) {
	jobInstance := &JobInstance{
		Id:       1,
		ServerID: "jbjb",
		Status:   "created",
		Provider: "linode",
	}

	db, _ := GetDatabaseConnection()

	db.Create(jobInstance)

}
