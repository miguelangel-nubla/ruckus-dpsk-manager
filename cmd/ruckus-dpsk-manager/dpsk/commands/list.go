package commands

import (
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/dpsk/commands/list"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

type List struct {
	client *client.Client
}

func init() {
	Register(&List{})
}

func (c *List) Name() string {
	return "list"
}

func (c *List) Description() string {
	return "List DPSK's"
}

func (c *List) Handle(rc *client.Client, args []string) error {
	return list.Handle(rc.Dpsk(), args)
}
