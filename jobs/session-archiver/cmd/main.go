package main

import (
	"log"
	"os"

	"github.com/browserkube/browserkube/session-archiver/internal/app"
)

func main() {
	archiver := app.New()
	if err := archiver.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
