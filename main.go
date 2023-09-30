package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Define command line flags
	var (
		serverFlag   = flag.String("server", "https://unleashed.ruckuswireless.com", "Ruckus controller server location")
		usernameFlag = flag.String("username", "dpsk", "Username for logging in to the Ruckus controller")
		passwordFlag = flag.String("password", "", "Password for logging in to the Ruckus controller")
		debugFlag    = flag.Bool("debug", false, "Enable debug output")
		helpFlag     = flag.Bool("help", false, "Print usage information")
	)

	// Parse command line flags
	flag.Parse()

	// Check for the help flag and print usage information
	if *helpFlag {
		printUsage()
		os.Exit(0)
	}

	// Check if password was provided
	if *passwordFlag == "" {
		printUsage()
		exitWithError("Error: password is required", 1)
	}

	// Create a new RuckusClient
	client := NewRuckusClient(*serverFlag)
	client.Debug = *debugFlag

	// Perform the login
	if err := client.Login(*usernameFlag, *passwordFlag); err != nil {
		exitWithError(fmt.Sprintf("Error logging in: %v", err), 1)
	}

	// Call the start function
	os.Exit(start(client, flag.Args()))
}

func printUsage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func start(client *RuckusClient, args []string) int {
	if len(args) < 1 {
		printUsage()
		exitWithError("Error: You must provide an operation (backup, dpsk)", 1)
	}

	operation := args[0]
	switch operation {
	case "backup":
		return backup(client, args[1:])
	case "dpsk":
		return dpsk(client, args[1:])
	default:
		printUsage()
		exitWithError("Error: Invalid operation specified", 1)
	}

	return 1
}

func backup(client *RuckusClient, args []string) int {
	if len(args) < 1 {
		exitWithError("Error: You must provide an output filename", 1)
	}

	outputFilename := args[0]
	err := client.SaveBackup(outputFilename)
	if err != nil {
        exitWithError(fmt.Sprintf("Error saving backup: %v", err), 1)
	}

	return 0
}

func dpsk(client *RuckusClient, args []string) int {
	if len(args) < 2 {
		exitWithError("Error: You must provide wlanID and username", 1)
	}

	wlanID, err := strconv.Atoi(args[0])
	if err != nil {
		exitWithError("Error: wlanID must be an integer", 1)
	}

	username := args[1]

	passphrase, err := autoUserPassphrase(client, wlanID, username)
	if err != nil {
		exitWithError(err.Error(), 1)
	}

	fmt.Println(passphrase)
	return 0
}

func findUserPassphrase(client *RuckusClient, wlanID int, username string) (string, error) {
	dpskData, err := client.GetDpskData()
	if err != nil {
		return "", fmt.Errorf("error getting DPSK list: %v", err)
	}

	passphrase, err := dpskData.FindPassphrase(wlanID, username)
	if err != nil {
		return "", fmt.Errorf("error finding user: %v", err)
	}

	return passphrase, nil
}

func autoUserPassphrase(client *RuckusClient, wlanID int, username string) (string, error) {
	passphrase, err := findUserPassphrase(client, wlanID, username)
	if err == nil {
		return passphrase, nil
	}

	// User not found, create a new user with a random passphrase
	err = client.CreateDpskUser(wlanID, username)
	if err != nil {
		return "", fmt.Errorf("error creating DPSK user: %v", err)
	}

	passphrase, err = findUserPassphrase(client, wlanID, username)
	if err != nil {
		return "", fmt.Errorf("error finding user right after creating it: %v", err)
	}

	return passphrase, nil
}

func exitWithError(message string, exitCode int) {
	fmt.Println(message)
	os.Exit(exitCode)
}
