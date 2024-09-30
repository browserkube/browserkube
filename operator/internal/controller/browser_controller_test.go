package controller

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	browserkubeapiv1 "github.com/browserkube/browserkube/operator/api/v1"
)

const (
	timeout           = time.Second * 40
	interval          = time.Millisecond * 250
	browserNamePrefix = "test"
	defaultNs         = "default"
	browserName       = "chrome"
)

var _ = Describe("Browser controller", Ordered, func() {
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)

	Context("When finding browser config", func() {
		It("find proper config", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}
			createBrowserSet(defaultNs, "check-browser-config-test")
			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName:    browserName,
					BrowserVersion: "version",
				},
			})
			Expect(err).Should(Succeed())
			Expect(config).Should(Not(BeNil()))
		})

		It("browser version is not supported", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName:    browserName,
					BrowserVersion: "non-existing-version",
				},
			})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("browser 'non-existing-version' is not supported"))
		})

		It("switch browser type to playwright", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName:    browserName,
					BrowserVersion: "version",
					Type:           browserkubeapiv1.TypePlaywright,
				},
			})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("browser 'chrome' is not supported"))
		})

		It("unknown browser type", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName:    browserName,
					BrowserVersion: "version",
					Type:           "unknown",
				},
			})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(&browserErr{reason: browserkubeapiv1.ReasonUnknownSessionType}))
		})

		It("browser 'firefox' is not supported", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName:    "firefox",
					BrowserVersion: "version",
				},
			})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("browser 'firefox' is not supported"))
		})

		It("browser platform not a linux", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName: browserName,
					Platform:    "default",
				},
			})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(&browserErr{reason: browserkubeapiv1.ReasonPlatformNotSupported}))
		})

		It("browser is not provided", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("browser is not provided"))
		})

		It("browser's config isn't found. Namespace not default", func() {
			bc := &BrowserReconciler{
				Client: k8sClient,
			}

			config, err := bc.findBrowserConfig(context.Background(), &browserkubeapiv1.Browser{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "notDefault",
				},
				Spec: browserkubeapiv1.BrowserSpec{
					BrowserName:    browserName,
					BrowserVersion: "non-existing-version",
				},
			})
			Expect(config).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(&browserErr{reason: browserkubeapiv1.ReasonConfigNotFound}))
		})
	})

	Context("When creating new browser", func() {
		It("create readiness configuration", func() {
			if err := k8sClient.Create(ctx, &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
					Name:      probeConfigMapName,
				},
				Data: map[string]string{"webdriver.enabled": "false"},
			}); err != nil {
				Expect(err).ToNot(HaveOccurred(), "failed to create readiness config map")
			}
		})

		It("check configuration variables", func() {
			if err := k8sClient.Update(ctx, &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: defaultNs,
					Name:      probeConfigMapName,
				},
				Data: map[string]string{
					"webdriver.enabled": "true",
				},
			}); err != nil {
				Expect(err).ToNot(HaveOccurred(), "failed to create readiness config map")
			}
		})

		It("creates browser configuration", func() {
			createBrowserSet(defaultNs, "test-config")
		})

		It("changes the status when browser pod is created", func() {
			browserName := fmt.Sprintf("%s-%d", browserNamePrefix, rnd.Intn(1000))
			b := createBrowser(browserName)
			podName := fmt.Sprintf("%s-%s", containerNameBrowser, browserName)
			By("Checking that the Browser has one active pod", func() {
				p := waitForBrowserPod(b.Namespace, podName)
				setPodContainersReady(p)
				setBrowserPodStatus(p, v1.PodRunning)
				waitForBrowserStatus(b.Namespace, b.Name, browserkubeapiv1.PhaseRunning)
			})
		})
		It("adds VNC to selenoid images if enabled", func() {
			browserName := fmt.Sprintf("%s-%d", browserNamePrefix, rnd.Intn(1000))
			b := createBrowser(browserName)
			podName := fmt.Sprintf("%s-%s", containerNameBrowser, browserName)
			By("Browser image has appropriate env variables", func() {
				p := waitForBrowserPod(b.Namespace, podName)
				container, found := find(p.Spec.Containers, func(item v1.Container) bool {
					return item.Name == containerNameBrowser
				})

				Expect(found).Should(BeTrue(), "browser container isn't found")

				vncEnabled, envFound := find(container.Env, func(item v1.EnvVar) bool {
					return item.Name == "ENABLE_VNC"
				})
				Expect(envFound).Should(BeTrue(), "ENABLE_VNC isn't found")
				Expect(vncEnabled.Value).Should(Equal("true"), "ENABLE_VNC isn't true")
			})
		})

		It("check playwright browser creation", func() {
			browserName := fmt.Sprintf("%s-%d", browserNamePrefix, rnd.Intn(1000))
			createBrowserPlaywright(defaultNs, browserName)
		})

		It("check browser without type creation", func() {
			browserName := fmt.Sprintf("%s-%d", browserNamePrefix, rnd.Intn(1000))
			createBrowserWithoutType(defaultNs, browserName)
		})

		It("create new browser for testing sidecar container and finalizer", func() {
			browserName := fmt.Sprintf("%s-%d", browserNamePrefix, rnd.Intn(1000))
			b := createBrowser(browserName)

			By("Set status for sidecar container", func() {
				podName := fmt.Sprintf("%s-%s", containerNameBrowser, browserName)
				p := waitForBrowserPod(b.Namespace, podName)
				setPodContainersReady(p)
				setBrowserPodContainerStatuses(p, v1.PodRunning)
			})

			By("Removing the test browser", func() {
				podName := fmt.Sprintf("%s-%s", containerNameBrowser, browserName)
				_ = waitForBrowserPod(b.Namespace, podName)

				err := k8sClient.Delete(ctx, b)
				Expect(err).Should(Succeed())
			})
		})
		It("deletes browser when sidecar is terminated", func() {
			browserName := fmt.Sprintf("%s-%s", browserNamePrefix, "sidecartest")
			b := createBrowser(browserName)

			podName := fmt.Sprintf("%s-%s", containerNameBrowser, browserName)
			p := waitForBrowserPod(b.Namespace, podName)
			setPodContainersReady(p)
			setBrowserPodStatus(p, v1.PodRunning)
			b = waitForBrowserStatus(b.Namespace, b.Name, browserkubeapiv1.PhaseRunning)

			helper, err := patch.NewHelper(p, k8sClient)
			Expect(err).Should(Succeed(), "patch ok")
			p.Status.ContainerStatuses = append(p.Status.ContainerStatuses, v1.ContainerStatus{
				Name: containerNameSidecar,
				State: v1.ContainerState{
					Terminated: &v1.ContainerStateTerminated{
						ExitCode:   -1,
						FinishedAt: metav1.Now(),
					},
				},
			})
			Expect(helper.Patch(ctx, p)).Should(Succeed(), "statuses ok")

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: b.Name, Namespace: b.Namespace}, &browserkubeapiv1.Browser{})
				return errors.IsNotFound(err)
			}).Should(BeTrue())
		})
		It("can be terminated", func() {
			browserName := fmt.Sprintf("%s-%d", browserNamePrefix, rnd.Intn(1000))
			b := createBrowser(browserName)

			By("Checking that the Browser has one active pod")

			podName := fmt.Sprintf("%s-%s", containerNameBrowser, browserName)
			p := waitForBrowserPod(b.Namespace, podName)
			setPodContainersReady(p)
			setBrowserPodStatus(p, v1.PodRunning)
			b = waitForBrowserStatus(b.Namespace, b.Name, browserkubeapiv1.PhaseRunning)

			helper, err := patch.NewHelper(b, k8sClient)
			Expect(err).Should(Succeed())
			b.Status.Phase = browserkubeapiv1.PhaseTerminated
			Expect(helper.Patch(ctx, b)).Should(Succeed())

			Eventually(func() bool {
				p := &v1.Pod{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: b.Name, Namespace: b.Namespace}, p)
				return errors.IsNotFound(err)
			}).Should(BeTrue())
		})
	})
})

