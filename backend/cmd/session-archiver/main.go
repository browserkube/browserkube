package main

import (
	"context"
	"log"
	"os"

	"github.com/browserkube/browserkube/cmd/session-archiver/internal/app"
)

func main() {
	archiver := app.New()
	if err := archiver.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
