package wd

import (
	"context"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/wd"
)

type ctxKey string

const (
	ctxSessionRemote ctxKey = "session-remote"
)

func withBrowser(ctx *wd.Context, sr *browserkubev1.Browser) *wd.Context {
	return ctx.WithValue(ctxSessionRemote, sr)
}

func getBrowser(ctx context.Context) (*browserkubev1.Browser, bool) {
	u, ok := ctx.Value(ctxSessionRemote).(*browserkubev1.Browser)
	return u, ok
}
