package controller

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	browserkubeapiv1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/operator/internal/controller/browserimage"
)

func buildContainerPort(name, port string) apiv1.ContainerPort {
	//nolint:gosec
	p, _ := strconv.Atoi(port)
	return apiv1.ContainerPort{
		Name:          name,
		Protocol:      apiv1.ProtocolTCP,
		ContainerPort: int32(p),
	}
}

func buildVolumeMounts(imageType browserimage.ImageType) []apiv1.VolumeMount {
	mounts := []apiv1.VolumeMount{
		{Name: "dshm", MountPath: "/dev/shm"},
		{Name: "usergroup", MountPath: "/etc/passwd", SubPath: "passwd"},
		{Name: "usergroup", MountPath: "/etc/group", SubPath: "group"},
		{Name: "videos", MountPath: filepath.Join(imageType.Homedir(), recorderVideosRelativePath)},
		{Name: "tmp", MountPath: "/tmp"},
		{Name: "userhome", MountPath: imageType.Homedir()}, // used by selenoid images
	}

	return mounts
}

func copySpec(p *apiv1.Pod, spec *browserkubeapiv1.BrowserPodSpec) {
	if spec == nil {
		return
	}
	p.Spec.NodeSelector = spec.NodeSelector
	p.Spec.NodeName = spec.NodeName
	p.Spec.DNSPolicy = spec.DNSPolicy
	p.Spec.Affinity = spec.Affinity
	p.Spec.HostAliases = spec.HostAliases
	p.Spec.Tolerations = spec.Tolerations
	p.Spec.PriorityClassName = spec.PriorityClassName
	p.Spec.Priority = spec.Priority
	p.Spec.ActiveDeadlineSeconds = spec.ActiveDeadlineSeconds
	p.Spec.TerminationGracePeriodSeconds = spec.TerminationGracePeriodSeconds
	p.Spec.ServiceAccountName = spec.ServiceAccountName
}

func buildVolumes(opts *BrowserCtrlOpts) []apiv1.Volume {
	return []apiv1.Volume{
		{
			Name:         "userhome",
			VolumeSource: apiv1.VolumeSource{EmptyDir: &apiv1.EmptyDirVolumeSource{}},
		},
		{
			Name: "usergroup",
			VolumeSource: apiv1.VolumeSource{
				ConfigMap: &apiv1.ConfigMapVolumeSource{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: opts.browserUserConfig,
					},
					Items: []apiv1.KeyToPath{
						{Key: "group", Path: "group"},
						{Key: "passwd", Path: "passwd"},
					},
				},
			},
		},
		{
			Name: "dshm",
			VolumeSource: apiv1.VolumeSource{EmptyDir: &apiv1.EmptyDirVolumeSource{
				Medium:    apiv1.StorageMediumMemory,
				SizeLimit: resource.NewQuantity(size1Gb, resource.BinarySI), // 1GB
			}},
		},
		{
			Name:         "videos",
			VolumeSource: apiv1.VolumeSource{EmptyDir: &apiv1.EmptyDirVolumeSource{}},
		},
		{
			Name:         "tmp",
			VolumeSource: apiv1.VolumeSource{EmptyDir: &apiv1.EmptyDirVolumeSource{}},
		},
		{
			Name:         "plugins",
			VolumeSource: apiv1.VolumeSource{EmptyDir: &apiv1.EmptyDirVolumeSource{}},
		},
	}
}

// CPU, in cores. (500m = .5 cores)
// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
func buildResources(
	limitResourceCPU,
	limitResourceMemory,
	requestResourceCPU,
	requestResourceMemory int64,
) apiv1.ResourceRequirements {
	return apiv1.ResourceRequirements{
		Limits: apiv1.ResourceList{
			apiv1.ResourceCPU:    *resource.NewMilliQuantity(limitResourceCPU, resource.DecimalSI),
			apiv1.ResourceMemory: *resource.NewQuantity(limitResourceMemory, resource.BinarySI),
		},
		Requests: apiv1.ResourceList{
			apiv1.ResourceCPU:    *resource.NewMilliQuantity(requestResourceCPU, resource.DecimalSI),
			apiv1.ResourceMemory: *resource.NewQuantity(requestResourceMemory, resource.BinarySI),
		},
	}
}

