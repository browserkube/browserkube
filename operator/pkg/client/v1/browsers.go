package v1

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

type BrowsersInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1.BrowserList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1.Browser, error)
	Create(context.Context, *v1.Browser) (*v1.Browser, error)
	Watch(ctx context.Context, pts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	WatchByName(ctx context.Context, name string) (watch.Interface, error)
}

type browserClient struct {
	restClient rest.Interface
	ns         string
}

func (c *browserClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.BrowserList, error) {
	result := v1.BrowserList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("browsers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *browserClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Browser, error) {
	result := v1.Browser{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("browsers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *browserClient) Create(ctx context.Context, project *v1.Browser) (*v1.Browser, error) {
	result := v1.Browser{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("browsers").
		Body(project).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *browserClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("browsers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

// Delete takes name of the browser and deletes it. Returns an error if one occurs.
func (c *browserClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("browsers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *browserClient) WatchByName(ctx context.Context, name string) (watch.Interface, error) {
	pWatch, err := c.Watch(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", name).String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find selenium pod to pWatch: %v", err)
	}
	return pWatch, nil
}
