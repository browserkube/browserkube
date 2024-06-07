package browserkubehttp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func WriteJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		zap.S().Errorf("Unable to write payload: %v", err)
		return errors.WithStack(err)
	}
	return nil
}

func WritePlainText(w http.ResponseWriter, statusCode int, text string) error {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(statusCode)

	_, err := io.WriteString(w, text)

	return err
}

type HTTPErr struct {
	error
	StatusCode int
}

func NewHTTPErr(statusCode int, err error) *HTTPErr {
	return &HTTPErr{error: err, StatusCode: statusCode}
}

func Handler(f func(w http.ResponseWriter, rq *http.Request) error) http.HandlerFunc {
	logger := zap.S()
	return func(w http.ResponseWriter, rq *http.Request) {
		if err := f(w, rq); err != nil {
			logger.Error(err)

			var statusCode int
			var httpErr *HTTPErr
			if errors.As(err, &httpErr) {
				statusCode = httpErr.StatusCode
			} else {
				statusCode = http.StatusInternalServerError
			}
			if wErr := WriteJSON(w, statusCode, map[string]string{"error": err.Error()}); wErr != nil {
				logger.Error(wErr)
			}
		}
	}
}
