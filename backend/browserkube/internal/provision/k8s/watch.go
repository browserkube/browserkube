package provisionk8s

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/util/broadcast"
)

const (
	sessionsIndex        = "webdriver-go-session-index"
	quotaSessionsNameTpl = "%s-sessions"

	broadcastBuflen = 100
)

// sessionWatch represents dynamic listener of k8s selenium-pods.
type sessionWatch struct {
	ns                string
	quotaLister       v12.ResourceQuotaNamespaceLister
	logger            *zap.SugaredLogger
	broadcast         broadcast.Broadcaster[*session.Session]
	browserkubeClient browserkubeclientv1.Interface

	quotaInformer    cache.SharedIndexInformer
	browsersInformer cache.SharedIndexInformer
}

// newSessionWatch creates new instance.
func newSessionWatch(
	clientset *kubernetes.Clientset,
	browserkubeClient browserkubeclientv1.Interface,
	ns string,
) (*sessionWatch, error) {
	logger := zap.S()

	browsersWatch := cache.NewListWatchFromClient(browserkubeClient.RESTClient(), "browsers", ns, fields.Everything())
	// Bind the workqueue to a browsersCache with the help of an browsersInformer. This way we make sure that
	// whenever the browserCache is updated, the browser key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Browser than the version which was responsible for triggering the update.
	browsersInformer := cache.NewSharedIndexInformer(browsersWatch, &browserkubev1.Browser{}, resyncPeriod,
		cache.Indexers{
			sessionsIndex: func(obj interface{}) ([]string, error) {
				browser, ok := obj.(*browserkubev1.Browser)
				if ok {
					if lv := browser.Labels[browserkubev1.LabelBrowserVisibility]; lv != "true" {
						return []string{}, nil
					}
					if sess := browser.Name; sess != "" {
						return []string{sess}, nil
					}
				}
				return []string{}, nil
			},
		})

	broadcaster := broadcast.NewBroadcaster[*session.Session](broadcastBuflen)
	broadcastF := func(obj interface{}) {
		b, ok := obj.(*browserkubev1.Browser)
		if !ok {
			return
		}
		if lv := b.Labels[browserkubev1.LabelBrowserVisibility]; lv != "true" {
			return
		}
		s, convErr := asSession(b)
		if convErr != nil {
			return
		}
		broadcaster.Submit(s)
	}

	if _, err := browsersInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    broadcastF,
		UpdateFunc: func(oldObj, newObj interface{}) { broadcastF(newObj) },
		DeleteFunc: broadcastF,
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	/// resource quotas informer
	quotaWatch := cache.NewFilteredListWatchFromClient(clientset.CoreV1().RESTClient(), "resourcequotas", ns,
		func(options *metav1.ListOptions) {
			options.LabelSelector = metav1.FormatLabelSelector(metav1.SetAsLabelSelector(map[string]string{
				browserkubev1.LabelApp: browserkubev1.LabelValueApplicationBrowserkube,
			}))
		})
	// Bind the workqueue to a podsCache with the help of an podsInformer. This way we make sure that
	// whenever the podsCache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	quotaInformer := cache.NewSharedIndexInformer(quotaWatch, &v1.ResourceQuota{}, resyncPeriod,
		cache.Indexers{})

	return &sessionWatch{
		logger:           logger,
		broadcast:        broadcaster,
		quotaLister:      v12.NewResourceQuotaLister(quotaInformer.GetIndexer()).ResourceQuotas(ns),
		ns:               ns,
		quotaInformer:    quotaInformer,
		browsersInformer: browsersInformer,
	}, nil
}

func (ac *sessionWatch) Start(ctx context.Context) {
	go ac.quotaInformer.Run(ctx.Done())
	go ac.browsersInformer.Run(ctx.Done())

	// We can now warm up the podsCache for initial synchronization.
	// Let's suppose that we knew about a pod "mypod" on our last run, therefore add it to the podsCache.
	// If this pod is not there anymore, the controller will be notified about the removal after the
	// podsCache has synchronized.
	if err := ac.quotaInformer.GetIndexer().Resync(); err != nil {
		ac.logger.Error(err)
	}
	if err := ac.browsersInformer.GetIndexer().Resync(); err != nil {
		ac.logger.Error(err)
	}
}

