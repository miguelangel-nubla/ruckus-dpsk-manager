package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager/commands"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

var (
	serverFlag     = flag.String("server", "https://unleashed.ruckuswireless.com", "Ruckus controller server location")
	usernameFlag   = flag.String("username", "dpsk", "Username for logging in to the Ruckus controller")
	passwordFlag   = flag.String("password", "", "Password for logging in to the Ruckus controller")
	caCertPathFlag = flag.String("cacert", "", "Path to a custom CA certificate")
	debugFlag      = flag.Bool("debug", false, "Enable debug output")
	helpFlag       = flag.Bool("help", false, "Print usage information")
)

func main() {
	flag.Parse()

	if *helpFlag {
		printUsage()
		os.Exit(0)
	}

	if *passwordFlag == "" {
		printUsage()
		exitWithError("Error: password is required")
	}

	ruckusClient, err := client.New(*serverFlag, *caCertPathFlag)
	if err != nil {
		exitWithError(fmt.Sprintf("Error initializing Ruckus client: %v", err))
	}

	ruckusClient.Debug = *debugFlag

	err = ruckusClient.Login(*usernameFlag, *passwordFlag)
	if err != nil {
		exitWithError(fmt.Sprintf("Error login with Ruckus client: %v", err))
	}

	args := flag.Args()
	os.Exit(start(ruckusClient, args))
}

func start(rc *client.Client, args []string) int {
	err := Handle(rc, args)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	return 0
}

func Handle(rc *client.Client, args []string) error {
	if len(args) < 1 {
		return &errors.CommandInvalidError{
			Msg:      "no command specified",
			Commands: commands.List,
		}
	}
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
