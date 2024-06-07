package v1

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
)

type Interface interface {
	RESTClient() rest.Interface
	Browsers(namespace string) BrowsersInterface
	BrowserSets(namespace string) BrowsersSetsInterface
	SessionResults(namespace string) SessionResultsInterface
}

type browserkubeV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (Interface, error) {
	config := *c
	setConfigDefaults(&config)
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &browserkubeV1Client{restClient: client}, nil
}

func setConfigDefaults(config *rest.Config) {
	gv := browserkubev1.GroupVersion

	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
}

func (c *browserkubeV1Client) Browsers(namespace string) BrowsersInterface {
	return &browserClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *browserkubeV1Client) BrowserSets(namespace string) BrowsersSetsInterface {
	return &browserSetClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *browserkubeV1Client) SessionResults(namespace string) SessionResultsInterface {
	return &sessionResultClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *browserkubeV1Client) RESTClient() rest.Interface {
	return c.restClient
}
