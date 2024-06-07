package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

type SessionResultsInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1.SessionResultList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1.SessionResult, error)
	Create(context.Context, *v1.SessionResult) (*v1.SessionResult, error)
	Watch(ctx context.Context, pts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

type sessionResultClient struct {
	restClient rest.Interface
	ns         string
}

func (c *sessionResultClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.SessionResultList, error) {
	result := v1.SessionResultList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("sessionresults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *sessionResultClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.SessionResult, error) {
	result := v1.SessionResult{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("sessionresults").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *sessionResultClient) Create(ctx context.Context, res *v1.SessionResult) (*v1.SessionResult, error) {
	result := v1.SessionResult{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("sessionresults").
		Body(res).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *sessionResultClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("sessionresults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

// Delete takes name of the browser and deletes it. Returns an error if one occurs.
func (c *sessionResultClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("sessionresults").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}
