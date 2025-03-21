package reaper

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	ErrReapingStopped = fmt.Errorf("reaping stopped")
)

// Reaper is a struct that reaps child processes.
//
// It listens for SIGCHLD signals and reaps child processes.
type Reaper struct {
	// logger is the logger of the Mischief instance
	logger *slog.Logger

	// ctx is the context of the Reaper instance
	ctx context.Context
	// cancel is the cancel function of the Reaper instance
	cancel context.CancelFunc
}

func NewReaper(logger *slog.Logger) *Reaper {
	ctx, cancel := context.WithCancel(context.Background())

	return &Reaper{
		logger: logger,

		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *Reaper) Start() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGCHLD)

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-sigCh:
				for {
					var status syscall.WaitStatus
					var rusage syscall.Rusage
					pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, &rusage)
					if pid <= 0 || err != nil {
						break
					}

					// Convert RAM usage to MB for easier reading
					ramUsage := fmt.Sprintf("%.2f MB", float64(rusage.Maxrss)/1024)
					r.logger.Info("Reaped child process", slog.Any("pid", pid), slog.Any("freed ram usage", ramUsage))
				}
			}
		}
	}()
}

func (r *Reaper) Shutdown() {
	r.cancel()
}
