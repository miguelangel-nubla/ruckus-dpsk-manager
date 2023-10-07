package commands

import (
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/dpsk/commands/modify"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

type Modify struct {
	client *client.Client
}

func init() {
	Register(&Modify{})
}

func (c *Modify) Name() string {
	return "modify"
}

func (c *Modify) Description() string {
	return "Modify DPSK's"
}

func (c *Modify) Handle(rc *client.Client, args []string) error {
	return modify.Handle(rc.Dpsk(), args)
}
