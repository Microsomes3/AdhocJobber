package servers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"microsomes.com/scheduler/cmd/scheduler/models"
)

type JobInstance struct {
	JI models.JobInstanceModel
}

func NewJobInstance(jobInstanceModel models.JobInstanceModel) *JobInstance {
	return &JobInstance{
		JI: jobInstanceModel,
	}
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

	sid, _ := strconv.Atoi(ji.JI.ServerCloudProviderID)

	instance, err := linode.GetServer(sid)

	if err != nil {
		return nil, err
	}

	if instance.Status != "running" {
		return nil, errors.New("server not running yet")

	}

	signer, err := ssh.ParsePrivateKey(ji.JI.SSHPrivate)

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

	ip := ji.JI.IPV4Address + ":22"

	client, err := ssh.Dial("tcp", ip, conf)

	if err != nil {

		fmt.Println(err.Error())

		return nil, err
	}

	return client, nil

}

func (ji *JobInstance) ExecuteCommands(cmds []string) error {

	client, err := ji.SSHConnection()

	if err != nil {
		return err
	}

	for _, command := range cmds {
		session, err := client.NewSession()

		var b bytes.Buffer

		if err != nil {
			return err
		}

		defer session.Close()

		session.Stdout = &b

		if err := session.Start(command); err != nil {
			return err
		}

		err = session.Wait()

		if err != nil {
			return err
		}

		fmt.Println(b.String())
	}

	return nil
}

func (ji *JobInstance) UploadFile(file bytes.Buffer, fnmame string) error {

	fmt.Println(file.String())

	client, err := ji.SSHConnection()

	if err != nil {
		return err
	}

	sftp, err := sftp.NewClient(client)

	if err != nil {
		return err
	}

	sf, err := sftp.Create(fmt.Sprintf("/root/%s", fnmame))

	if err != nil {
		return err
	}

	_, err = io.Copy(sf, &file)

	if err != nil {
		return err
	}

	// file.
	return nil
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

	fmt.Println("pass")

	s := strings.Replace(publicKey, "\n\n", "\n", -1)

	return privateKey, s, nil
}
