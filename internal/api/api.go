package api

import (
	"context"
	"net/http"
	"skillfactory/news_agregator/internal/repository"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 //512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Opts struct {
	Addr       string
	Log        zerolog.Logger
	Repository *repository.Repository
}

type API struct {
	l          zerolog.Logger
	server     *http.Server
	router     *mux.Router
	repository *repository.Repository
}

func NewAPI(opts *Opts) *API {
	router := mux.NewRouter()

	api := &API{
		l: opts.Log,
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		router:     router,
		repository: opts.Repository,
	}

	go api.StartParseUrl()

	api.setupEndpoints()

	return api
}

func (api *API) setupEndpoints() {
	api.router.HandleFunc("/news/{limit}", api.Feeds).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./cmd/webapp"))))
}

func (api *API) Serve() error {
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to start api server")
		return err
	}
	return nil
}

func (api *API) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, ServerShutdownDefaultDelay)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to stop api server")
	}
}
