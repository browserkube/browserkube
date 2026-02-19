package reportportal

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
)

const label = "io.browserkube.rp-project"

type settingsRepo interface {
	FindByProjectName(ctx context.Context, n string) (*ProjectSettings, error)
}

type ProjectSettings struct {
	Host        string
	ProjectName string
	AuthToken   string
}

type k8sSettingsRepo struct {
	client corev1.SecretInterface
}

func newSettingsRepo(clientset *kubernetes.Clientset, envConfig *provision.Config) settingsRepo {
	return &k8sSettingsRepo{client: clientset.CoreV1().Secrets(envConfig.BrowserNS)}
}

func (sr *k8sSettingsRepo) FindByProjectName(ctx context.Context, name string) (*ProjectSettings, error) {
	secrets, err := sr.client.List(
		ctx, metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(
				metav1.SetAsLabelSelector(map[string]string{label: name}),
			),
		},
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(secrets.Items) == 0 {
		return nil, errors.New("project settings isn't found")
	}
	secretData := secrets.Items[0].Data
	return &ProjectSettings{
		ProjectName: name,
		Host:        string(secretData["host"]),
		AuthToken:   string(secretData["authToken"]),
	}, nil
}
