package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	browserkubeapiv1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/operator/internal/controller/browserimage"
)

const (
	aerokubeDisplayNum    = "0"
	aerokubeDisplay       = ":0"
	aerokubeRemoteDisplay = "127.0.0.1:0"
)

type aerokubePodBuilder struct {
	browserConfig *browserkubeapiv1.BrowserConfig
}

func (a *aerokubePodBuilder) Build(ctx context.Context, b *browserkubeapiv1.Browser, opts *BrowserCtrlOpts, readinessProbe *apiv1.Probe) (*apiv1.Pod, error) {
	volumeMounts := buildVolumeMounts(browserimage.ImageTypeAerokube)

	spec := &apiv1.PodSpec{
		Hostname:      b.Name,
		RestartPolicy: apiv1.RestartPolicyNever,
		Containers: []apiv1.Container{
			{
				Name:  containerNameSidecar,
				Image: opts.sidecarImage,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("sidecar", opts.sidecarPort),
				},
				Env:          buildSidecarEnvVar(opts.sidecarPort, browserimage.ImageTypeAerokube, a.browserConfig.Port, a.browserConfig.Path),
				VolumeMounts: volumeMounts,
				Resources:    buildResources(200, memory128Mi, 100, memory128Mi),
			},
			{
				Name:  containerNameBrowser,
				Image: a.browserConfig.Image,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("browser", a.browserConfig.Port),
					buildContainerPort("vnc", ports.VNC),
				},
				Env:          a.buildBrowserEnvVars(b.Spec, a.browserConfig),
				VolumeMounts: volumeMounts,
				// ReadinessProbe: readinessProbe,
				Resources: buildResources(1000, size2Gi, 500, size2Gi),
			},
			{
				Name:         containerNameClipboard,
				Image:        opts.clipboardImage,
				VolumeMounts: volumeMounts,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("p", ports.Clipboard),
				},
				Env:   []apiv1.EnvVar{{Name: "DISPLAY", Value: aerokubeRemoteDisplay}},
				Stdin: true,
				TTY:   true,
			},
		},
		Volumes: buildVolumes(opts),
	}

	if b.Spec.EnableVNC {
		spec.Containers = append(spec.Containers,
			apiv1.Container{
				Name:         "x-server",
				Image:        opts.xServerImage,
				VolumeMounts: volumeMounts,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("p", "6000"),
				},
				Env: []apiv1.EnvVar{
					{Name: "SCREEN_RESOLUTION", Value: GetResolution(b.Spec.ScreenResolution)},
					{Name: "DISPLAY", Value: aerokubeDisplay},
				},
			},
			apiv1.Container{
				Name:         "vnc-server",
				Image:        opts.vncServerImage,
				VolumeMounts: volumeMounts,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("p", "5900"),
				},
				// Env: []apiv1.EnvVar{{Name: "DISPLAY", Value: a.remoteDisplay()}},
			},
		)
	}
	if a.browserConfig.EnableVideo {
		addContainerRecorder(opts, spec, browserimage.ImageTypeAerokube, aerokubeDisplayNum, volumeMounts)
	}

	logger := log.FromContext(ctx)

	if b.Spec.Extensions != nil || len(b.Spec.Extensions) != 0 {
		logger.Info("Browser Extension Capabilities: ", "Capabilities", fmt.Sprintf("%+v", b.Spec.Extensions))
		installPlugins(spec,
			b.Spec.BrowserName,
			browserimage.ImageTypeAerokube,
			b.Spec.Extensions,
			opts.extensionInstallerImage,
			opts.browserExtensionConfig,
		)
	}

	//if a.browserConfig.RegistrySecret != "" {
	//	spec.ImagePullSecrets = []apiv1.LocalObjectReference{
	//		{Name: a.browserConfig.RegistrySecret},
	//	}
	//}

	browserPod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getBrowserPodName(b.Name),
			Labels:    getBrowserPodLabels(b.Name),
			Namespace: b.Namespace,
		},
		Spec: *spec,
	}
	copySpec(browserPod, a.browserConfig.Spec)
	return browserPod, nil
}

func (a *aerokubePodBuilder) buildBrowserEnvVars(browser browserkubeapiv1.BrowserSpec, browserConfig *browserkubeapiv1.BrowserConfig) []apiv1.EnvVar {
	vars := []apiv1.EnvVar{
		{Name: "TZ", Value: browserConfig.Timezone},
	}

	if browser.EnableVNC {
		vars = append(vars, apiv1.EnvVar{Name: "DISPLAY", Value: aerokubeDisplay})
	}

	return vars
}
