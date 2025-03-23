package mischief

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-rod/rod"
	"github.com/yyewolf/rodent/rat"
)

// Mischief is the main struct of the Mischief package.
//
// It is used to take screenshots of URLs.
type Mischief struct {
	// externalBrowser is a flag to use external browsers
	externalBrowser bool
	// browserUrls is a list of URLs to connect to external browsers, it should be set when externalBrowser is true
	browserUrls []string
	// browserConcurrency is the number of browsers to use to take screenshots concurrently
	browserConcurrency int
	// pageConcurrency is the number of pages to use to take screenshots concurrently
	pageConcurrency int

	// logger is the logger of the Mischief instance
	logger *slog.Logger

	// ratPool is the pool of browsers
	ratPool rod.Pool[rat.Rat]
	rats    []*rat.Rat

	// browserRetakeTimeout is the timeout used when taking a browser from the pool
	browserRetakeTimeout time.Duration
	// pageRetakeTimeout is the timeout used when taking a page from the pool
	pageRetakeTimeout time.Duration
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
		WithBrowserConcurrency(1),
		WithPageConcurrency(1),
		WithLogger(slog.Default()),
		WithBrowserRetakeTimeout(5 * time.Second),
		WithPageRetakeTimeout(5 * time.Second),
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
	mischief.ratPool = rod.NewPool[rat.Rat](mischief.browserConcurrency * mischief.pageConcurrency)

	var rats []*rat.Rat = make([]*rat.Rat, mischief.browserConcurrency)

	// Instanciate every browser
	for i := 0; i < mischief.browserConcurrency; i++ {
		var err error
		var uri *string

		if mischief.externalBrowser {
			uri = &mischief.browserUrls[i]
		}

		mischief.logger.Info("mischief is creating rat", slog.Any("uri", uri))
		rat, err := rat.New(
			rat.WithPagePoolLength(mischief.pageConcurrency),
			rat.WithPageRetakeTimeout(mischief.pageRetakeTimeout),
			rat.WithCreateBrowserFunc(createBrowser(uri)),
		)
		if err != nil {
			return err
		}

		rats[i] = rat
	}

	mischief.rats = rats

	// Put them back in the pool for usage later on
	for _, r := range rats {
		for i := 0; i < mischief.pageConcurrency; i++ {
			r, _ = mischief.ratPool.Get(func() (*rat.Rat, error) {
				return r, nil
			})

			mischief.ratPool.Put(r)
		}
	}

	return nil
}
