package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/browserkube/browserkube/pkg/storage"
)

//nolint:gosec // not a credentials
const nsSecret = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func New() *cli.Command {
	return &cli.Command{
		Name:      "session-archiver",
		Usage:     "Archives old session results into separate storage",
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Action:    Archive,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "timeout",
				Sources:     cli.EnvVars("EXECUTION_TIMEOUT"),
				Required:    false,
				DefaultText: "1m",
			},
			&cli.StringFlag{
				Name:     "blob-url",
				Sources:  cli.EnvVars("BLOB_URL"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "blob-url-archive",
				Sources:  cli.EnvVars("BLOB_URL_ARCHIVE"),
				Required: true,
			},
		},
	}
}

func Archive(ctx context.Context, cmd *cli.Command) error {
	contextTimeout, err := time.ParseDuration(cmd.String("timeout"))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	go handleSignals(cancel)

	ns, err := getCurrentNamespace()
	if err != nil {
		return err
	}

	return archiveSessionResults(ctx, cmd, ns)
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile(nsSecret)
	if err != nil {
		return "", err
	}
	return string(ns), nil
}

func archiveSessionResults(ctx context.Context, cmd *cli.Command, ns string) error {
	client, err := provideClient()
	if err != nil {
		return err
	}

	blobSessionStorage, err := storage.New(ctx, cmd.String("blob-url"))
	if err != nil {
		return fmt.Errorf("failed to open blob storage: %w", err)
	}

	blobSessionArchiveStorage, err := storage.New(ctx, cmd.String("blob-url-archive"))
	if err != nil {
		return fmt.Errorf("failed to open blob archive storage: %w", err)
	}

	archiver := &SessionResultArchiver{
		SessionResults:            client.SessionResults(ns),
		BlobSessionStorage:        blobSessionStorage,
		BlobSessionArchiveStorage: blobSessionArchiveStorage,
	}

	err = archiver.Archive(ctx)
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
