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

type selenoidPodBuilder struct {
	browserConfig *browserkubeapiv1.BrowserConfig
}

func (s *selenoidPodBuilder) Build(ctx context.Context, b *browserkubeapiv1.Browser, opts *BrowserCtrlOpts, readinessProbe *apiv1.Probe) (*apiv1.Pod, error) {
	volumeMounts := buildVolumeMounts(browserimage.ImageTypeSelenoid)

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
				Env:          buildSidecarEnvVar(opts.sidecarPort, browserimage.ImageTypeSelenoid, s.browserConfig.Port, s.browserConfig.Path),
				VolumeMounts: volumeMounts,
				Resources:    buildResources(200, memory128Mi, 100, memory128Mi),
			},
			{
				Name:  containerNameBrowser,
				Image: s.browserConfig.Image,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("browser", s.browserConfig.Port),
					buildContainerPort("vnc", ports.VNC),
					//buildContainerPort("devtools", ports.DevTools),
				},
				Env:            s.buildBrowserEnvVars(b, s.browserConfig, b.Spec.EnableVNC),
				VolumeMounts:   volumeMounts,
				ReadinessProbe: readinessProbe,
				Resources:      buildResources(1000, size2Gi, 500, size2Gi),
			},
			{
				Name:         containerNameClipboard,
				Image:        opts.clipboardImage,
				VolumeMounts: volumeMounts,
				Ports: []apiv1.ContainerPort{
					buildContainerPort("p", ports.Clipboard),
				},
				Env:   []apiv1.EnvVar{{Name: "DISPLAY", Value: ":99"}},
				Stdin: true,
				TTY:   true,
			},
		},
		Volumes: buildVolumes(opts),
	}

	if s.browserConfig.EnableVideo {
		addContainerRecorder(opts, spec, browserimage.ImageTypeSelenoid, "99", volumeMounts)
	}

	logger := log.FromContext(ctx)

	if b.Spec.Extensions != nil || len(b.Spec.Extensions) != 0 {
		logger.Info("Browser Extension Capabilities: ", "Capabilities", fmt.Sprintf("%+v", b.Spec.Extensions))
		installPlugins(spec,
			b.Spec.BrowserName,
			browserimage.ImageTypeSelenoid,
			b.Spec.Extensions,
			opts.extensionInstallerImage,
			opts.browserExtensionConfig,
		)
	}

	//if s.browserConfig.RegistrySecret != "" {
	//	spec.ImagePullSecrets = []apiv1.LocalObjectReference{
	//		{Name: s.browserConfig.RegistrySecret},
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
	copySpec(browserPod, s.browserConfig.Spec)

	return browserPod, nil
}

func (s *selenoidPodBuilder) buildBrowserEnvVars(b *browserkubeapiv1.Browser, browserConfig *browserkubeapiv1.BrowserConfig, browserEnableVNC bool) []apiv1.EnvVar {
	vars := []apiv1.EnvVar{
		{Name: "TZ", Value: browserConfig.Timezone},
		{Name: "DISPLAY", Value: ":99"},
		{Name: "SCREEN_RESOLUTION", Value: GetResolution(b.Spec.ScreenResolution)},
	}

	if browserEnableVNC {
		vars = append(vars, apiv1.EnvVar{Name: "ENABLE_VNC", Value: "true"})
		//vars = append(vars, apiv1.EnvVar{Name: "ENABLE_VIDEO", Value: "false"})
	}

	return vars
}
