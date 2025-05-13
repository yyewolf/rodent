package mischief

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/yyewolf/rodent/rat"
	"golang.org/x/exp/rand"
)

func randomBackoff() time.Duration {
	return time.Duration(1+rand.Intn(1000)) * time.Millisecond
}

func (mischief *Mischief) clearRat(ctx context.Context, rat *rat.Rat, recreate bool) error {
	if recreate {
		err := rat.Recreate()
		if err != nil {
			return fmt.Errorf("failed to recreate rat: %w", err)
		}
	} else {
		err := rat.Close()
		if err != nil {
			return fmt.Errorf("failed to close rat: %w", err)
		}
	}

	select {
	case <-time.After(randomBackoff()):
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	}

	return nil
}

func (mischief *Mischief) cleanBrowserPool(ctx context.Context, recreate bool) error {
	var errorList error

	for _, rat := range mischief.rats {
		rat.Lock()
		err := mischief.clearRat(ctx, rat, recreate)
		if err != nil {
			mischief.logger.Error("mischief failed to clear rat", slog.Any("error", err))
			errorList = errors.Join(errorList, err)
		}
		rat.Unlock()
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
	err := mischief.cleanBrowserPool(ctx, false)
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

	return mischief.cleanBrowserPool(ctx, true)
}
