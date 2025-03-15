package mischief

import (
	"log/slog"
	"time"
)

// WithExternalBrowsers is an option to use external browsers
// instead of the default browser.
//
// This is useful when you want to use multiple browsers
// to take screenshots concurrently.
//
// Example:
//
//	m := mischief.New(
//		mischief.WithExternalBrowsers([]string{
//			"http://localhost:9222",
//			"http://localhost:9223",
//		}),
//	)
func WithExternalBrowsers(urls []string) MischiefOpt {
	return func(m *Mischief) {
		m.externalBrowser = true
		m.browserUrls = urls
		m.concurrency = len(urls)
	}
}

// WithConcurrency is an option to set the concurrency of the
// screenshotting process.
//
// By default, this will create a pool of <n> browsers to
// take screenshots concurrently.
//
// A browser will be used only once to take a screenshot.
// They should not open multiple tabs.
//
// A browser also has an expiration date of 5 hours.
// After which it will be closed and a new one will be created.
//
// Example:
//
//	m := mischief.New(
//		mischief.WithConcurrency(5),
//	)
func WithConcurrency(c int) MischiefOpt {
	return func(m *Mischief) {
		m.concurrency = c
	}
}

// WithLogger is an option to set the logger of the Mischief instance.
//
// By default, the logger is set to slog.Default().
//
// Example:
//
//	m := mischief.New(
//		mischief.WithLogger(slog.Default()),
//	)
func WithLogger(logger *slog.Logger) MischiefOpt {
	return func(m *Mischief) {
		m.logger = logger
	}
}

// WithBrowserRetakeTimeout is an option to set the timeout
// when taking a browser from the pool.
//
// By default, this is set to 5 seconds.
//
// Example:
//
//	m := mischief.New(
//		mischief.WithBrowserRetakeTimeout(10*time.Second),
//	)
func WithBrowserRetakeTimeout(timeout time.Duration) MischiefOpt {
	return func(m *Mischief) {
		m.browserRetakeTimeout = timeout
	}
}

// WithPageStabilityTimeout is an option to set the timeout
// when waiting for the page to be stable.
//
// By default, this is set to 3 seconds.
//
// Example:
//
//	m := mischief.New(
//		mischief.WithPageStabilityTimeout(5*time.Second),
//	)
func WithPageStabilityTimeout(timeout time.Duration) MischiefOpt {
	return func(m *Mischief) {
		m.pageStabilityTimeout = timeout
	}
}
