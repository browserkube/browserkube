package provisionk8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/ptr"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/browserkube/browserkube/pkg/session"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

const (
	browserUPTimeout           = time.Minute
	podGracefulShutdownTimeout = 30
)

// container names
const (
	browserContainerName = "browser"
)

type CreationErr struct {
	error
}

// k8sWebDriverProvisioner provisioner for k8s
type k8sWebDriverProvisioner struct {
	logger            *zap.SugaredLogger
	clientset         *kubernetes.Clientset
	podClient         corev1.PodInterface
	envConfig         *provision.Config
	browsersClient    browserkubeclientv1.BrowsersInterface
	browserSetsClient browserkubeclientv1.BrowsersSetsInterface
}

func newK8sWebDriverProvisioner(
	clientset kubernetes.Interface,
	browserkubeClient browserkubeclientv1.Interface,
	envConfig *provision.Config,
) *k8sWebDriverProvisioner {
	logger := zap.S()
	logger.Infof("Browser Namespace: %s", envConfig.BrowserNS)
	return &k8sWebDriverProvisioner{
		logger:            logger,
		envConfig:         envConfig,
		browsersClient:    browserkubeClient.Browsers(envConfig.BrowserNS),
		browserSetsClient: browserkubeClient.BrowserSets(envConfig.BrowserNS),
		podClient:         clientset.CoreV1().Pods(envConfig.BrowserNS),
	}
}

func (kp *k8sWebDriverProvisioner) Available(ctx context.Context) (*browserkubev1.BrowserSetList, error) {
	browserSets, err := kp.browserSetsClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return browserSets, nil
}

func (kp *k8sWebDriverProvisioner) Provision(
	ctx context.Context,
	id string,
	opts *session.Capabilities,
) (*browserkubev1.Browser, error) {
	kp.logger.Infow("Starting browser", "opts", opts)

	capsRaw, err := json.Marshal(opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if opts.BrowserKubeOpts.Type == "" {
		opts.BrowserKubeOpts.Type = browserkubev1.TypeWebDriver
	}
	tracingContext := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, tracingContext)

	browser := &browserkubev1.Browser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      id,
			Namespace: kp.envConfig.BrowserNS,
			Labels: map[string]string{
				browserkubev1.LabelBrowserVisibility: "true",
			},
			Annotations: tracingContext,
		},
		Spec: browserkubev1.BrowserSpec{
			Platform:       opts.Platform,
			BrowserVersion: opts.BrowserVersion,
			BrowserName:    opts.BrowserName,
			Timezone:       opts.Timezone,
			Type:           opts.BrowserKubeOpts.Type,
			Caps:           json.RawMessage(capsRaw),

			EnableVNC:   opts.BrowserKubeOpts.EnableVNC,
			EnableVideo: opts.BrowserKubeOpts.EnableVideo,
			Extensions:  opts.BrowserKubeOpts.Extensions,
		},
	}

	browser, err = kp.browsersClient.Create(ctx, browser)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if browser, err = kp.waitForBrowser(ctx, browser, browserUPTimeout); err != nil {
		return browser, errors.WithStack(err)
	}

	if browser.Spec.Type == browserkubev1.TypeWebDriver {
		if err := browserkubeutil.SeleniumUP(ctx, browser.Status.SeleniumURL); err != nil {
			return browser, errors.Wrapf(err, "Timeout on waiting of pod running: %v", err)
		}
	}

	kp.logger.Debugf("Browser [%s] is UP", browser.GetName())
	return browser, nil
}

func (kp *k8sWebDriverProvisioner) Delete(ctx context.Context, id string) error {
	kp.logger.Infof("Deleting Browser [%s]", id)
	return errors.WithStack(kp.browsersClient.Delete(ctx, id, metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To(int64(podGracefulShutdownTimeout)),
	}))
}

func (kp *k8sWebDriverProvisioner) Logs(ctx context.Context, name string, follow bool) (io.ReadCloser, error) {
	req := kp.podClient.GetLogs(name, &apiv1.PodLogOptions{
		Container:  browserContainerName,
		Follow:     follow,
		Previous:   false,
		Timestamps: true,
	})
	logs, err := req.Stream(ctx)
	return logs, errors.WithStack(err)
}

// Update sends a MergeType PATCH request to k8s API to update our running browserset resource
func (kp *k8sWebDriverProvisioner) Update(ctx context.Context, bs *browserkubev1.BrowserSet) error {
	b, err := json.Marshal(bs)
	if err != nil {
		return err
	}
	err = kp.browserSetsClient.Patch(ctx, bs.Name, b, metav1.PatchOptions{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (kp *k8sWebDriverProvisioner) waitForBrowser(ctx context.Context, browser *browserkubev1.Browser, timeout time.Duration) (*browserkubev1.Browser, error) {
	pWatch, err := kp.browsersClient.WatchByName(ctx, browser.Name)
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
			case browserkubev1.PhasePending, "":
				// no status yet
				continue
			case browserkubev1.PhaseFailed:
				if p.Status.Reason != "" {
					return nil, &CreationErr{error: errors.New(string(p.Status.Reason))}
				}
				return nil, fmt.Errorf("browser can't be created [%s][%s]", p.Status.Phase, p.Status.Reason)
			case browserkubev1.PhaseTerminated:
				return nil, fmt.Errorf("browser has been terminated already [%s]", p.Status.Phase)
			default:
				return nil, errors.Errorf("unknown browser state: %s", p.Status.Phase)
			}
		}
	}
}
