package servers

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"

	"net/http"
	"os"
)

type Linode struct {
	Client linodego.Client
}

func NewLinodeClient() *Linode {
	return &Linode{
		Client: GetClient(),
	}
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

func (lin *Linode) CreateServer(label string) (int, error) {

	instance, err := lin.Client.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region:   "eu-west",
		Image:    "linode/ubuntu22.04",
		Label:    label,
		Type:     "g6-nanode-1",
		RootPass: os.Getenv("DEFAULT_PASSWORD_LIN"),
	})

	if err != nil {
		return -1, err
	}

	fmt.Println(instance.ID)

	return instance.ID, nil

}

func (lin *Linode) GetServer(id int) (*linodego.Instance, error) {
	return lin.Client.GetInstance(context.Background(), id)
}

func (lin *Linode) DeleteServer(id int) error {
	return lin.Client.DeleteInstance(context.Background(), id)
}

func (lin *Linode) ExecuteCommandOnServer() {}
