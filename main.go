package main

import (
	"leasepeg/cmd"
	"os"
)

var version = "dev"

func main() {
	if err := cmd.Execute(version); err != nil {
		os.Exit(1)
	}
}
