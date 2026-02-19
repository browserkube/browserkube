package main

import (
	"context"
	"log"
	"os"

	extensioninstaller "github.com/browserkube/browserkube/extension-installer"
)

func main() {
	app := extensioninstaller.NewApp()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
