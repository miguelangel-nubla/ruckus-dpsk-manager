package command

import "github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"

type Command interface {
	Name() string                          // Name of the subcommand
	Description() string                   // Short description for help output
	Handle(*client.Client, []string) error // Logic for handling the command
}
