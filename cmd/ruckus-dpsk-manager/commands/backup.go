package commands

import (
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/backup"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

type Backup struct {
	client *client.Client
}

func init() {
	Register(&Backup{})
}

func (c *Backup) Name() string {
	return "config"
}

func (c *Backup) Description() string {
	return "Manage backups"
}

func (c *Backup) Handle(rc *client.Client, args []string) error {
	return backup.Handle(rc, args)
}
