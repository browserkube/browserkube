package main

import (
	"github.com/browserkube/browserkube/browser-updater/internal"
	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

//nolint:gosec // not a credentials
const nsSecret = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func main() {
	app := &cli.App{
		Name:      "browser-updater",
		Usage:     "Checks remote registries for new versions of browser images and updates accordingly",
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "namespace",
				Aliases:  []string{"ns"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "kubeconfig",
				Aliases:  []string{"k"},
				Required: false,
			},
		},
		Action: runUpdate,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runUpdate(ctx *cli.Context) error {
	var ns string
	var err error
	if ns = ctx.String("namespace"); ns == "" {
		ns, err = getCurrentNamespace()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	clientSet, bkClient, err := buildClientSet(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = internal.UpdateBrowserImages(ctx.Context, clientSet, bkClient, ns)
	return err
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile(nsSecret)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(ns), nil
}

func buildClientSet(ctx *cli.Context) (*kubernetes.Clientset, browserkubeclientv1.Interface, error) {
	var clientset *kubernetes.Clientset
	var err error

	var config *rest.Config
	if kubeconfig := ctx.String("kubeconfig"); kubeconfig != "" {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if err = browserkubeclientv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	browserkubeClient, err := browserkubeclientv1.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return clientset, browserkubeClient, errors.WithStack(err)
}
