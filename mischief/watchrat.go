package mischief

import (
	"context"
	"log/slog"
	"time"
)

func (mischief *Mischief) watchrat(ctx context.Context) {
	for {
		select {
		// case <-time.After(mischief.watchratInterval):
		case <-time.After(5 * time.Minute):
			err := mischief.Cleanup(ctx)
			if err != nil {
				mischief.logger.Error("mischief failed to cleanup", slog.Any("error", err))
			}
		case <-ctx.Done():
			mischief.logger.Info("mischief is stopping watchrat")
			return
		}
	}
}
