package servers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"net/http"
	"os"
)

type Linode struct {
	Client linodego.Client
	Db     *gorm.DB
}

func NewLinodeClient() (*Linode, error) {

	db, err := GetDatabaseConnection()
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

func (lin *Linode) CreateSSH(label string) (string, error) {
	_, pub, err := generateSSHKeyPair(label)

	singleLinePubKey := strings.Join(strings.Split(pub, "\n"), "")

	if err != nil {
		return "", err
	}

	_, err = lin.Client.CreateSSHKey(context.Background(), linodego.SSHKeyCreateOptions{
		Label:  label,
		SSHKey: singleLinePubKey,
	})

	if err != nil {
		return "", err
	}

	return singleLinePubKey, nil

}

func (lin *Linode) CreateServer(label string) (int, error) {

	sskKey, err := lin.CreateSSH(label)

	if err != nil {
		return -1, errors.New("cannot create ssh key")
	}

	instance, err := lin.Client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:         "eu-west",
		Image:          "linode/ubuntu22.04",
		Label:          label,
		Type:           "g6-nanode-1",
		RootPass:       os.Getenv("DEFAULT_PASSWORD_LIN"),
		AuthorizedKeys: []string{sskKey},
	})

	if err != nil {
		return -1, err
	}

	lin.Db.Create(&JobInstance{
		ServerID: fmt.Sprint(instance.ID),
		Status:   string(instance.Status),
		Provider: "linode",
	})

	return instance.ID, nil
}

func (lin *Linode) GetServer(id int) (*linodego.Instance, error) {
	return lin.Client.GetInstance(context.Background(), id)
}

func (lin *Linode) DeleteServer(id int) error {
	return lin.Client.DeleteInstance(context.Background(), id)
}

func (lin *Linode) ExecuteCommandOnServer() {}
