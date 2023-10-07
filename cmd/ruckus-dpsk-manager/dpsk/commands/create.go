package commands

import (
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/dpsk/commands/create"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

type Create struct {
	client *client.Client
}

func init() {
	Register(&Create{})
}

func (c *Create) Name() string {
	return "create"
}

func (c *Create) Description() string {
	return "Create DPSK's"
}

func (c *Create) Handle(rc *client.Client, args []string) error {
	return create.Handle(rc.Dpsk(), args)
}
