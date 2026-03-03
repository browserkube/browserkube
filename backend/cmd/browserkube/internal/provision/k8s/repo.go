package provisionk8s

import (
	"context"

	"github.com/browserkube/browserkube/pkg/session"
)

type k8sSessionRepository struct {
	sessionWatch SessionWatchInterface
}

func newK8SessionRepository(sessionWatch SessionWatchInterface) session.Repository {
	return &k8sSessionRepository{sessionWatch: sessionWatch}
}

func (pps *k8sSessionRepository) FindByID(id string) (*session.Session, error) {
	sess, err := pps.sessionWatch.LoadByID(id)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (pps *k8sSessionRepository) FindAll() ([]*session.Session, error) {
	sessionIDs := pps.sessionWatch.GetSessions()
	sessions := make([]*session.Session, len(sessionIDs))
	for i, sID := range sessionIDs {
		sess, err := pps.FindByID(sID)
		if err != nil {
			return nil, err
		}
		sessions[i] = sess
	}
	return sessions, nil
}

func (pps *k8sSessionRepository) Quota() (int, int, error) {
	currentQ, maxQ := pps.sessionWatch.GetQuotas()
	current, _ := currentQ.AsInt64()
	maxQuotaInt, _ := maxQ.AsInt64()
	return int(current), int(maxQuotaInt), nil
}

func (pps *k8sSessionRepository) Watch(ctx context.Context) <-chan *session.Session {
	return pps.sessionWatch.Watch(ctx)
}

func (pps *k8sSessionRepository) Save(_ *session.Session) error {
	// do nothing since we use k8s sessionWatch for persisting session info
	return nil
}

func (pps *k8sSessionRepository) Delete(_ string) error {
	// do nothing
	return nil
}
