package api

import (
	"log/slog"

	"github.com/yyewolf/rodent/mischief"
)

// WithHost sets the host to run the API server on.
//
// Example:
//
//	a := api.New(
//		api.WithHost("127.0.0.1"),
//	)
func WithHost(host string) ApiServerOpt {
	return func(a *ApiServer) {
		a.host = host
	}
}

// WithPort sets the port to run the API server on.
//
// Example:
//
//	a := api.New(
//		api.WithPort("8080"),
//	)
func WithPort(port string) ApiServerOpt {
	return func(a *ApiServer) {
		a.port = port
	}
}

// WithMischief sets the Mischief instance to use.
//
// Example:
//
//	a := api.New(
//		api.WithMischief(m),
//	)
func WithMischief(mischief *mischief.Mischief) ApiServerOpt {
	return func(a *ApiServer) {
		a.mischief = mischief
	}
}

// WithLogger is an option to set the logger of the API server.
//
// By default, the logger is set to slog.Default().
//
// Example:
//
//	a := api.New(
//		api.WithLogger(slog.Default()),
//	)
func WithLogger(logger *slog.Logger) ApiServerOpt {
	return func(m *ApiServer) {
		m.logger = logger
	}
}
