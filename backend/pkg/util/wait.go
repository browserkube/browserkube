package browserkubeutil

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	seleniumUpTimeout      = time.Minute
	seleniumUPRetryTimeout = time.Millisecond * 500 // 0.5 Sec
)

func SeleniumUP(ctx context.Context, u string) error {
	logger := zap.S()

	uri, _ := url.Parse(u)
	uri.Path = path.Join(uri.Path, "/status")
	if _, err := RetryWithTimeout(
		seleniumUpTimeout, seleniumUPRetryTimeout, seleniumUPRetryTimeout,
		func() (interface{}, error) {
			rq, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
			rs, err := http.DefaultClient.Do(rq)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			defer func() {
				if err := rs.Body.Close(); err != nil {
					logger.Error(err)
				}
			}()

			if rs.StatusCode != http.StatusOK {
				return nil, errors.New("incorrect http status")
			}
			return http.NoBody, nil
		}); err != nil {
		return errors.Wrapf(err, "Browser wait has failed with error: %v", err)
	}
	return nil
}
