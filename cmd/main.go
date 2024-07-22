package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"skillfactory/news_agregator/internal/api"
	"skillfactory/news_agregator/internal/config"
	"skillfactory/news_agregator/internal/repository"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to init config")
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	logger := zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()

	ctx := context.Background()

	dbLog := logger.With().Str("module", "repository").Logger()

	db, err := repository.NewPostgres(ctx, cfg.PgConnString, dbLog)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init database")
	}

	repositories := repository.NewRepository(ctx, db, dbLog)

	apiLog := logger.With().Str("module", "api").Logger()
	fmt.Println(cfg.APIPort)

	APIConfig := &api.Opts{
		Addr:       fmt.Sprintf("localhost:%v", cfg.APIPort),
		Log:        apiLog,
		Repository: repositories,
	}

	server := api.NewAPI(APIConfig)

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start api server")
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGABRT)

	logger.Info().Msg("awaiting signal")

	sig := <-sigs

	log.Info().Str("signal", sig.String()).Msg("signal received")

	server.Stop(context.Background())
	repositories.StopPG()

	logger.Info().Msg("exiting")

}
