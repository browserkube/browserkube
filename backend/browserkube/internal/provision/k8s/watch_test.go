package provisionk8s

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/browserkube/browserkube/browserkube/internal/provision/k8s/mocks"
	v1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
	mocks2 "github.com/browserkube/browserkube/pkg/util/broadcast/mocks"
)

func Test_newSessionWatch(t *testing.T) {
	mockInterface := mocks.NewRestMock(t)
	mockInterfaceV1 := mocks.NewInterface(t)
	mockInterfaceV1.On("RESTClient").Return(mockInterface).Maybe()
	got, err := newSessionWatch(&kubernetes.Clientset{}, mockInterfaceV1, "ns")
	require.NoError(t, err)
	require.NotEmpty(t, got)
}

func Test_sessionWatch_LoadByID(t *testing.T) {
	type args struct {
		sessionID string
	}
	tests := []struct {
		name                     string
		args                     args
		want                     *session.Session
		wantErr                  bool
		prepareMockIndex         func(mockIndex *mocks.Indexer)
		prepareMockIndexInformer func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer)
	}{
		{
			name: "LoadByID: success",
			args: args{
				sessionID: "sessionID",
			},
			want: &session.Session{
				ID:    "sessionID",
				State: "pending",
				Browser: &v1.Browser{
					ObjectMeta: metav1.ObjectMeta{Name: "sessionID"},
				},
				Caps: &session.Capabilities{},
			},
			wantErr: false,
			prepareMockIndex: func(mockIndex *mocks.Indexer) {
				mockIndex.On("ByIndex", mock.Anything, mock.Anything).Return(
					[]interface{}{&v1.Browser{ObjectMeta: metav1.ObjectMeta{Name: "sessionID"}}},
					nil,
				).Maybe()
			},
			prepareMockIndexInformer: func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer) {
				mockIndexInformer.On("GetIndexer").Return(mockIndex).Maybe()
			},
		},
		{
			name: "LoadByID: return error, error expected",
			args: args{
				sessionID: "sessionID",
			},
			wantErr: true,
			prepareMockIndex: func(mockIndex *mocks.Indexer) {
				mockIndex.On("ByIndex", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
			prepareMockIndexInformer: func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer) {
				mockIndexInformer.On("GetIndexer").Return(mockIndex).Maybe()
			},
		},
		{
			name: "LoadByID: len(browsers) < 1, error expected",
			args: args{
				sessionID: "sessionID",
			},
			wantErr: true,
			prepareMockIndex: func(mockIndex *mocks.Indexer) {
				mockIndex.On("ByIndex", mock.Anything, mock.Anything).Return([]interface{}{}, nil).Maybe()
			},
			prepareMockIndexInformer: func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer) {
				mockIndexInformer.On("GetIndexer").Return(mockIndex).Maybe()
			},
		},
		{
			name: "LoadByID: string(browser.Spec.Caps) != '', error expected",
			args: args{
				sessionID: "sessionID",
			},
			wantErr: true,
			prepareMockIndex: func(mockIndex *mocks.Indexer) {
				mockIndex.On("ByIndex", mock.Anything, mock.Anything).Return(
					[]interface{}{&v1.Browser{
						ObjectMeta: metav1.ObjectMeta{Name: "sessionID"},
						Spec: v1.BrowserSpec{
							Caps: []byte("capabilities"),
						},
					}},
					nil,
				).Maybe()
			},
			prepareMockIndexInformer: func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer) {
				mockIndexInformer.On("GetIndexer").Return(mockIndex).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIndex := mocks.NewIndexer(t)
			tt.prepareMockIndex(mockIndex)

			mockIndexInformer := mocks.NewSharedIndexInformer(t)
			tt.prepareMockIndexInformer(mockIndexInformer, mockIndex)

			ac := &sessionWatch{
				browsersInformer: mockIndexInformer,
			}
			got, err := ac.LoadByID(tt.args.sessionID)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_sessionWatch_Exists(t *testing.T) {
	type args struct {
		sessionID string
	}
	tests := []struct {
		name                     string
		args                     args
		want                     bool
		wantErr                  bool
		prepareMockIndex         func(mockIndex *mocks.Indexer)
		prepareMockIndexInformer func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer)
	}{
		{
			name: "Exists: success",
			args: args{
				sessionID: "sessionID",
			},
			want:    true,
			wantErr: false,
			prepareMockIndex: func(mockIndex *mocks.Indexer) {
				mockIndex.On("ByIndex", mock.Anything, mock.Anything).Return(
					[]interface{}{&v1.Browser{ObjectMeta: metav1.ObjectMeta{Name: "sessionID"}}},
					nil,
				).Maybe()
			},
			prepareMockIndexInformer: func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer) {
				mockIndexInformer.On("GetIndexer").Return(mockIndex).Maybe()
			},
		},
		{
			name: "Exists: return error, error expected",
			args: args{
				sessionID: "sessionID",
			},
			wantErr: true,
			prepareMockIndex: func(mockIndex *mocks.Indexer) {
				mockIndex.On("ByIndex", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
			prepareMockIndexInformer: func(mockIndexInformer *mocks.SharedIndexInformer, mockIndex *mocks.Indexer) {
				mockIndexInformer.On("GetIndexer").Return(mockIndex).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIndex := mocks.NewIndexer(t)
			tt.prepareMockIndex(mockIndex)

			mockIndexInformer := mocks.NewSharedIndexInformer(t)
			tt.prepareMockIndexInformer(mockIndexInformer, mockIndex)

			ac := &sessionWatch{
				browsersInformer: mockIndexInformer,
			}
			got, err := ac.Exists(tt.args.sessionID)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_sessionWatch_GetQuotas(t *testing.T) {
	tests := []struct {
		name        string
		want        resource.Quantity
		want1       resource.Quantity
		prepareFunc func(mockQuotaLister *mocks.ResourceQuotaNamespaceLister)
	}{
		{
			name: "GetQuotas: success",
			want: resource.Quantity{
				Format: "1",
			},
			want1: resource.Quantity{
				Format: "1",
			},
			prepareFunc: func(mockQuotaLister *mocks.ResourceQuotaNamespaceLister) {
				mockQuotaLister.On("Get", mock.Anything).Return(&coreV1.ResourceQuota{
					Status: coreV1.ResourceQuotaStatus{
						Hard: coreV1.ResourceList{
							"count/browsers.api.browserkube.io": resource.Quantity{
								Format: "1",
							},
						},
						Used: coreV1.ResourceList{
							"count/browsers.api.browserkube.io": resource.Quantity{
								Format: "1",
							},
						},
					},
				}, nil)
			},
		},
		{
			name:  "GetQuotas: return error, error expected",
			want:  resource.Quantity{},
			want1: resource.Quantity{},
			prepareFunc: func(mockQuotaLister *mocks.ResourceQuotaNamespaceLister) {
				mockQuotaLister.On("Get", mock.Anything).Return(nil, errors.New("error"))
			},
		},
		{
			name:  "GetQuotas: return empty resource.Quantity{}",
			want:  resource.Quantity{},
			want1: resource.Quantity{},
			prepareFunc: func(mockQuotaLister *mocks.ResourceQuotaNamespaceLister) {
				mockQuotaLister.On("Get", mock.Anything).Return(&coreV1.ResourceQuota{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQuotaLister := mocks.NewResourceQuotaNamespaceLister(t)
			tt.prepareFunc(mockQuotaLister)

			logger := fakeLogger{
				SugaredLogger: zap.NewExample().Sugar(),
			}

			ac := &sessionWatch{
				logger:      logger.SugaredLogger,
				ns:          "browser",
				quotaLister: mockQuotaLister,
			}
			got, got1 := ac.GetQuotas()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sessionWatch.GetQuotas() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("sessionWatch.GetQuotas() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type fakeLogger struct {
	*zap.SugaredLogger
}

func (l fakeLogger) Errorf(template string, args ...interface{}) {
	l.SugaredLogger.Errorf(template, args...)
}

func (l fakeLogger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

func Test_sessionWatch_Start(t *testing.T) {
	tests := []struct {
		name        string
		prepareFunc func(mockStore *mocks.Store)
	}{
		{
			name: "Start: success",
			prepareFunc: func(mockStore *mocks.Store) {
				mockStore.On("Resync").Return(nil).Maybe()
			},
		},
		{
			name: "Start: return error, error expected",
			prepareFunc: func(mockStore *mocks.Store) {
				mockStore.On("Resync").Return(errors.New("error")).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer leaktest.Check(t)()

			mockStore := mocks.NewStore(t)
			tt.prepareFunc(mockStore)

			mockIndexInformer := mocks.NewSharedIndexInformer(t)
			mockIndexInformer.On("Run", mock.Anything).Maybe()
			mockIndexInformer.On("GetIndexer").Return(&fakeIndex{
				Store: *mockStore,
			}).Maybe()

			logger := fakeLogger{
				SugaredLogger: zap.NewExample().Sugar(),
			}

			ac := &sessionWatch{
				logger:           logger.SugaredLogger,
				quotaInformer:    mockIndexInformer,
				browsersInformer: mockIndexInformer,
			}
			ac.Start(context.Background())
		})
	}
}

type fakeIndex struct {
	mocks.Store
}

func (fakeIndex) Index(indexName string, obj interface{}) ([]interface{}, error) {
	return []interface{}{}, nil
}

func (fakeIndex) IndexKeys(indexName, indexedValue string) ([]string, error) { return []string{}, nil }

func (fakeIndex) ListIndexFuncValues(indexName string) []string { return []string{} }

func (fakeIndex) ByIndex(indexName, indexedValue string) ([]interface{}, error) {
	return []interface{}{}, nil
}

func (fakeIndex) GetIndexers() cache.Indexers { return cache.Indexers{} }

func (fakeIndex) AddIndexers(newIndexers cache.Indexers) error { return nil }

func Test_sessionWatch_Watch(t *testing.T) {
	mockBroadcaster := mocks2.NewBroadcaster[*session.Session](t)
	mockBroadcaster.On("Register", mock.Anything)
	ac := &sessionWatch{
		broadcast: mockBroadcaster,
	}
	got := ac.Watch(context.Background())
	var want <-chan *session.Session
	require.IsType(t, want, got)
}

func Test_sessionWatch_GetSessions(t *testing.T) {
	mockIndexer := mocks.NewIndexer(t)
	mockIndexer.On("ListIndexFuncValues", mock.Anything).Return([]string{"1", "2"}).Maybe()

	mockIndexInformer := mocks.NewSharedIndexInformer(t)
	mockIndexInformer.On("GetIndexer").Return(mockIndexer).Maybe()

	ac := &sessionWatch{
		browsersInformer: mockIndexInformer,
	}
	got := ac.GetSessions()
	require.NotEmpty(t, got)
}