func buildSidecarEnvVar(sidecarPort string, imageType browserimage.ImageType, browserConfigPort, browserConfigPath string) []apiv1.EnvVar {
	proxyURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("localhost", browserConfigPort),
		Path:   browserConfigPath,
	}

	return []apiv1.EnvVar{
		{Name: "PORT", Value: sidecarPort},
		{Name: "PROXY_URL", Value: proxyURL.String()},
		{Name: "BROWSER_HOME_DIR", Value: imageType.Homedir()},
	}
}

func GetResolution(res string) string {
	if res == "" {
		return "1920x1080x24"
	}
	resParts := strings.Split(res, "x")
	if len(resParts) == 2 {
		resParts = append(resParts, "24")
	}
	return strings.Join(resParts, "x")
}

func addContainerRecorder(
	opts *BrowserCtrlOpts,
	spec *apiv1.PodSpec,
	imageType browserimage.ImageType,
	displayNum string,
	volumeMounts []apiv1.VolumeMount,
) {
	recorderURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("localhost", "5555"),
	}

	spec.Containers[0].Env = append(spec.Containers[0].Env, apiv1.EnvVar{
		Name: "RECORDER_URL", Value: recorderURL.String(),
	})
	spec.Containers = append(spec.Containers,
		apiv1.Container{
			Name:  containerNameRecorder,
			Image: opts.recorderImage,
			SecurityContext: &apiv1.SecurityContext{
				AllowPrivilegeEscalation: ptr.To(false),
				RunAsNonRoot:             ptr.To(true),
			},
			Ports: []apiv1.ContainerPort{
				buildContainerPort("http", "5555"),
			},
			Args: []string{
				"--video-size=1360x1020",
				"--frame-rate=12",
				"--display-num=" + displayNum,
				"--codec=libx264",
				"--file-path=" + filepath.Join(imageType.Homedir(), "/videos"),
			},
			VolumeMounts: volumeMounts, // re-use volume mounts
		})
}

func getBrowserPodName(n string) string {
	return fmt.Sprintf("browser-%s", strings.ToLower(n))
}

func getBrowserPodLabels(sessionID string) map[string]string {
	return map[string]string{
		browserkubeapiv1.LabelComponent: browserkubeapiv1.LabelValueComponentBrowserkubeBrowser,
		browserkubeapiv1.LabelApp:       browserkubeapiv1.LabelValueApplicationBrowserkube,
		browserkubeapiv1.LabelSessionID: sessionID,
	}
}

func installPlugins(
	spec *apiv1.PodSpec,
	browserName string,
	imageType browserimage.ImageType,
	extensions []browserkubeapiv1.BrowserExtension,
	extensionInstallerImage,
	browserExtensionConfig string,
) {
	extMounts := buildExtensionMounts(imageType)[browserName]

	for i, c := range spec.Containers {
		if c.Name == containerNameBrowser {
			spec.Containers[i].VolumeMounts = append(c.VolumeMounts, extMounts...)
			break
		}
	}

	args := make([]string, len(extensions))
	for _, extension := range extensions {
		args = append(args,
			"--browserName="+browserName,
			"--extensionId="+extension.ExtensionID,
			"--updateUrl="+extension.UpdateURL,
		)
	}
	spec.InitContainers = []apiv1.Container{
		{
			Name:  extensionInstallerContainerName,
			Image: extensionInstallerImage,
			Env: []apiv1.EnvVar{
				{
					Name: fmt.Sprintf("%s-WHITELIST-%s", strings.ToUpper(browserExtensionConfig), strings.ToUpper(browserName)),
					ValueFrom: &apiv1.EnvVarSource{
						ConfigMapKeyRef: &apiv1.ConfigMapKeySelector{
							LocalObjectReference: apiv1.LocalObjectReference{
								Name: browserExtensionConfig,
							},
							Key: fmt.Sprintf("whitelist.%s", strings.ToLower(browserName)),
						},
					},
				},
			},
			Args:         args,
			VolumeMounts: extMounts,
		},
	}
}

func buildExtensionMounts(imageType browserimage.ImageType) map[string][]apiv1.VolumeMount {
	return map[string][]apiv1.VolumeMount{
		"firefox": {
			{Name: "plugins", MountPath: filepath.Join(imageType.Homedir(), "/.mozilla/extensions")},
			{Name: "plugins", MountPath: "/opt/firefox"},
		},
		"chrome": {
			{Name: "plugins", MountPath: "/opt/google/chrome/extensions"},
		},
	}
}
