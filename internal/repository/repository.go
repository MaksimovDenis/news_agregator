package repository

import (
	"context"
	"skillfactory/news_agregator/internal/models"
	feedsdb "skillfactory/news_agregator/internal/repository/feedsDB"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Feeds interface {
	StoreFeeds(feeds []models.Feeds) error             //Save news to storage
	Feeds(limit int) (feeds []models.Feeds, err error) //Get news from storge
}

type Repository struct {
	Feeds      Feeds
	Log        zerolog.Logger
	PostgresDB *pgxpool.Pool
}

func NewRepository(ctx context.Context, postgresDB *pgxpool.Pool, log zerolog.Logger) *Repository {
	return &Repository{
		Feeds:      feedsdb.NewFeedsPostgres(postgresDB, log),
		Log:        log,
		PostgresDB: postgresDB,
	}
}
