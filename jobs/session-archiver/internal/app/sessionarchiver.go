package app

import (
	"context"
	"k8s.io/utils/env"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/browserkube/browserkube/storage"
)

//nolint:gosec // not a credentials
const nsSecret = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func New() *cli.App {
	return &cli.App{
		Name:      "session-archiver",
		Usage:     "Archives old session results into separate storage",
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Action:    Archive,
	}
}

func Archive(c *cli.Context) error {
	contextTimeout, err := time.ParseDuration(env.GetString("CONTEXT_TIMEOUT", ""))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Context, contextTimeout)
	go handleSignals(cancel)

	ns, err := getCurrentNamespace()
	if err != nil {
		return err
	}

	return archiveSessionResults(ctx, ns)
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile(nsSecret)
	if err != nil {
		return "", err
	}
	return string(ns), nil
}

func archiveSessionResults(ctx context.Context, ns string) error {
	client, err := provideClient()
	if err != nil {
		return err
	}

	blobSessionStorage, err := storage.New(context.Background(), env.GetString("BLOB_URL", ""))
	if err != nil {
		return err
	}

	blobSessionArchiveStorage, err := storage.New(context.Background(), env.GetString("BLOB_URL_ARCHIVE", ""))
	if err != nil {
		return err
	}

	archiver := &SessionResultArchiver{
		SessionResults:            client.SessionResults(ns),
		BlobSessionStorage:        blobSessionStorage,
		BlobSessionArchiveStorage: blobSessionArchiveStorage,
		ctx:                       ctx,
	}

	err = archiver.Archive()
	if err != nil {
		return err
	}

	return nil
}

func provideClient() (browserkubeclientv1.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	if err = browserkubeclientv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	browserkubeClient, err := browserkubeclientv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return browserkubeClient, nil
}

func handleSignals(cancel context.CancelFunc) {
	sigChn := make(chan os.Signal, 1)
	signal.Notify(sigChn, os.Interrupt, syscall.SIGTERM)

	for {
		sig := <-sigChn
		switch sig {
		default:
			cancel()
			return
		}
	}
}
