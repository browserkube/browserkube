package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/fleetframework/goga/cmds/recorder/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	rcrd := internal.NewApp()
	wait := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/recorder/stop", stop(cancel, wait))
		err := http.ListenAndServe(":5555", mux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := rcrd.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
	wait <- struct{}{}
	wg.Wait()
}

func stop(cancel context.CancelFunc, wait chan struct{}) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		cancel()
		<-wait
	}
}
