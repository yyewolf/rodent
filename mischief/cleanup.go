package mischief

import (
	"context"
	"os"
)

func (mischief *Mischief) cleanBrowserPool(ctx context.Context) error {
	for i := 0; i < mischief.concurrency; i++ {
		select {
		case browser := <-mischief.browserPool:
			if browser == nil {
				continue
			}

			browser.Close()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Clean up the profile directories in /tmp/rod
	return os.RemoveAll("/tmp/rod")
}

// Destroy destroys the Mischief instance.
//
// It waits for all the browsers in the pool and closes them.
func (mischief *Mischief) Destroy(ctx context.Context) error {
	mischief.logger.Info("mischief is destroying")
	err := mischief.cleanBrowserPool(ctx)
	if err != nil {
		return err
	}

	mischief.watchratCancel()

	return nil
}

// Cleanup destroys the Mischief instance.
//
// It waits for all the browsers in the pool and closes them.
func (mischief *Mischief) Cleanup(ctx context.Context) error {
	mischief.logger.Info("mischief is cleaning up")
	err := mischief.cleanBrowserPool(ctx)
	if err != nil {
		return err
	}

	return mischief.initialize()
}
