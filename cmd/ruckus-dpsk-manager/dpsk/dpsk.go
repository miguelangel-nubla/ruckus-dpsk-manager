package dpsk

import (
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/dpsk/commands"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

func Handle(rc *client.Client, args []string) error {
	if len(args) < 1 {
		return &errors.CommandInvalidError{
			Msg:      "no operation specified",
			Commands: commands.CommandList,
		}
	}

	operation := args[0]

	for _, cmd := range commands.CommandList {
		if cmd.Name() == operation {
			return cmd.Handle(rc, args[1:])
		}
	}

	return &errors.CommandInvalidError{
		Msg:      fmt.Sprintf("invalid operation specified: %s", operation),
		Commands: commands.CommandList,
	}
}
