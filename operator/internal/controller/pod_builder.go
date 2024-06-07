package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"

	browserkubeapiv1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/operator/internal/controller/browserimage"
)

type BrowserPodBuilder interface {
	Build(ctx context.Context, b *browserkubeapiv1.Browser, opts *BrowserCtrlOpts, readinessProbe *apiv1.Probe) (*apiv1.Pod, error)
}

func NewPodBuilder(browserConfig *browserkubeapiv1.BrowserConfig) (BrowserPodBuilder, browserimage.ImageType, error) {
	// return appropriate pod builder for given browserConfig
	imageType, err := browserimage.ParseImageType(browserConfig.Image)
	if err != nil {
		return nil, -1, fmt.Errorf("unsupported browser image type: %s", browserConfig.Image)
	}

	var builder BrowserPodBuilder
	switch imageType {
	case browserimage.ImageTypeSelenium:
		builder = &seleniumPodBuilder{
			browserConfig: browserConfig,
		}
	case browserimage.ImageTypeSelenoid:
		builder = &selenoidPodBuilder{
			browserConfig: browserConfig,
		}
	case browserimage.ImageTypeAerokube:
		builder = &aerokubePodBuilder{
			browserConfig: browserConfig,
		}
	default:
		return nil, -1, fmt.Errorf("unknown image type")
	}
	return builder, imageType, nil
}
