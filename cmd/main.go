package main

import (
	"log"
	"os"

	"github.com/sedind/cli"
)

func main() {

	app := cli.New()
	app.Description = "A sample application"
	app.UseHelpFlag = true
	app.UseVersionFlag = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
