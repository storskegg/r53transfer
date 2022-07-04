package main

import (
	"os"

	"github.com/storskegg/r53transfer/internal/application"
)

func main() {
	app := application.New()
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
