package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ory/graceful"
)

func main() {
	mux := chi.NewRouter()
	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(middleware.Logger)
	mux.Get("/", cCopy)
	mux.Post("/", cPaste)

	log.Println("main: Starting the server")
	server := graceful.WithDefaults(&http.Server{
		Addr:    ":9191",
		Handler: mux,
	})
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		log.Fatalln("main: Failed to gracefully shutdown")
	}
	log.Println("main: Server was shutdown gracefully")
}

func cCopy(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("xsel", "-b", "-o")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go func() {
		_, _ = io.Copy(w, stdout)
	}()
	err = cmd.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func cPaste(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("xsel", "-b", "-i")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	go func() {
		defer func() {
			_ = stdin.Close()
		}()
		_, _ = io.Copy(stdin, r.Body)
	}()
	w.WriteHeader(http.StatusOK)

	err = cmd.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
