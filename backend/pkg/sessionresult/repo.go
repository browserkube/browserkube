package sessionresult

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

const (
	VideosPath = "/videos/"
)

const (
	BrowserLogFileName = "browser.log"
	VideoFileName      = "video.mp4"
	MessageLogFileName = "message.log"
)

type Repository interface {
	Create(ctx context.Context, req *Result) (*Result, error)
	FindAll(ctx context.Context, limit int, continueToken string) (*browserkubeutil.Page[*Result], error)
	FindByID(ctx context.Context, name string) (*Result, error)
	Delete(cctx context.Context, name string, opts metav1.DeleteOptions) error
}

type Result struct {
	browserkubev1.SessionResult
}

func (r *Result) GetCaps() (session.Capabilities, error) {
	var caps session.Capabilities
	capsStr := string(r.Spec.Browser.Caps)
	if capsStr != "" {
		if err := json.Unmarshal([]byte(capsStr), &caps); err != nil {
			return caps, errors.WithStack(err)
		}
	}
	return caps, nil
}
