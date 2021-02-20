package main

import (
	"os"

	"github.com/JustinChen1975/lsbeat/cmd"

	_ "github.com/JustinChen1975/lsbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
