package internal

import (
	"context"
	"fmt"
	"github.com/browserkube/browserkube/browser-updater/internal/registry"
	"github.com/browserkube/browserkube/browser-updater/utils"
	apiv1 "github.com/browserkube/browserkube/operator/api/v1"
	clientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/distribution/reference"
	"log/slog"
)

type BrowserImageUpdater struct {
	// This keeps unique browser versions, so we can start a pod and 'cache' them to kubernetes cluster.
	// Does not affect functionality if backend pod stops/restarts. Just for making provisioners job easier.
	BrowserCache   []string
	DoneChan       chan struct{}
	BrowserSets    clientv1.BrowsersSetsInterface
	Browsers       clientv1.BrowsersInterface
	RegistryClient registry.RegistryClient
}

type DockerTagsResponse struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

type QuayTagsResponse struct {
	Tags []Result `json:"tags"`
}

type Result struct {
	Name string `json:"name"`
}

func (b *BrowserImageUpdater) UpdateBrowserSet(ctx context.Context, browserSet *apiv1.BrowserSetList) (*apiv1.BrowserSetList, error) {
	// Check browser versions
	if len(browserSet.Items) == 0 {
		return nil, fmt.Errorf("error :no browsersets found")
	}
	spec := browserSet.Items[0].Spec

	updateFunc := func(spec map[string]apiv1.BrowsersConfig) error {
		images := b.getImagesToUpdate(spec)
		for browser, imgRefs := range images {
			for _, imgRef := range imgRefs {
				resp, err := b.RegistryClient.Tags(ctx, imgRef)
				if err != nil {
					return fmt.Errorf("error while requesting registry. img:%s, err: %w", imgRef.String(), err)
				}

				// update our browserset
				for _, result := range resp.Tags {
					var taggedRef reference.NamedTagged
					taggedRef, err = reference.WithTag(imgRef, result)
					if err != nil {
						return err
					}
					spec[browser].Versions[result] = apiv1.BrowserConfig{
						Provider: "k8s",
						Image:    taggedRef.String(),
						Port:     "4444",
					}
					b.toBeCached(browser, result)
				}
			}
		}
		return nil
	}

	// Webdriver
	if err := updateFunc(spec.WebDriver); err != nil {
		return nil, err
	}

	// Playwright
	if err := updateFunc(spec.Playwright); err != nil {
		return nil, err
	}
	return browserSet, nil
}

// checkBrowserVersion checks the browsersets and gets the unique image names without the image versions
func (b *BrowserImageUpdater) getImagesToUpdate(spec map[string]apiv1.BrowsersConfig) map[string][]reference.Named {
	uniqueImages := make(map[string][]reference.Named)
	for browserName, v := range spec {
		for _, version := range v.Versions {
			imgRef, err := reference.ParseNormalizedNamed(version.Image)
			if err != nil {
				slog.Error("unable to parse image", "error", err)
				continue
			}

			if _, ok := uniqueImages[browserName]; !ok {
				uniqueImages[browserName] = []reference.Named{}
			}

			// check image name only, ignore tags
			if utils.SliceContains(uniqueImages[browserName], imgRef, func(o1 reference.Named, o2 reference.Named) bool {
				return o1.Name() == o2.Name()
			}) {
				continue
			}
			uniqueImages[browserName] = append(uniqueImages[browserName], imgRef)
		}
	}
	return uniqueImages
}
