package mischief

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-rod/rod"
)

// Mischief is the main struct of the Mischief package.
//
// It is used to take screenshots of URLs.
type Mischief struct {
	// externalBrowser is a flag to use external browsers
	externalBrowser bool
	// browserUrls is a list of URLs to connect to external browsers, it should be set when externalBrowser is true
	browserUrls []string
	// concurrency is the number of browsers to use to take screenshots concurrently
	concurrency int

	// logger is the logger of the Mischief instance
	logger *slog.Logger

	// browserPool is the pool of browsers
	browserPool rod.Pool[Rat]

	// browserRetakeTimeout is the timeout used when taking a browser from the pool
	browserRetakeTimeout time.Duration
	// pageStabilityTimeout is the timeout used when waiting for the page to be stable
	pageStabilityTimeout time.Duration

	// watchratCancel is the context used to watch the rat
	watchratCancel context.CancelFunc
}

type MischiefOpt func(*Mischief)

// New creates a new Mischief instance.
// A Mischief instance is used to take screenshots of URLs.
//
// Example (and default values):
//
//	m := mischief.New(
//		mischief.WithConcurrency(1),
//		mischief.WithLogger(slog.Default()),
//		mischief.WithBrowserRetakeTimeout(5*time.Second),
//		mischief.WithPageStabilityTimeout(3*time.Second),
//	)
func New(opts ...MischiefOpt) (*Mischief, error) {
	var m Mischief

	var defaultOpts = []MischiefOpt{
		WithConcurrency(1),
		WithLogger(slog.Default()),
		WithBrowserRetakeTimeout(5 * time.Second),
		WithPageStabilityTimeout(3 * time.Second),
	}

	opts = append(defaultOpts, opts...)

	for _, opt := range opts {
		opt(&m)
	}

	err := m.initialize()
	if err != nil {
		return nil, err
	}

	watchratCtx, watchratCancel := context.WithCancel(context.Background())
	m.watchratCancel = watchratCancel

	go m.watchrat(watchratCtx)

	return &m, nil
}

// initialize initializes the Mischief instance.
//
// It creates a pool of browsers to take screenshots concurrently.
func (mischief *Mischief) initialize() error {
	mischief.browserPool = rod.NewPool[Rat](mischief.concurrency)

	var rats []*Rat = make([]*Rat, mischief.concurrency)

	// Instanciate every browser
	for i := 0; i < mischief.concurrency; i++ {
		var rat *Rat
		var err error

		if mischief.externalBrowser {
			mischief.logger.Info("mischief is connecting to external browser", slog.Any("url", mischief.browserUrls[i]))
			rat, err = mischief.browserPool.Get(createBrowser(&mischief.browserUrls[i]))
		} else {
			mischief.logger.Info("mischief is creating a new browser")
			rat, err = mischief.browserPool.Get(createBrowser(nil))
		}
		if err != nil {
			return err
		}

		rats[i] = rat
	}

	// Put them back in the pool for usage later on
	for _, browser := range rats {
		mischief.browserPool.Put(browser)
	}

	return nil
}
