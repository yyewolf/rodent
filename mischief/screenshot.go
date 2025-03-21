package mischief

import (
	"errors"
	"log/slog"

	"github.com/go-rod/rod/lib/proto"
)

// TakeScreenshot takes a screenshot of the given URL.
// It returns the screenshot as a byte slice.
//
// In order :
//
// - It gets a browser from the pool
//
// - It opens a new page with the given URL
//
// - It waits for the page to be stable
//
// - It takes the screenshot
//
// It puts the browser back to the pool after returning.
func (mischief *Mischief) TakeScreenshot(url string) ([]byte, error) {
	mischief.logger.Info("mischief is taking a screenshot", slog.Any("url", url))

	browser, err := mischief.getBrowser()
	if err != nil {
		return nil, errors.Join(ErrGettingBrowser, err)
	}
	defer mischief.browserPool.Put(browser)

	page, err := browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return nil, errors.Join(ErrOpeningPage, err)
	}
	defer page.Close()

	err = page.WaitDOMStable(mischief.pageStabilityTimeout, 30)
	if err != nil {
		return nil, errors.Join(ErrWaitingForPageToBeStable, err)
	}

	bytes, err := page.Screenshot(true, nil)
	if err != nil {
		return nil, errors.Join(ErrWhileTakingScreenshot, err)
	}

	return bytes, nil
}
