package browserkubeutil

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Retry executes callback func until it executes successfully
func Retry(attempts int, timeout time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	log := zap.S()
	var err error
	for i := 0; i < attempts; i++ {
		var res interface{}
		res, err = callback()
		if err == nil {
			return res, nil
		}

		<-time.After(timeout)
		log.Infof("Retrying... Attempt: %d. Left: %d. Err: %v", i+1, attempts-1-i, err)
	}
	return nil, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}

// RetryWithTimeout executes action with defined timeout until receives timeout signal
func RetryWithTimeout(waitTimeout, initialDelay, timeout time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	log := zap.S()
	var err error
	var attempts int

	waitTimer := time.NewTimer(waitTimeout)
	defer waitTimer.Stop()

	attemptTimer := time.After(initialDelay)
	for {
		select {
		case <-waitTimer.C:
			return nil, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
		case <-attemptTimer:
			var res interface{}
			res, err = callback()
			if err == nil {
				return res, nil
			}
			attempts++
			attemptTimer = time.After(timeout)
			log.Debugf("Retrying... Attempt: %d. Err: %v", attempts, err)
		}
	}
}
