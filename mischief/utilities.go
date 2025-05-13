package mischief

import (
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/yyewolf/rodent/pool"
	"github.com/yyewolf/rodent/rat"
)

func createBrowser(controlUrl *string) func(*rat.Rat) error {
	return func(rat *rat.Rat) error {
		var newControlUrl string

		if controlUrl == nil {
			uri, err := launcher.New().Bin(os.Getenv("BROWSER_PATH")).Launch()
			if err != nil {
				return err
			}

			newControlUrl = uri
		} else {
			newControlUrl = *controlUrl
		}

		browser := rod.New().ControlURL(newControlUrl)

		err := browser.Connect()
		if err != nil {
			return err
		}

		rat.Browser = browser

		return rat.Initialize()
	}
}

func (mischief *Mischief) getRat() (*rat.Rat, error) {
	rat, err := pool.GetFromPoolWithTimeout(mischief.ratPool, mischief.browserRetakeTimeout)
	if err != nil {
		return nil, err
	}

	return rat, nil
}
