/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"

	errors2 "github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	browserkubeapiv1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/operator/internal/controller/utils"
)

const (
	size1Gb     = 1 * 1024 * 1024 * 1024
	size2Gi     = 2 * size1Gb
	memory128Mi = 128 * 1024 * 1024
)

const (
	podGracefulShutdownTimeout = 30

	// container names
	containerNameBrowser            = "browser"
	containerNameSidecar            = "sidecar"
	containerNameRecorder           = "recorder"
	containerNameClipboard          = "clipboard"
	extensionInstallerContainerName = "extension-installer"
)

// selenium constants
const sidecarSeleniumPath = "/wd/hub"

// recorder constants
const (
	recorderVideosRelativePath = "/videos"
)

var ports = browserkubeapiv1.PortConfig{
	VNC:        "5900",
	DevTools:   "7070",
	FileServer: "8080",
	Clipboard:  "9191",
	Sidecar:    "9999",
	Browser:    "4444",
}

type browserErr struct {
	error
	reason browserkubeapiv1.Reason
}

func (br *browserErr) Error() string {
	if br.error != nil {
		return br.error.Error()
	}
	return ""
}

// BrowserReconciler reconciles a Browser object
type BrowserReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	finalizerName string
	opts          *BrowserCtrlOpts
}

func NewBrowserReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	opts *BrowserCtrlOpts,
) *BrowserReconciler {
	return &BrowserReconciler{
		Client: client,
		Scheme: scheme,
		opts:   opts,
		finalizerName: fmt.Sprintf("%s/finalizer",
			browserkubeapiv1.GroupVersion.WithResource("browsers").GroupResource().String()),
	}
}

//+kubebuilder:rbac:groups=api.browserkube.io,namespace=browserkube,resources=browsersets,verbs=get;list;watch
//+kubebuilder:rbac:groups=api.browserkube.io,namespace=browserkube,resources=browsers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.browserkube.io,namespace=browserkube,resources=browsers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.browserkube.io,namespace=browserkube,resources=browsers/finalizers,verbs=update
//+kubebuilder:rbac:groups=api.browserkube.io,namespace=browserkube,resources=sessionresults,verbs=get;list;watch
//+kubebuilder:rbac:groups=api.browserkube.io,namespace=browserkube,resources=sessionresults,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,namespace=browserkube,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,namespace=browserkube,resources=pods/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *BrowserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	instance := &browserkubeapiv1.Browser{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		logger.Error(err, "unable to get browser")
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if res, err := r.checkFinalizer(ctx, instance); res != nil {
		return *res, err
	}

	var browserkubePod apiv1.Pod
	err = r.Get(context.TODO(), types.NamespacedName{
		Namespace: req.NamespacedName.Namespace,
		Name:      getBrowserPodName(req.NamespacedName.Name),
	}, &browserkubePod)
	if err != nil {
		if errors.IsNotFound(err) {
			if !instance.DeletionTimestamp.IsZero() {
				// already terminated
				return ctrl.Result{}, nil
			}
			logger.Info("creating pod for browser", "name", browserkubePod.Name)

			// create the browser
			if cErr := r.createBrowser(ctx, instance); cErr != nil {
				var bErr *browserErr
				if errors2.As(cErr, &bErr) {
					instance.Status.Reason = bErr.reason
				} else {
					instance.Status.Reason = browserkubeapiv1.ReasonUnknown
				}
				instance.Status.Message = cErr.Error()
				instance.Status.Phase = browserkubeapiv1.PhaseFailed

				if suErr := r.Status().Update(ctx, instance); suErr != nil {
					return ctrl.Result{Requeue: false}, suErr
				}
				return ctrl.Result{}, cErr
			}
			// pod is created. requeue to wait until it's running
			return reconcile.Result{Requeue: true}, nil
		}
		logger.Error(err, "unable to get browser pod")
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if res, resErr := r.checkPending(ctx, instance, &browserkubePod); res != nil {
		return *res, resErr
	}
	if res, resErr := r.checkTerminated(ctx, instance, &browserkubePod); res != nil {
		return *res, resErr
	}
	if res, resErr := r.checkSidecarRunning(ctx, instance, &browserkubePod); res != nil {
		return *res, resErr
	}

	return ctrl.Result{}, nil
}