func createBrowserSet(ns, name string) {
	bConfig := &browserkubeapiv1.BrowserSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: browserkubeapiv1.BrowserSetSpec{
			WebDriver: map[string]browserkubeapiv1.BrowsersConfig{
				browserName: {
					DefaultPath:    "/",
					DefaultVersion: "version",
					Versions: map[string]browserkubeapiv1.BrowserConfig{
						"version": {
							Image:    "selenoid/vnc_chrome:103.0",
							Provider: "k8s",
							Port:     "4444",
						},
					},
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, bConfig)).Should(Succeed())
}

func createBrowser(name string) *browserkubeapiv1.Browser {
	b := browserkubeapiv1.Browser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: defaultNs,
		},
		Spec: browserkubeapiv1.BrowserSpec{
			BrowserName: browserName,
			Type:        browserkubeapiv1.TypeWebDriver,
			EnableVideo: true,
			EnableVNC:   true,
			Extensions: []browserkubeapiv1.BrowserExtension{
				{
					ExtensionID: "extID",
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, &b)).Should(Succeed())

	return &b
}

func createBrowserPlaywright(ns, name string) *browserkubeapiv1.Browser {
	b := browserkubeapiv1.Browser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: browserkubeapiv1.BrowserSpec{
			BrowserName: browserName,
			Type:        browserkubeapiv1.TypePlaywright,
			EnableVideo: true,
			Extensions: []browserkubeapiv1.BrowserExtension{
				{
					ExtensionID: "extID",
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, &b)).Should(Succeed())

	return &b
}

func createBrowserWithoutType(ns, name string) *browserkubeapiv1.Browser {
	b := browserkubeapiv1.Browser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: browserkubeapiv1.BrowserSpec{
			BrowserName: browserName,
			Type:        "",
			EnableVideo: true,
			Extensions: []browserkubeapiv1.BrowserExtension{
				{
					ExtensionID: "extID",
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, &b)).Should(Succeed())

	return &b
}

func waitForBrowserPod(ns, name string) *v1.Pod {
	p := &v1.Pod{}
	Eventually(func() error {
		err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, p)
		if err != nil {
			return err
		}
		return nil
	}, timeout, interval).Should(Succeed(), "should get browser pod  %s ", name)
	return p
}

func waitForBrowserStatus(ns, name string, status browserkubeapiv1.Phase) *browserkubeapiv1.Browser {
	browserLookupKey := types.NamespacedName{Name: name, Namespace: ns}
	createdBrowser := &browserkubeapiv1.Browser{}
	Eventually(func() (browserkubeapiv1.Phase, error) {
		err := k8sClient.Get(ctx, browserLookupKey, createdBrowser)
		if err != nil {
			return browserkubeapiv1.PhasePending, err
		}

		return createdBrowser.Status.Phase, nil
	}, timeout, interval).Should(Equal(status), "should list browser %s in the active jobs list in status", name)
	return createdBrowser
}

func setBrowserPodStatus(p *v1.Pod, phase v1.PodPhase) {
	helper, err := patch.NewHelper(p, k8sClient)
	Expect(err).Should(Succeed(), "patch ok")
	p.Status.Phase = phase
	p.Status.PodIP = "127.0.0.1"
	Expect(helper.Patch(ctx, p)).Should(Succeed(), "phase ok")
}

func setPodContainersReady(p *v1.Pod) {
	helper, err := patch.NewHelper(p, k8sClient)
	Expect(err).Should(Succeed(), "patch ok")
	for i := range p.Status.ContainerStatuses {
		p.Status.ContainerStatuses[i].Ready = true
	}
	Expect(helper.Patch(ctx, p)).Should(Succeed(), "statuses ok")
}

func setBrowserPodContainerStatuses(p *v1.Pod, phase v1.PodPhase) {
	helper, err := patch.NewHelper(p, k8sClient)
	Expect(err).Should(Succeed(), "patch ok")
	p.Status.Phase = phase
	p.Status.ContainerStatuses = []v1.ContainerStatus{
		{
			Name: containerNameSidecar,
			State: v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{},
			},
		},
		{
			Name:  containerNameBrowser,
			State: v1.ContainerState{},
		},
	}
	Expect(helper.Patch(ctx, p)).Should(Succeed(), "phase ok")
}

func find[T any](collection []T, predicate func(item T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}

	var result T
	return result, false
}

func TestNewBrowserReconciler(t *testing.T) {
	type args struct {
		client client.Client
		scheme *runtime.Scheme
		opts   *BrowserCtrlOpts
	}
	tests := []struct {
		name string
		args args
		want *BrowserReconciler
	}{
		{
			name: "NewBrowserReconciler",
			args: args{
				client: k8sClient,
				scheme: &runtime.Scheme{},
			},
			want: &BrowserReconciler{
				Client: k8sClient,
				Scheme: &runtime.Scheme{},
				finalizerName: fmt.Sprintf("%s/finalizer",
					browserkubeapiv1.GroupVersion.WithResource("browsers").GroupResource().String()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBrowserReconciler(tt.args.client, tt.args.scheme, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBrowserReconciler() = %v, want %v", got, tt.want)
			}
		})
	}
}
