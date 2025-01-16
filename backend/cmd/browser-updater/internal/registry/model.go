package registry

import (
	"encoding/json"
	"time"
)

type k8sRegistrySecret struct {
	Auths map[string]registryServerCreds `json:"auths"`
}

type registryServerCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}

func (k *k8sRegistrySecret) UnmarshalJSON(in []byte) error {
	tmp := map[string]json.RawMessage{}
	if err := json.Unmarshal(in, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp["auths"], &k.Auths); err != nil {
		return err
	}
	return nil
}

type DockerTagsResponse struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

type Result struct {
	Name   string `json:"name"`
	Digest string `json:"digest"`
}

type RegistryImageListResp struct {
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Digests []string `json:"-"`
}

type registryAuthResp struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

type dockerHubAuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type dockerHubAuthResp struct {
	Token string `json:"token"`
}

type manifestResp struct {
	Name         string            `json:"name"`
	Tag          string            `json:"tag"`
	Architecture string            `json:"architecture"`
	History      []manifestHistory `json:"history"`
}

type manifestHistory struct {
	V1Compatibility string `json:"v1Compatibility"`
}

type manifestV1Compatibility struct {
	Created time.Time `json:"created,omitempty"`
}
