package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/distribution/reference"
	"net/http"
	"time"

	"github.com/browserkube/browserkube/browser-updater/utils"
)

type dockerRegistry struct {
}

func NewDockerRegistryClient() RegistryClient {
	registry := &dockerRegistry{}
	return registry
}

func (r *dockerRegistry) Tags(ctx context.Context, ref reference.Named) (*RegistryImageListResp, error) {
	domain := reference.Domain(ref)
	if domain != dockerDomain && domain != "" {
		return nil, errors.New("incorrect registry type")
	}
	// If the secret points to docker hub or if there are no secrets defined:
	return r.getTags(ctx, ref, "", "")
}

// CheckDockerHubRegistry used for defaulting to docker hub registry in case of no config provided for something else
// or if secret specifically refers to docker hub uri
func (r *dockerRegistry) getTags(ctx context.Context, ref reference.Named, uname, password string) (*RegistryImageListResp, error) {
	// Perform Auth if any secrets are given
	var token string
	if uname != "" && password != "" {
		tokenUri := fmt.Sprintf("%s/v2/users/login", dockerHubURI)
		authVals := &dockerHubAuthReq{
			Username: uname,
			Password: password,
		}
		m, err := json.Marshal(authVals)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenUri, bytes.NewBuffer(m))
		if err != nil {
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		tokenResp := &dockerHubAuthResp{}
		err = json.NewDecoder(resp.Body).Decode(tokenResp)
		if err != nil {
			return nil, err
		}
		token = tokenResp.Token
	}

	uri := fmt.Sprintf("%s/v2/repositories/%s/tags?page_size=%d", dockerHubURI, reference.FamiliarName(ref), dockerHubPageSize)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tags := &DockerTagsResponse{}
	err = json.NewDecoder(resp.Body).Decode(tags)
	if err != nil {
		return nil, err
	}

	listResp := &RegistryImageListResp{Name: ref.Name()}
	for _, tag := range tags.Results {
		listResp.Tags = append(listResp.Tags, tag.Name)
		listResp.Digests = append(listResp.Digests, tag.Digest)
	}

	r.dockerHubFilterDigest(listResp)

	return listResp, nil
}

func (r *dockerRegistry) dockerHubFilterDigest(list *RegistryImageListResp) {
	uniqueDigests := make(map[string]interface{})
	var selected []string
	for index, tag := range list.Tags {
		if _, ok := uniqueDigests[list.Digests[index]]; !ok {
			selected = append(selected, tag)
			uniqueDigests[list.Digests[index]] = nil
		}
	}
	list.Tags = selected
}

func (r *dockerRegistry) dockerV2Auth(authGuide, uname, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if authGuide == "" {
		return "", fmt.Errorf("no header named %s is present", authGuideHeader)
	}

	headers := utils.ParseStringMap(authGuide, ",")
	if err := r.checkAuthHeaders(headers); err != nil {
		return "", err
	}
	fmt.Printf("Authenticating for %s\n", headers[authGuideRealm])
	authUrl := fmt.Sprintf("%s?service=%s&scope=%s", headers[authGuideRealm], headers[authGuideService], headers[authGuideScope])

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error while creating request: %w", err)
	}
	if uname != "" && password != "" {
		req.Header.Add("Authorization", "Basic "+utils.BasicAuth(uname, password))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while requesting auth token: %w", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	token := &registryAuthResp{}
	if err = json.NewDecoder(resp.Body).Decode(token); err != nil {
		return "", fmt.Errorf("error while decoding token response: %w", err)
	}

	if token == nil || token.AccessToken == "" {
		return "", fmt.Errorf("access token is empty")
	}
	return token.AccessToken, nil
}

func (r *dockerRegistry) checkAuthHeaders(headers map[string]string) error {
	if len(headers) == 0 {
		return fmt.Errorf("headers map len is 0")
	}
	if headers[authGuideRealm] == "" {
		return fmt.Errorf("realm key is empty")
	}
	if headers[authGuideService] == "" {
		return fmt.Errorf("service key is empty")
	}
	if headers[authGuideScope] == "" {
		return fmt.Errorf("scope key is empty")
	}
	return nil
}

// filterTags filters the tags that are received according to length, blacklist, date etc.
func (r *dockerRegistry) filterTags(tags *RegistryImageListResp) {
	var selected []string
	var selectedDigest []string
	for i, v := range tags.Tags {
		if len(v) <= filterCharlength &&
			!utils.StringSliceContains(tagsBlacklist, v) {
			selected = append(selected, v)
			if tags.Digests != nil && len(tags.Digests) > 0 {
				selectedDigest = append(selectedDigest, tags.Digests[i])
			}
		}
	}
	tags.Tags = selected
	tags.Digests = selectedDigest
}

// dockerV2SortTags sorts the tag response returned from the manifest request
// Note: Currently(blame date) docker registry v2 does not supports ordering tags by date via 1 request. Tags always returns alphabetically ordered,
// so in order to sort them by tag date we must send a separate request for every tag and process them manually
func (r *dockerRegistry) dockerV2SortTags(url, token string, tags *RegistryImageListResp) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	// Get manifests for the selected
	manifestDates := make(map[string]time.Time, len(tags.Tags))
	for _, tag := range tags.Tags {
		fmt.Printf("Processing repo: %s for tag: %s\n", tags.Name, tag)
		manifestUrl := fmt.Sprintf("%s/v2/%s/manifests/%s", url, tags.Name, tag)
		manifestReq, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestUrl, nil)
		if err != nil {
			return fmt.Errorf("error while preparing manifest req for %s: %w", tag, err)
		}
		if token != "" {
			manifestReq.Header.Add("Authorization", "Bearer "+token)
		}
		resp, err := http.DefaultClient.Do(manifestReq)
		if err != nil {
			return fmt.Errorf("error while sending manifest req for %s: %w", tag, err)
		}
		if resp != nil {
			defer resp.Body.Close()
		}
		manifest := &manifestResp{}
		if err := json.NewDecoder(resp.Body).Decode(manifest); err != nil {
			return fmt.Errorf("error while decoding manifest for %s: %w", tag, err)
		}
		// Unmarshall string history fields
		for _, field := range manifest.History {
			v1Comp := &manifestV1Compatibility{}
			if err := json.Unmarshal([]byte(field.V1Compatibility), v1Comp); err != nil {
				return fmt.Errorf("error while unmarshalling history fields: %w", err)
			}
			// Get the tag date from the first layer of history that has a created field
			if !v1Comp.Created.IsZero() {
				manifestDates[tag] = v1Comp.Created
				break
			}
		}
	}
	// Sort the times(newest to oldest)
	tags.Tags = utils.SortStringTime(manifestDates)
	return nil
}
