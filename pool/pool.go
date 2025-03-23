package pool

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
)

func GetFromPoolWithTimeout[K any](pool rod.Pool[K], timeout time.Duration) (*K, error) {
	select {
	case elem := <-pool:
		return elem, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
}
