package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distribution/reference"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type v2Registry struct {
}

func NewV2Registry() RegistryClient {
	return &v2Registry{}
}

func (r *v2Registry) Tags(ctx context.Context, ref reference.Named) (*RegistryImageListResp, error) {
	domain := reference.Domain(ref)
	path := reference.Path(ref)

	registryUrl := fmt.Sprintf("https://%s", domain)
	imgUrl, _ := url.JoinPath(registryUrl, "/v2", path, "/tags/list")
	// Ping the registry to ensure that is supports docker registry api v2
	if err := r.ping(registryUrl); err != nil {
		return nil, fmt.Errorf("error while pinging registry %s: %w", ref.Name(), err)
	}

	fmt.Printf("Image URL for Image %s: %s\n", ref.String(), imgUrl)

	imageList := &RegistryImageListResp{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request: %w", err)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad request")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(imageList); err != nil {
		return nil, fmt.Errorf("error while decoding json: %w", err)
	}

	return imageList, nil
}

func (r *v2Registry) ping(baseUrl string) error {
	pingUrl, err := url.JoinPath(baseUrl, "/v2/_catalog")
	if err != nil {
		return err
	}
	fmt.Printf("registry.ping url=%s\n", pingUrl)
	resp, err := http.Get(pingUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return err
}
