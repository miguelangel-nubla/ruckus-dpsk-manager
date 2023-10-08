package errors

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/command"
)

type CommandError struct {
	Msg     string
	FlagSet *flag.FlagSet
}

func (e *CommandError) Error() string {
	var buf bytes.Buffer
	e.FlagSet.SetOutput(&buf)
	e.FlagSet.PrintDefaults()
	return fmt.Sprintf("%s\navailable %s:\n%s", e.Msg, e.FlagSet.Name(), buf.String())
}

type CommandInvalidError struct {
	Msg      string
	Commands []command.Command
}

func (e *CommandInvalidError) Error() string {
	var helpOutput string
	for _, cmd := range e.Commands {
		helpOutput = helpOutput + "    " + cmd.Name() + ": " + cmd.Description() + "\n"
	}

	return fmt.Sprintf("%s\navailable commands:\n%s", e.Msg, helpOutput)
}
