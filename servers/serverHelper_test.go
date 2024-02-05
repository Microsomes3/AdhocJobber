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
		log.Fatal("error loading .env file")
	}
}

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
	priv, pub, err := generateSSHKeyPair("example123")

	if err != nil {
		t.Fail()
	}

	jobInstance := &JobInstance{
		ServerID:   "jbjb",
		Status:     "created",
		Provider:   "linode",
		SSHPrivate: []byte(priv),
		SSHPublic:  []byte(pub),
	}

	db, _ := GetDatabaseConnection()

	db.Create(jobInstance)

}

func TestEstabilishSSHConnection(t *testing.T) {

	db, _ := GetDatabaseConnection()

	var JobInstance JobInstance

	result := db.Last(&JobInstance)

	if result.Error != nil {
		t.Fail()
	}

	_, err := JobInstance.SSHConnection()

	if err != nil {
		t.Fail()
	}

}
