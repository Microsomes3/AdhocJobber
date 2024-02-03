package servers

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"

	"net/http"
	"os"
)

type Linode struct{}

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

func (lin *Linode) CreateServer() (bool, error) {

	linodeClient := GetClient()

	instance, err := linodeClient.CreateInstance(context.Background(), linodego.InstanceCreateOptions{
		Region: "eu-west",
		Image:  "linode/ubuntu22.04",
		Label:  "scheduler",
		Type:   "g6-nanode-1",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(instance.ID)

}

func (lin *Linode) GetServer() {}

func (lin *Linode) DeleteServer() {}

func (lin *Linode) ExecuteCommandOnServer() {}
