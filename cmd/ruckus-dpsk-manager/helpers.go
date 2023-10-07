package main

import (
	"flag"
	"fmt"
	"os"
)

func printUsage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func exitWithError(message string) {
	fmt.Println(message)
	os.Exit(1)
}
