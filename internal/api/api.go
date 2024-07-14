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

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./webapp"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./webapp/index.html")
	})

	api.setupEndpoints()

	return api
}

func (api *API) setupEndpoints() {
	api.router.HandleFunc("/api/feeds/{limit}", api.Feeds).Methods("GET")
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
