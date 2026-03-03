package provisionk8s

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	browserkubeclientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

type k8sResultsRepository struct {
	resultsClient browserkubeclientv1.SessionResultsInterface
}

func newK8sResultsRepository(resultsClient browserkubeclientv1.SessionResultsInterface) sessionresult.Repository {
	return &k8sResultsRepository{resultsClient: resultsClient}
}

func (pps *k8sResultsRepository) FindByID(ctx context.Context, id string) (*sessionresult.Result, error) {
	res, err := pps.resultsClient.List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", id).String(),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(res.Items) == 0 {
		return nil, errors.Errorf("Session result with ID=%s not found", id)
	}
	if len(res.Items) > 1 {
		return nil, errors.New("Multiple objects found for given name")
	}
	result := res.Items[0]
	return &sessionresult.Result{SessionResult: result}, nil
}

func (pps *k8sResultsRepository) FindAll(ctx context.Context, limit int, continueToken string) (*browserkubeutil.Page[*sessionresult.Result], error) {
	res, err := pps.resultsClient.List(ctx, metav1.ListOptions{
		Limit:                int64(limit),
		Continue:             continueToken,
		ResourceVersion:      "",
		ResourceVersionMatch: "",
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(res.Items) == 0 {
		return browserkubeutil.Empty[*sessionresult.Result](), nil
	}

	results := make([]*sessionresult.Result, len(res.Items))
	for i, val := range res.Items {
		results[i] = &sessionresult.Result{SessionResult: val}
	}
	page := browserkubeutil.AsPage(res, results)
	return page, nil
}

func (pps *k8sResultsRepository) Create(ctx context.Context, req *sessionresult.Result) (*sessionresult.Result, error) {
	res, err := pps.resultsClient.Create(ctx, &req.SessionResult)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &sessionresult.Result{SessionResult: *res}, nil
}

func (pps *k8sResultsRepository) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	err := pps.resultsClient.Delete(ctx, name, opts)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