func (r *BrowserReconciler) createBrowser(
	ctx context.Context,
	browser *browserkubeapiv1.Browser,
) error {
	logger := log.FromContext(ctx)

	browserConfig, err := r.findBrowserConfig(ctx, browser)
	if err != nil {
		logger.Error(err, fmt.Sprintf("error while finding config: %s", err))
		return err
	}
	if uErr := r.Update(ctx, browser); uErr != nil {
		logger.Error(err, "error while updating browser")
		return uErr
	}

	logger.Info("Starting browser pod", "image", browserConfig.Image)

	readinessProbe, err := r.getReadinessProbe(ctx, browser.Namespace, browser.Spec.Type, browserConfig.Path, browserConfig.Port)
	if err != nil {
		logger.Error(err, fmt.Sprintf("error while getting readiness probe: %s", err))
	}

	browserPodBuilder, imgType, err := NewPodBuilder(browserConfig)
	if err != nil {
		logger.Error(err, "error while creating browser pod builder", "error", err.Error())
		return err
	}

	browserPod, err := browserPodBuilder.Build(ctx, browser, r.opts, readinessProbe)
	if err != nil {
		logger.Error(err, "error while creating browser pod", "error", err.Error())
		return err
	}

	if err = controllerutil.SetControllerReference(browser, browserPod, r.Scheme); err != nil {
		logger.Error(err, "error while setting controller reference")
		return err
	}

	err = r.Create(ctx, browserPod, &client.CreateOptions{})
	if err != nil {
		logger.Error(err, "error while creating browser object", "error", err.Error())
		return err
	}

	browser.Status.PortConfig = ports
	browser.Status.Phase = browserkubeapiv1.PhasePending
	browser.Status.Image = browserConfig.Image
	browser.Status.PodName = browserPod.Name
	browser.Status.VncPass = imgType.VncPass()
	logger.Info("updating browser resource", "resource", fmt.Sprintf("%+v", browser))

	if err := r.Status().Update(ctx, browser); err != nil {
		logger.Error(err, "error while updating browser")
		return err
	}
	logger.Info("resource updated to new status pending", "resource", browser.Name)
	return nil
}

func (r *BrowserReconciler) getReadinessProbe(ctx context.Context, namespace, browserType, path, port string) (*apiv1.Probe, error) {
	var livenessConfigMap apiv1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      r.opts.browserReadinessConfig,
	},
		&livenessConfigMap,
	)
	if err != nil {
		return nil, err
	}

	prefix := strings.ToLower(browserType)

	enabled := livenessConfigMap.Data[utils.CreateReadinessConfigmapKey(prefix, "enabled")]
	if enabled == "false" {
		return nil, nil
	}

	probeAction := r.getReadinessProbeAction(browserType, path, port)
	if probeAction == nil {
		return nil, nil
	}

	initialDelay, err := strconv.Atoi(livenessConfigMap.Data[utils.CreateReadinessConfigmapKey(prefix, "initialDelaySeconds")])
	if err != nil {
		initialDelay = 2
	}
	timeoutSecond, err := strconv.Atoi(livenessConfigMap.Data[utils.CreateReadinessConfigmapKey(prefix, "timeoutSecond")])
	if err != nil {
		timeoutSecond = 10
	}
	failureThreshold, err := strconv.Atoi(livenessConfigMap.Data[utils.CreateReadinessConfigmapKey(prefix, "failureThreshold")])
	if err != nil {
		failureThreshold = 10
	}

	return &apiv1.Probe{
		ProbeHandler: apiv1.ProbeHandler{
			HTTPGet: probeAction,
		},
		InitialDelaySeconds: int32(initialDelay),
		TimeoutSeconds:      int32(timeoutSecond),
		FailureThreshold:    int32(failureThreshold),
	}, nil
}

