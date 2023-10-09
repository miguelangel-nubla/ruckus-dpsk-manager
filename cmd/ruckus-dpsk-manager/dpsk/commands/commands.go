package commands

import command "github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/command"

var CommandList []command.Command

func Register(cmd command.Command) {
	CommandList = append(CommandList, cmd)
}
