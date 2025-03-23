package rat

import "time"

// WithPagePoolLength is an option to set the length of the page pool
// when creating a new Rat instance.
//
// By default, this is set to 10.
//
// Example:
//
//	rat, err := rat.New(
//		rat.WithPagePoolLength(5),
//	)
func WithPagePoolLength(l int) RatOpt {
	return func(rat *Rat) {
		rat.pagePoolLength = l
	}
}

// WithPageRetakeTimeout is an option to set the timeout
// when taking a page from the pool.
//
// By default, this is set to 5 seconds.
//
// Example:
//
//	rat, err := rat.New(
//		rat.WithPageRetakeTimeout(10*time.Second),
//	)
func WithPageRetakeTimeout(t time.Duration) RatOpt {
	return func(rat *Rat) {
		rat.pageRetakeTimeout = t
	}
}

// WithCreateBrowserFunc is an option to set the function that creates a new browser
// when creating a new Rat instance.
//
// Example:
//
//	rat, err := rat.New(
//		rat.WithCreateBrowserFunc(func(rat *rat.Rat) error {
//			browser, err := rod.New()
//			if err != nil {
//				return err
//			}
//			rat.Browser = browser
//			return nil
//		}),
//	)
func WithCreateBrowserFunc(f func(*Rat) error) RatOpt {
	return func(rat *Rat) {
		rat.createBrowserFunc = f
	}
}