func (r *BrowserReconciler) findBrowserConfig(ctx context.Context, browser *browserkubeapiv1.Browser) (*browserkubeapiv1.BrowserConfig, error) {
	instances := &browserkubeapiv1.BrowserSetList{}
	err := r.List(ctx, instances, client.InNamespace(browser.Namespace))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, &browserErr{reason: browserkubeapiv1.ReasonConfigNotFound}
		}
		return nil, fmt.Errorf("browsers config can't be loaded: %w", err)
	}
	if len(instances.Items) == 0 {
		return nil, &browserErr{reason: browserkubeapiv1.ReasonConfigNotFound}
	}
	instance := instances.Items[0]

	var platform string
	if browser.Spec.Platform == "" {
		platform = "linux"
	} else {
		platform = strings.ToLower(browser.Spec.Platform)
	}

	// TODO validate
	browserName := strings.ToLower(browser.Spec.BrowserName)
	if browserName == "" {
		return nil, errors2.New("browser is not provided")
	}
	browserType := browser.Spec.Type
	if browserType == "" {
		browserType = browserkubeapiv1.TypeWebDriver
	}

	var browsers map[string]browserkubeapiv1.BrowsersConfig
	switch browserType {
	case browserkubeapiv1.TypePlaywright:
		browsers = instance.Spec.Playwright
	case browserkubeapiv1.TypeWebDriver:
		browsers = instance.Spec.WebDriver
	default:
		return nil, &browserErr{reason: browserkubeapiv1.ReasonUnknownSessionType}
	}
	switch platform {
	case "linux":
		browserMapping, ok := browsers[browserName]
		if !ok {
			return nil, &browserErr{reason: browserkubeapiv1.ReasonVersionNotSupported, error: fmt.Errorf("browser '%s' is not supported", browserName)}
		}

		browser.Spec.BrowserVersion = utils.FirstNonEmpty(browser.Spec.BrowserVersion, browserMapping.DefaultVersion)
		browserConfig, ok := browserMapping.Versions[browser.Spec.BrowserVersion]
		if !ok {
			return nil, &browserErr{reason: browserkubeapiv1.ReasonVersionNotSupported, error: fmt.Errorf("browser '%s' is not supported", browser.Spec.BrowserVersion)}
		}
		browserConfig.Path = utils.FirstNonEmpty(browserConfig.Path, browserMapping.DefaultPath)
		browserConfig.Timezone = utils.FirstNonEmpty(browser.Spec.Timezone, browserConfig.Timezone, instance.Spec.DefaultTimezone, "UTC")

		// video options
		browserConfig.EnableVideo = browser.Spec.EnableVideo

		// TODO write some unit tests
		if podSpec := instance.Spec.PodSpec; podSpec != nil {
			if browserConfig.Spec == nil {
				browserConfig.Spec = podSpec
			} else {
				merr := mergo.Merge(&browserConfig.Spec, podSpec)
				if merr != nil {
					return nil, merr
				}
			}
		}
		return &browserConfig, nil

	default:
		return nil, &browserErr{reason: browserkubeapiv1.ReasonPlatformNotSupported}
	}
}

func (r *BrowserReconciler) deletePod(ctx context.Context, pod *apiv1.Pod) error {
	logger := log.FromContext(ctx)

	logger.Info("Deleting Browser pod", "name", pod.Name)
	return r.Delete(ctx, pod, &client.DeleteOptions{
		GracePeriodSeconds: ptr.To(int64(podGracefulShutdownTimeout)),
	})
}

// nolint:unparam
func (r *BrowserReconciler) checkTerminated(ctx context.Context, instance *browserkubeapiv1.Browser, browserkubePod *apiv1.Pod) (*ctrl.Result, error) {
	if instance.Status.Phase == browserkubeapiv1.PhaseTerminated {
		if err := r.deletePod(context.Background(), browserkubePod); err != nil {
			log.FromContext(ctx).Error(err, "unable to delete pod browser")
		}
	}
	return nil, nil
}

func (r *BrowserReconciler) checkPending(ctx context.Context, instance *browserkubeapiv1.Browser, browserkubePod *apiv1.Pod) (*ctrl.Result, error) {
	// status is Pending. Awaiting pod creation
	// once pod containers are ready, changing the status accordingly
	if instance.Status.Phase != browserkubeapiv1.PhasePending {
		return nil, nil
	}

	var readyCount int
	for _, status := range browserkubePod.Status.ContainerStatuses {
		if status.Ready {
			readyCount++
		}
	}

	host := browserkubePod.Status.PodIP

	if readyCount == len(browserkubePod.Status.ContainerStatuses) && host != "" {
		instance.Status.Phase = browserkubeapiv1.PhaseRunning
		seleniumURL := (&url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(host, ports.Sidecar),
			Path:   sidecarSeleniumPath,
		}).String()
		instance.Status.Host = host
		instance.Status.SeleniumURL = seleniumURL

		if err := r.Status().Update(ctx, instance); err != nil {
			return &ctrl.Result{}, err
		}
		return &ctrl.Result{}, nil
	}
	// cleanup resource if it can't get up and running
	if instance.GetCreationTimestamp().Add(5 * time.Minute).Before(time.Now()) {
		return &ctrl.Result{}, r.Delete(ctx, instance)
	}

	// requeue to wait for pod to get up and running
	return &ctrl.Result{Requeue: true}, nil
}

