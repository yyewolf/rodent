package rat

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/yyewolf/rodent/pool"
)

type Rat struct {
	pagePoolLength    int
	pageRetakeTimeout time.Duration
	createBrowserFunc func(*Rat) error

	CreatedAt time.Time

	pagePool rod.Pool[rod.Page]

	*rod.Browser
	sync.Mutex
}

type RatOpt func(*Rat)

// New creates a new Rat instance.
// A Rat instance is used to interact with a browser.
//
// Example (and default values):
//
//	rat, err := rat.New(
//		rat.WithPagePoolLength(10),
//		rat.WithPageRetakeTimeout(5*time.Second),
//	)
func New(opts ...RatOpt) (*Rat, error) {
	rat := &Rat{
		CreatedAt: time.Now(),
	}

	var defaultOpts = []RatOpt{
		WithPagePoolLength(10),
		WithPageRetakeTimeout(5 * time.Second),
	}

	opts = append(defaultOpts, opts...)

	for _, opt := range opts {
		opt(rat)
	}

	if rat.createBrowserFunc == nil {
		return nil, fmt.Errorf("create browser function is required")
	}

	err := rat.createBrowserFunc(rat)
	if err != nil {
		return nil, err
	}

	return rat, nil
}

func (rat *Rat) Initialize() error {
	rat.pagePool = rod.NewPagePool(rat.pagePoolLength)
	return nil
}

func (rat *Rat) Close() error {
	pages, err := rat.Pages()
	if err != nil {
		return err
	}

	for _, page := range pages {
		_ = page.Close()
	}

	err = rat.Browser.Close()
	if err != nil {
		return err
	}

	return nil
}

func (rat *Rat) Recreate() error {
	err := rat.Close()
	if err != nil {
		return fmt.Errorf("failed to close rat: %w", err)
	}

	err = rat.createBrowserFunc(rat)
	if err != nil {
		return fmt.Errorf("failed to recreate browser: %w", err)
	}

	return nil
}

func (rat *Rat) createPageFunc() (*rod.Page, error) {
	page, err := rat.Browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (rat *Rat) GetPage() (*rod.Page, error) {
	page, err := pool.GetFromPoolWithTimeout(rat.pagePool, rat.pageRetakeTimeout)
	if err != nil {
		return nil, err
	}

	if page == nil {
		return rat.createPageFunc()
	}

	return page, nil
}

func (rat *Rat) PutPage(page *rod.Page) {
	rat.pagePool.Put(page)
}
