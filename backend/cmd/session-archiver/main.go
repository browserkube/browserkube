package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/browserkube/browserkube/cmd/session-archiver/internal/app"
)

func main() {
	// Create base context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	archiver := app.New()
	if err := archiver.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