//nolint:unparam
func (r *BrowserReconciler) checkSidecarRunning(ctx context.Context, instance *browserkubeapiv1.Browser, browserkubePod *apiv1.Pod) (*ctrl.Result, error) {
	// sidecar container exited which signals that browser termination is requested
	if browserkubePod.Status.Phase == apiv1.PodRunning {
		logger := log.FromContext(ctx)
		logger.Info("checking quit session:", "container status", browserkubePod.Status.ContainerStatuses)
		for _, c := range browserkubePod.Status.ContainerStatuses {
			if c.Name == containerNameSidecar && c.State.Terminated != nil {
				logger.Info("Browser seems to be timed out. Deleting...")
				instance.Status.Phase = browserkubeapiv1.PhaseTerminated
				if err := r.Status().Update(ctx, instance); err != nil {
					logger.Error(err, "unable to update browser status")
				}
				if instance.DeletionTimestamp.IsZero() {
					if err := r.Delete(ctx, instance); err != nil {
						logger.Error(err, "unable to delete timed out browser")
					}
				}
				logger.Info("Browser is scheduled for deletion")
				return &ctrl.Result{Requeue: true}, nil
			}
		}
	}
	return nil, nil
}

func (r *BrowserReconciler) checkFinalizer(ctx context.Context, instance *browserkubeapiv1.Browser) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithValues("instance", instance.Name)

	// examine DeletionTimestamp to determine if object is under deletion
	// The object is not being deleted, so if it does not have our finalizer,
	// then lets add the finalizer and update the object. This is equivalent
	// registering our finalizer.
	if !controllerutil.ContainsFinalizer(instance, r.finalizerName) {
		if ok := controllerutil.AddFinalizer(instance, r.finalizerName); !ok {
			logger.Info("Failed to add finalizer into the custom resource")
			return &ctrl.Result{Requeue: true}, nil
		}
		if err := r.Update(ctx, instance); err != nil {
			logger.Error(err, "Add finalizer error", "error", err.Error())
			return &ctrl.Result{}, err
		}
	}

	if !instance.DeletionTimestamp.IsZero() {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(instance, r.finalizerName) {
			logger.Info("Executing pre-delete hook")
			// our finalizer is present, so lets handle any external dependency
			// update browser status if not set
			if err := r.preDeleteHook(ctx, instance); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return &ctrl.Result{}, err
			}
			logger.Info("Pre-delete hook has been executed")
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(instance, r.finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				logger.Error(err, "Finalizer remove error", "error", err.Error())
				return &ctrl.Result{}, err
			}
			logger.Info("Finalizer has been removed")
		}
		// Stop reconciliation as the item is being deleted
		return &ctrl.Result{}, nil
	}
	return nil, nil
}

func (r *BrowserReconciler) preDeleteHook(ctx context.Context, instance *browserkubeapiv1.Browser) error {
	// make sure the status of browser is terminated
	if instance.Status.Phase != browserkubeapiv1.PhaseTerminated {
		instance.Status.Phase = browserkubeapiv1.PhaseTerminated
		if err := r.Status().Update(ctx, instance); err != nil {
			return err
		}
	}
	return nil
}

func (r *BrowserReconciler) getReadinessProbeAction(browserType, path, port string) *apiv1.HTTPGetAction {
	switch browserType {
	case browserkubeapiv1.TypeWebDriver:
		hPath, _ := url.JoinPath(path, "/status")
		return &apiv1.HTTPGetAction{
			Scheme: apiv1.URISchemeHTTP,
			Port:   intstr.Parse(port),
			Path:   hPath,
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BrowserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&browserkubeapiv1.Browser{}).
		Owns(&apiv1.Pod{}).
		Complete(r)
}
