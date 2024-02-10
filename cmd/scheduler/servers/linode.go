package servers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"microsomes.com/scheduler/cmd/scheduler/database"
	"microsomes.com/scheduler/cmd/scheduler/models"

	"net/http"
	"os"
)

type Linode struct {
	Client linodego.Client
	Db     *gorm.DB
}

func NewLinodeClient() (*Linode, error) {

	db, err := database.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	return &Linode{
		Client: GetClient(),
		Db:     db,
	}, nil
}

func GetClient() linodego.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("LINODE_API_KEY")})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	linodeClient := linodego.NewClient(oauth2Client)
	linodeClient.SetDebug(true)
	return linodeClient
}

func (lin *Linode) CreateSSH(label string) (string, string, error) {
	priv, pub, err := generateSSHKeyPair(label)

	singleLinePubKey := strings.Join(strings.Split(pub, "\n"), "")

	if err != nil {
		return "", "", err
	}

	_, err = lin.Client.CreateSSHKey(context.Background(), linodego.SSHKeyCreateOptions{
		Label:  label,
		SSHKey: singleLinePubKey,
	})

	if err != nil {
		return "", "", err
	}

	return singleLinePubKey, priv, nil

}

func (lin *Linode) CreateServer(label string, taskRunID uint) (int, error) {

	pubKey, privKey, err := lin.CreateSSH(label)

	if err != nil {
		return -1, errors.New("cannot create ssh key")
	}

	instance, err := lin.Client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:         "eu-west",
		Image:          "linode/ubuntu22.04",
		Label:          label,
		Type:           "g6-nanode-1",
		RootPass:       os.Getenv("DEFAULT_PASSWORD_LIN"),
		AuthorizedKeys: []string{pubKey},
	})

	if err != nil {
		return -1, err
	}

	ipAddress := instance.IPv4[0]

	ip := ipAddress.String()

	jobInstanceModel := &models.JobInstanceModel{
		ServerCloudProviderID: fmt.Sprint(instance.ID),
		InstanceStatus:        string(instance.Status),
		CloudProvider:         "linode",
		SSHPublic:             []byte(pubKey),
		SSHPrivate:            []byte(privKey),
		IPV4Address:           ip,
	}

	tx := lin.Db.Create(jobInstanceModel)

	if tx.Error != nil {
		return -1, err
	}

	var taskRun models.TaskRunsModel

	tx = lin.Db.Find(&taskRun, "id=?", taskRunID)

	if tx.Error != nil {
		return -1, nil
	}

	taskRun.JobInstanceModelId = jobInstanceModel.ID
	taskRun.Started = int32(time.Now().Unix())
	taskRun.Status = "PROVISIONED"

	lin.Db.Save(&taskRun)

	return instance.ID, nil
}

func (lin *Linode) GetServer(id int) (*linodego.Instance, error) {
	return lin.Client.GetInstance(context.Background(), id)
}

func (lin *Linode) DeleteServer(id int) error {
	return lin.Client.DeleteInstance(context.Background(), id)
}

func (lin *Linode) ExecuteCommandOnServer() {}
