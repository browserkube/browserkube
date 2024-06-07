package provisionk8s

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/browserkube/browserkube/browserkube/internal/provision/k8s/mocks"
	v1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
)

func Test_k8sSessionRepository_FindByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		prepareFunc func(mockSessionWatch *mocks.SessionWatchInterface)
	}{
		{
			name: "FindByID: success",
			args: args{
				id: "1",
			},
			wantErr: false,
			prepareFunc: func(mockSessionWatch *mocks.SessionWatchInterface) {
				mockSessionWatch.On("LoadByID", mock.Anything).Return(mockSession, nil).Maybe()
			},
		},
		{
			name: "FindByID: return error, error expected",
			args: args{
				id: "1",
			},
			wantErr: true,
			prepareFunc: func(mockSessionWatch *mocks.SessionWatchInterface) {
				mockSessionWatch.On("LoadByID", mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionWatch := mocks.NewSessionWatchInterface(t)
			tt.prepareFunc(mockSessionWatch)

			repo := newK8SessionRepository(mockSessionWatch)

			session, err := repo.FindByID(tt.args.id)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Contains(t, session.ID, tt.args.id)
			} else {
				require.Error(t, err)
			}
		})
	}
}

var mockSession = &session.Session{
	ID:      "1",
	State:   "running",
	Browser: &v1.Browser{},
	Caps:    &session.Capabilities{},
}

func Test_k8sSessionRepository_FindAll(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     bool
		prepareFunc func(mockSessionWatch *mocks.SessionWatchInterface)
	}{
		{
			name:    "FindAll: success",
			wantErr: false,
			prepareFunc: func(mockSessionWatch *mocks.SessionWatchInterface) {
				mockSessionWatch.On("LoadByID", mock.Anything).Return(mockSession, nil).Maybe()
				mockSessionWatch.On("GetSessions").Return([]string{"1"}).Maybe()
			},
		},
		{
			name:    "FindAll: return error, error expected",
			wantErr: true,
			prepareFunc: func(mockSessionWatch *mocks.SessionWatchInterface) {
				mockSessionWatch.On("LoadByID", mock.Anything).Return(nil, errors.New("error")).Maybe()
				mockSessionWatch.On("GetSessions").Return([]string{"1"}).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionWatch := mocks.NewSessionWatchInterface(t)
			tt.prepareFunc(mockSessionWatch)

			repo := newK8SessionRepository(mockSessionWatch)

			sessions, err := repo.FindAll()
			if !tt.wantErr {
				require.NoError(t, err)
				require.NotEmpty(t, sessions)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_k8sSessionRepository_Quota(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     bool
		prepareFunc func(mockSessionWatch *mocks.SessionWatchInterface)
	}{
		{
			name:    "Quota: success",
			wantErr: false,
			prepareFunc: func(mockSessionWatch *mocks.SessionWatchInterface) {
				mockSessionWatch.On("GetQuotas").Return(resource.Quantity{}, resource.Quantity{}).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionWatch := mocks.NewSessionWatchInterface(t)
			tt.prepareFunc(mockSessionWatch)

			repo := newK8SessionRepository(mockSessionWatch)

			currentQ, maxQ, err := repo.Quota()
			if !tt.wantErr {
				require.NoError(t, err)
				require.NotNil(t, currentQ)
				require.NotNil(t, maxQ)
			} else {
				require.Error(t, err)
			}
		})
	}
}
