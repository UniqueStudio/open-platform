package main

import (
	"os"

	"github.com/UniqueStudio/open-platform/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
