package main

import (
	"log"
	"os"

	extensioninstaller "github.com/browserkube/browserkube/extension-installer"
)

func main() {
	app := extensioninstaller.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
