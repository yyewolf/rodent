package mischief

import (
	"errors"
	"log/slog"
	"time"

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

	rat, err := mischief.getRat()
	if err != nil {
		return nil, errors.Join(ErrGettingBrowser, err)
	}
	defer mischief.ratPool.Put(rat)

	rat.Lock()
	defer rat.Unlock()

	page, err := rat.GetPage()
	if err != nil {
		return nil, errors.Join(ErrGettingPage, err)
	}
	defer rat.PutPage(page)

	page = page.Timeout(mischief.pageStabilityTimeout)

	err = page.Navigate(url)
	if err != nil {
		return nil, errors.Join(ErrNavigatingToPage, err)
	}

	err = page.WaitDOMStable(time.Millisecond, 0)
	if err != nil {
		return nil, errors.Join(ErrWaitingForPageToBeStable, err)
	}

	page = page.CancelTimeout()

	screenshotParams := &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	}

	bytes, err := page.Screenshot(false, screenshotParams)
	if err != nil {
		return nil, errors.Join(ErrWhileTakingScreenshot, err)
	}

	return bytes, nil
}
