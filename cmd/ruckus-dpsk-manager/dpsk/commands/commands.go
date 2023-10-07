package commands

import (
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/command"
)

var List []command.Command

func Register(cmd command.Command) {
	List = append(List, cmd)
}
