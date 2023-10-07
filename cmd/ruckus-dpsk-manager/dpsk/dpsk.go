package dpsk

import (
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/dpsk/commands"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

func Handle(rc *client.Client, args []string) error {
	operation := args[0]

	for _, cmd := range commands.List {
		if cmd.Name() == operation {
			return cmd.Handle(rc, args[1:])
		}
	}

	return &errors.CommandInvalidError{
		Msg:      fmt.Sprintf("invalid operation specified: %s", operation),
		Commands: commands.List,
	}
}
