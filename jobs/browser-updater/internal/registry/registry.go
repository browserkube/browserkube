package registry

import (
	"context"
	"github.com/distribution/reference"
)

const (
	v2PageSize = 450
	// Default docker hub registry. Updater turns to this registry if any custom one is not provided
	dockerDomain       = "docker.io"
	dockerHubURI       = "https://hub.docker.com"
	browserSetPageSize = 4
	dockerHubPageSize  = 50
)

const (
	authGuideHeader  = "www-authenticate"
	authGuideRealm   = "Bearer realm"
	authGuideService = "service"
	authGuideScope   = "scope"
)

const (
	filterCharlength = 10
)

var tagsBlacklist = []string{
	"develop",
	"beta",
	"test",
	"focal",
	"jammy",
	"next",
	"next-vrt",
	"dev",
	"latest",
}

type registryManager struct {
}

type RegistryClient interface {
	Tags(ctx context.Context, ref reference.Named) (*RegistryImageListResp, error)
}

func NewRegistryManager() (RegistryClient, error) {
	return &registryManager{}, nil
}

func (r *registryManager) Tags(ctx context.Context, ref reference.Named) (*RegistryImageListResp, error) {
	domain := reference.Domain(ref)

	var imageList *RegistryImageListResp
	var err error
	if domain == "" || domain == dockerDomain {
		imageList, err = NewDockerRegistryClient().Tags(ctx, ref)
	} else {
		imageList, err = NewV2Registry().Tags(ctx, ref)
	}

	if err != nil {
		return nil, err
	}

	// Remove oldest n tags
	if len(imageList.Tags) >= browserSetPageSize {
		imageList.Tags = imageList.Tags[:browserSetPageSize]
	}
	return imageList, nil
}
