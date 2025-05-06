package main

import (
	"os"

	"github.com/seaung/yarx-go/internal/niuma"
)

func main() {
	if err := niuma.NewNiuMaCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
