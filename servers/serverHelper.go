package servers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server interface {
	CreateServer() (bool, error)
	GetServer()
	DeleteServer()
	ExecuteCommandOnServer()
}

type JobInstance struct {
	Id          uint `gorm:"primaryKey"`
	ServerID    string
	Status      string
	Provider    string
	SSHPublic   []byte
	SSHPrivate  []byte
	IPV4Address string
}

func InitDB() {

	fmt.Println("hi")

	db, err := GetDatabaseConnection()

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&JobInstance{})
	fmt.Println("auto migrate completed")
}

func (ki *JobInstance) TableName() string {
	return "servers"
}

func (ji *JobInstance) ExecuteCommand() {}

func (ji *JobInstance) SSHConnection() (*ssh.Client, error) {

	linode, err := NewLinodeClient()

	if err != nil {
		return nil, err
	}

	sid, _ := strconv.Atoi(ji.ServerID)

	instance, err := linode.GetServer(sid)

	if err != nil {
		return nil, err
	}

	if instance.Status != "running" {
		return nil, errors.New("server not running yet")

	}

	signer, err := ssh.ParsePrivateKey(ji.SSHPrivate)

	if err != nil {
		return nil, err
	}

	conf := &ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			ssh.Password(os.Getenv("DEFAULT_PASSWORD_LIN")),
		},
	}

	ip := ji.IPV4Address + ":22"

	client, err := ssh.Dial("tcp", ip, conf)

	if err != nil {

		fmt.Println(err.Error())

		return nil, err
	}

	return client, nil

}

func (ji *JobInstance) ExecuteCommands(cmds []string) error {

	cliemt, err := ji.SSHConnection()

	var b bytes.Buffer

	session.Stdout = &b

	if err := session.Start("apt-get update -y"); err != nil {
		return nil, err
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	session.Wait()

	fmt.Println(b.String())
}

func generateSSHKeyPair(label string) (privateKey, publicKey string, err error) {
	// Generate a new RSA key pair
	privateKeyRSA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKeyRSA),
	})
	privateKey = string(privateKeyBytes)

	// Generate public key string
	publicKeyBytes, err := ssh.NewPublicKey(&privateKeyRSA.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate public key: %v", err)
	}
	publicKey = string(ssh.MarshalAuthorizedKey(publicKeyBytes))

	f, err := os.Create("keys_" + label)

	f.Write([]byte(privateKey))

	if err != nil {
		return "", "", err
	}

	fmt.Println("pass")

	s := strings.Replace(publicKey, "\n\n", "\n", -1)

	return privateKey, s, nil
}

func GetDatabaseConnection() (*gorm.DB, error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/scheduler", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_IP"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