func (ac *sessionWatch) Watch(ctx context.Context) <-chan *session.Session {
	sCh := make(chan *session.Session)
	ac.broadcast.Register(sCh)
	go func(c chan *session.Session) {
		<-ctx.Done()
		ac.broadcast.Deregister(c)
		close(c)
	}(sCh)
	return sCh
}

// GetSessions loads all provisioned selenium services
func (ac *sessionWatch) GetSessions() []string {
	return ac.browsersInformer.GetIndexer().ListIndexFuncValues(sessionsIndex)
}

// LoadByID loads provisioned selenium service by session ID
func (ac *sessionWatch) LoadByID(sessionID string) (*session.Session, error) {
	browser, err := ac.getBrowser(sessionID)
	if err != nil {
		return nil, err
	}

	return asSession(browser)
}

// getBrowser returns browser from sessionWatch or an error
func (ac *sessionWatch) getBrowser(sessionID string) (*browserkubev1.Browser, error) {
	browsers, err := ac.browsersInformer.GetIndexer().ByIndex(sessionsIndex, sessionID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//nolint:gomnd // checks browser number
	if l := len(browsers); l < 1 {
		return nil, errors.Errorf("Browser [%s] is not accessible. Available browsers: [%d]", sessionID, l)
	}
	//nolint:forcetypeassert
	browser := browsers[0].(*browserkubev1.Browser)

	return browser, nil
}

// Exists checks whether browser for the given session ID exists
func (ac *sessionWatch) Exists(sessionID string) (bool, error) {
	browsers, err := ac.browsersInformer.GetIndexer().ByIndex(sessionsIndex, sessionID)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return len(browsers) > 0, nil
}

// GetQuotas returns defined quotas
func (ac *sessionWatch) GetQuotas() (resource.Quantity, resource.Quantity) {
	sessionsQuota, err := ac.quotaLister.Get(fmt.Sprintf(quotaSessionsNameTpl, ac.ns))
	if err != nil {
		ac.logger.Errorf("Unable to get browsers resource quotas: %v", err)
		return resource.Quantity{}, resource.Quantity{}
	}
	resourceName := fmt.Sprintf("count/%s", browserkubev1.GroupVersion.WithResource("browsers").GroupResource().String())

	if qHard, ok := sessionsQuota.Status.Hard[v1.ResourceName(resourceName)]; ok {
		if qUsed, uOk := sessionsQuota.Status.Used[v1.ResourceName(resourceName)]; uOk {
			return qUsed, qHard
		}
	}

	return resource.Quantity{}, resource.Quantity{}
}

// IsNewSessionAllowed returns TRUE if creation is allowed
func (ac *sessionWatch) IsNewSessionAllowed() bool {
	used, hard := ac.GetQuotas()
	if hard.IsZero() {
		return true
	}
	return hard.Cmp(used) > 0
}

func asSession(browser *browserkubev1.Browser) (*session.Session, error) {
	var caps session.Capabilities
	capsStr := string(browser.Spec.Caps)
	if capsStr != "" {
		if err := json.Unmarshal([]byte(capsStr), &caps); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &session.Session{
		ID:      browser.Name,
		State:   asStatus(browser),
		Browser: browser,
		Caps:    &caps,
	}, nil
}

func asStatus(b *browserkubev1.Browser) string {
	switch b.Status.Phase {
	case browserkubev1.PhasePending:
		return "pending"
	case browserkubev1.PhaseRunning:
		if b.DeletionTimestamp != nil {
			return "terminating"
		}
		return "running"
	case browserkubev1.PhaseTerminated:
		return "terminated"
	default:
		return "pending"
	}
}

//go:generate mockery --name SessionWatchInterface --output mocks
type SessionWatchInterface interface {
	Watch(ctx context.Context) <-chan *session.Session
	GetSessions() []string
	LoadByID(sessionID string) (*session.Session, error)
	Exists(sessionID string) (bool, error)
	GetQuotas() (resource.Quantity, resource.Quantity)
}
