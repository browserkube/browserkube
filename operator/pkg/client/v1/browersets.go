package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

type BrowsersSetsInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1.BrowserSetList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1.BrowserSet, error)
	Watch(ctx context.Context, pts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, resourceName string, b []byte, opts metav1.PatchOptions) error
}

type browserSetClient struct {
	restClient rest.Interface
	ns         string
}

func (c *browserSetClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.BrowserSetList, error) {
	result := v1.BrowserSetList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("browsersets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *browserSetClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.BrowserSet, error) {
	result := v1.BrowserSet{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("browsersets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *browserSetClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("browsersets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *browserSetClient) Patch(ctx context.Context, resourceName string, b []byte, opts metav1.PatchOptions) error {
	return c.restClient.
		Patch(types.MergePatchType).
		Namespace(c.ns).
		Resource("browsersets").
		Name(resourceName).
		Body(b).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Error()
}
