package main

import (
	"os"

	_ "go.uber.org/automaxprocs"

	yarxgo "github.com/seaung/yarx-go/internal/yarx-go"
)

func main() {
	cmd := yarxgo.NewYarxgoCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
