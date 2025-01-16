package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
)

const PlatformLinux = "linux"

func (b *BrowserImageUpdater) toBeCached(name, version string) {
	b.BrowserCache = append(b.BrowserCache, name+":"+version)
}

// CacheBrowser starts and stops pods with the newly updated pod images, so kubernetes can cache them beforehand.
func (b *BrowserImageUpdater) CacheBrowser(ns string, cache chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	wg := &sync.WaitGroup{}
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := range cache {
				select {
				case <-ctx.Done():
					return
				default:
				}
				s := strings.Split(k, ":")
				bName := s[0]
				bVersion := s[1]

				id := "update-browser-" + uuid.NewString()
				browser := &browserkubev1.Browser{
					ObjectMeta: v12.ObjectMeta{
						Name:      id,
						Namespace: ns,
					},
					Spec: browserkubev1.BrowserSpec{
						Platform:       PlatformLinux,
						BrowserVersion: bVersion,
						BrowserName:    bName,
					},
				}
				browser, err := b.Browsers.Create(ctx, browser)
				if err != nil {
					slog.Error("Err while creating.:" + err.Error())
					cancel()
					return
				}
				if browser, err = b.waitForBrowser(ctx, browser, 3*time.Minute); err != nil {
					slog.Error("Err while waiting.:" + err.Error())
					cancel()
					return
				}
				// Image cached in kubernetes. Stop browser:
				if err = b.Browsers.Delete(ctx, browser.Name, v12.DeleteOptions{}); err != nil {
					slog.Error("Err while deleting browser.:" + err.Error())
					cancel()
					return
				}
				slog.Debug("Cached: " + k)
			}
		}()
	}
	go func() {
		wg.Wait()
		cancel()
		close(b.DoneChan)
	}()
}

func (b *BrowserImageUpdater) waitForBrowser(ctx context.Context, browser *browserkubev1.Browser, timeout time.Duration) (*browserkubev1.Browser, error) {
	pWatch, err := b.Browsers.WatchByName(ctx, browser.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find browser to pWatch: %v", err)
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	var lastState *browserkubev1.Browser
	for {
		select {
		case <-timer.C:
			return lastState, errors.New("timeout exception while waiting for browser")
		case ev := <-pWatch.ResultChan():
			p, ok := ev.Object.(*browserkubev1.Browser)
			if !ok {
				continue
			}
			if ev.Type == watch.Deleted {
				return nil, errors.New("browser has been deleted after creation")
			}
			lastState = p
			switch p.Status.Phase {
			case browserkubev1.PhaseRunning:
				return p, nil
			case browserkubev1.PhasePending:
				continue
			case browserkubev1.PhaseFailed:
				return nil, fmt.Errorf("browser can't be created [%s]", p.Status.Phase)
			case "":
				// no status yet
				continue
			case browserkubev1.PhaseTerminated:
				return nil, fmt.Errorf("browser has been terminated already [%s]", p.Status.Phase)

			default:
				return nil, errors.Errorf("unknown browser state: %s", p.Status.Phase)
			}
		}
	}
}
