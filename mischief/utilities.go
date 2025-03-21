package mischief

import (
	"errors"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func createBrowser(controlUrl *string) func() (*Rat, error) {
	return func() (b *Rat, err error) {
		if controlUrl == nil {
			uri, err := launcher.New().Bin(os.Getenv("BROWSER_PATH")).Launch()
			if err != nil {
				return nil, err
			}

			controlUrl = &uri
		}

		browser := rod.New().ControlURL(*controlUrl)

		err = browser.Connect()
		if err != nil {
			return nil, err
		}

		return &Rat{
			CreatedAt: time.Now(),
			Browser:   browser,
		}, nil
	}
}

func getFromPoolWithTimeout[K any](pool rod.Pool[K], timeout time.Duration) (*K, error) {
	select {
	case elem := <-pool:
		return elem, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
}

func (mischief *Mischief) getBrowser() (*Rat, error) {
	rat, err := getFromPoolWithTimeout(mischief.browserPool, mischief.browserRetakeTimeout)
	if err != nil {
		return nil, err
	}

	return rat, nil
}
