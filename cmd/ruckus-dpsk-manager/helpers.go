package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.PrintDefaults()
}

// Errors go to stderr so a caller capturing stdout (a passphrase, a JSON list) never takes an error message for data.
func exitWithError(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// resolvePassword returns the password from -password or -password-file (exactly one), trailing newline removed.
func resolvePassword(password, passwordFile string) (string, error) {
	if password != "" && passwordFile != "" {
		return "", errors.New("use either -password or -password-file, not both")
	}
	if passwordFile != "" {
		data, err := os.ReadFile(passwordFile)
		if err != nil {
			return "", fmt.Errorf("reading password file: %v", err)
		}
		password = strings.TrimRight(string(data), "\r\n")
		if password == "" {
			return "", fmt.Errorf("password file %s is empty", passwordFile)
		}
	}
	if password == "" {
		return "", errors.New("password is required (-password-file or -password)")
	}
	return password, nil
}
