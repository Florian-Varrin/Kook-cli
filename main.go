package main

import (
	"fmt"
	"os"

	"kook/internal/cli"
)

var version = "dev" // Default version, will be overridden at build time

func main() {
	if err := cli.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
