package mischief

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

func (mischief *Mischief) cleanBrowserPool(ctx context.Context) error {
	var errorList error

	for i := 0; i < mischief.concurrency; i++ {
		select {
		case rat := <-mischief.browserPool:
			if rat == nil {
				continue
			}

			err := rat.Close()
			if err != nil {
				mischief.logger.Error("mischief failed to close browser", slog.Any("error", err))
				errorList = errors.Join(errorList, err)
				continue
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Clean up the profile directories in /tmp/rod
	err := os.RemoveAll("/tmp/rod")
	if err != nil {
		mischief.logger.Warn("mischief failed to clean up profile directories", slog.Any("error", err))
	}

	return errorList
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
	mischief.cleanBrowserPool(ctx)

	return mischief.initialize()
}
