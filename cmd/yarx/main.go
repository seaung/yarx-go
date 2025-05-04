package main

import (
	"os"

	"github.com/seaung/yarx-go/internal/yarx"
)

func main() {
	if err := yarx.NewAppCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
