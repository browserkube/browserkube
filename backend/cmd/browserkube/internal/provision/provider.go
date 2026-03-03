package provision

import (
	"context"
	"io"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
)

const PlatformLinux = "linux"

//go:generate mockery --name Provisioner --filename ../../playwright/mocks/Provisioner.go
type Provisioner interface {
	Provision(ctx context.Context, name string, opts *session.Capabilities) (*browserkubev1.Browser, error)
	Delete(ctx context.Context, id string) error
	Logs(ctx context.Context, id string, follow bool) (io.ReadCloser, error)
	Available(ctx context.Context) (*browserkubev1.BrowserSetList, error)
	Update(ctx context.Context, bs *browserkubev1.BrowserSet) error
}
