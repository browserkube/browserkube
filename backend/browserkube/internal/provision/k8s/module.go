package provisionk8s

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	wdsession "github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/sessionresult"
)

//
//go:generate mockery --name Indexer --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/tools/cache --output mocks
//go:generate mockery --name SharedIndexInformer --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/tools/cache --output mocks
//go:generate mockery --name Store --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/tools/cache --output mocks
//go:generate mockery --name PodInterface --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/kubernetes/typed/core/v1 --output mocks
//go:generate mockery --name ResourceQuotaNamespaceLister --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/listers/core/v1 --output mocks
//go:generate mockery --name CoreV1Interface --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/kubernetes/typed/core/v1 --output mocks
//go:generate mockery --all --dir ../../../../../operator/pkg/client/v1 --output mocks
//go:generate mockery --name Interface --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/rest --output mocks --filename RestMock.go --structname RestMock
var resyncPeriod = 1 * time.Minute

var Module = fx.Options(
	fx.Provide(
		provideClientSet,
		provideProvisioner,
		provideSessionRepository,
		provideResultsRepository,
	),
)

func provideSessionRepository(lc fx.Lifecycle,
	clientset *kubernetes.Clientset,
	browserkubeClient browserkubeclientv1.Interface,
	env *provision.Config,
) (wdsession.Repository, error) {
	sw, err := newSessionWatch(clientset, browserkubeClient, env.BrowserNS)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	watchCtx, cancelFunc := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sw.Start(watchCtx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancelFunc()
			return nil
		},
	})
	return newK8SessionRepository(sw), nil
}

func provideResultsRepository(browserkubeClient browserkubeclientv1.Interface, env *provision.Config,
) sessionresult.Repository {
	return newK8sResultsRepository(browserkubeClient.SessionResults(env.BrowserNS))
}

func provideProvisioner(
	clientset *kubernetes.Clientset,
	browserkubeClient browserkubeclientv1.Interface,
	envConf *provision.Config,
) provision.Provisioner {
	return newK8sWebDriverProvisioner(clientset, browserkubeClient, envConf)
}

func provideClientSet() (*kubernetes.Clientset, browserkubeclientv1.Interface, error) {
	var clientset *kubernetes.Clientset
	var err error

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	if err = browserkubeclientv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	browserkubeClient, err := browserkubeclientv1.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return clientset, browserkubeClient, errors.WithStack(err)
}
