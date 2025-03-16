package api

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-fuego/fuego"
	"github.com/yyewolf/rodent/mischief"
)

// ApiServer is the main struct of the API server.
//
// It is used to start the API server.
type ApiServer struct {
	// port is the port to run the API server on
	port string
	// host is the host to run the API server on
	host string

	// mischief is the Mischief instance to use
	mischief *mischief.Mischief
	// logger is the logger of the API server
	logger *slog.Logger

	// server is the Fuego server
	server *fuego.Server
}

type ApiServerOpt func(*ApiServer)

// New creates a new ApiServer instance.
// An ApiServer instance is used to start the API server.
//
// Example (and default values):
//
//	a := api.New(
//		api.WithHost("0.0.0.0"),
//		api.WithPort("8080"),
//	)
//
// By default, a new Mischief instance is created with default values.
func New(opts ...ApiServerOpt) (*ApiServer, error) {
	var apiServer ApiServer

	var defaultOpts = []ApiServerOpt{
		WithHost("0.0.0.0"),
		WithPort("8080"),
		WithLogger(slog.Default()),
	}

	opts = append(defaultOpts, opts...)

	for _, opt := range opts {
		opt(&apiServer)
	}

	if apiServer.mischief == nil {
		mischief, err := mischief.New()
		if err != nil {
			return nil, errors.Join(ErrCreatingMischiefInstance, err)
		}

		apiServer.mischief = mischief
	}

	apiServer.server = fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf("%s:%s", apiServer.host, apiServer.port)),
		fuego.WithLogHandler(apiServer.logger.Handler()),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				DisableLocalSave: true,
			}),
		),
	)

	apiServer.register()

	return &apiServer, nil
}

// register registers the API server routes.
func (apiServer *ApiServer) register() {
	var repositories = []Repository{
		NewScreenshotRepository(apiServer.mischief, apiServer.logger),
	}

	group := fuego.Group(apiServer.server, "/api")

	for _, repository := range repositories {
		repository.Register(fuego.Group(group, repository.Group()))
	}
}

// Start starts the API server.
func (apiServer *ApiServer) Start() error {
	return apiServer.server.Run()
}
