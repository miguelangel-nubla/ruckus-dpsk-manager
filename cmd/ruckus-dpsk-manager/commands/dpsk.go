package commands

import (
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/dpsk"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

type Dpsk struct {
	client *client.Client
}

func init() {
	Register(&Dpsk{})
}

func (c *Dpsk) Name() string {
	return "dpsk"
}

func (c *Dpsk) Description() string {
	return "Manage DPSK's"
}

func (c *Dpsk) Handle(rc *client.Client, args []string) error {
	return dpsk.Handle(rc, args)
}
