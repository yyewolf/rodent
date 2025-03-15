package mischief

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
)

func createBrowser(controlUrl *string) func() (*rod.Browser, error) {
	return func() (b *rod.Browser, err error) {
		browser := rod.New()
		if controlUrl != nil {
			browser = rod.New().ControlURL(*controlUrl)
		}

		err = browser.Connect()
		if err != nil {
			return nil, err
		}

		return browser, nil
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

func (mischief *Mischief) getBrowser() (*rod.Browser, error) {
	browser, err := getFromPoolWithTimeout(mischief.browserPool, mischief.browserRetakeTimeout)
	if err != nil {
		return nil, err
	}

	return browser, nil
}
