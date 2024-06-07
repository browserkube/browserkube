package provisionk8s

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/browserkube/browserkube/browserkube/internal/provision/k8s/mocks"
	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

func Test_k8sResultsRepository_FindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		prepareFunc func(mockSessionResults *mocks.SessionResultsInterface)
	}{
		{
			name: "FindByID: success",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: false,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(sessionResult, nil).Maybe()
			},
		},
		{
			name: "FindByID: return error, error expected",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: true,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
		},
		{
			name: "FindByID: len(res.Items) == 0, error expected",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: true,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(&v1.SessionResultList{}, nil).Maybe()
			},
		},
		{
			name: "FindByID: len(res.Items) > 1, error expected",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: true,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(sessionResultList, nil).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionResults := mocks.NewSessionResultsInterface(t)
			tt.prepareFunc(mockSessionResults)

			repo := newK8sResultsRepository(mockSessionResults)

			result, err := repo.FindByID(tt.args.ctx, tt.args.id)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Contains(t, result.Name, tt.args.id)
			} else {
				require.Error(t, err)
			}
		})
	}
}

var sessionResult = &v1.SessionResultList{
	TypeMeta: metav1.TypeMeta{},
	ListMeta: metav1.ListMeta{
		Continue: "continueToken",
	},
	Items: []v1.SessionResult{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "1",
				Namespace: "namespace",
			},
		},
	},
}

var sessionResultList = &v1.SessionResultList{
	TypeMeta: metav1.TypeMeta{},
	ListMeta: metav1.ListMeta{},
	Items: []v1.SessionResult{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "1",
				Namespace: "namespace",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "1",
				Namespace: "namespace",
			},
		},
	},
}

func Test_k8sResultsRepository_FindAll(t *testing.T) {
	type args struct {
		ctx           context.Context
		limit         int
		continueToken string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		prepareFunc func(mockSessionResults *mocks.SessionResultsInterface)
	}{
		{
			name: "FindAll: success",
			args: args{
				ctx:           context.Background(),
				limit:         1,
				continueToken: "continueToken",
			},
			wantErr: false,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(sessionResult, nil).Maybe()
			},
		},
		{
			name: "FindAll: return error, error expected",
			args: args{
				ctx:           context.Background(),
				limit:         1,
				continueToken: "continueToken",
			},
			wantErr: true,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
		},
		{
			name: "FindAll: len(res.Items) == 0, error not expected",
			args: args{
				ctx:           context.Background(),
				limit:         1,
				continueToken: "continueToken",
			},
			wantErr: false,
			prepareFunc: func(mockSessionResults *mocks.SessionResultsInterface) {
				mockSessionResults.On("List", context.Background(), mock.Anything).Return(&v1.SessionResultList{}, nil).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionResults := mocks.NewSessionResultsInterface(t)
			tt.prepareFunc(mockSessionResults)

			repo := newK8sResultsRepository(mockSessionResults)

			result, err := repo.FindAll(tt.args.ctx, tt.args.limit, tt.args.continueToken)
			if !tt.wantErr && len(result.Items) > 0 {
				require.NoError(t, err)
				require.NotEmpty(t, result)
				require.Contains(t, result.ContinueToken, tt.args.continueToken)
			} else if !tt.wantErr && len(result.Items) == 0 {
				require.Empty(t, result)
			} else {
				require.Error(t, err)
			}
		})
	}
}
