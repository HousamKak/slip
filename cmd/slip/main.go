package main

import (
	"os"

	"slip/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
