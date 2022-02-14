package main

import (
	"os"

	"github.com/Ali-A-A/loran/internal/app/loran/cmd"

	_ "go.uber.org/automaxprocs"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
