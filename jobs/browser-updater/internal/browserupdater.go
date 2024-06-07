package internal

import (
	"context"
	"encoding/json"
	"fmt"
	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log/slog"

	"github.com/browserkube/browserkube/browser-updater/internal/registry"
)

func UpdateBrowserImages(ctx context.Context, clientset *kubernetes.Clientset, bkClient browserkubeclientv1.Interface, ns string) error {
	slog.Info("Running browser image updater")

	browserSets := bkClient.BrowserSets(ns)

	registryClient, err := registry.NewRegistryManager()
	if err != nil {
		return err
	}

	bsi, err := browserSets.List(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error while running image updater. err: %s", err.Error())
	}

	imageUpdater := &BrowserImageUpdater{
		BrowserCache:   []string{},
		DoneChan:       make(chan struct{}),
		Browsers:       bkClient.Browsers(ns),
		RegistryClient: registryClient,
	}

	bsi, err = imageUpdater.UpdateBrowserSet(ctx, bsi)
	if err != nil {
		return fmt.Errorf("error while running image updater. err: %s", err.Error())
	}

	b, err := json.Marshal(bsi.Items[0])
	if err != nil {
		return fmt.Errorf("error while marshalling json: %s", err.Error())
	}

	if err = browserSets.Patch(ctx, bsi.Items[0].Name, b, v1.PatchOptions{}); err != nil {
		return fmt.Errorf("error while patching browserset. err: %s", err.Error())
	}
	slog.Info("BrowserSet patched. Moving to caching...")

	cache := make(chan string)
	imageUpdater.CacheBrowser(ns, cache)

	go func() {
		for _, j := range imageUpdater.BrowserCache {
			cache <- j
		}
		close(cache)
	}()
	<-imageUpdater.DoneChan

	slog.Info("Browser image updater is done")
	return nil
}
