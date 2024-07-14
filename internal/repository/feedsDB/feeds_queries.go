package feedsdb

import (
	"context"
	"skillfactory/news_agregator/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type FeedsPostgres struct {
	db *pgxpool.Pool
	l  zerolog.Logger
}

func NewFeedsPostgres(db *pgxpool.Pool, log zerolog.Logger) *FeedsPostgres {
	return &FeedsPostgres{
		db: db,
		l:  log,
	}
}

func (f *FeedsPostgres) Feeds(limit int) (feeds []models.Feeds, err error) {
	rows, err := f.db.Query(context.Background(), `
		SELECT 
			id, 
			title,
			content,
			pub_date,
			link
		FROM feeds
		ORDER BY pub_date DESC
		LIMIT $1;
	`, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var f models.Feeds
		err = rows.Scan(
			&f.Id,
			&f.Title,
			&f.Content,
			&f.PubDate,
			&f.Link,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to get feeds from storage")
			return nil, err
		}

		feeds = append(feeds, f)
	}

	return feeds, nil
}

func (f *FeedsPostgres) StoreFeeds(feeds []models.Feeds) error {
	for _, feed := range feeds {
		_, err := f.db.Exec(context.Background(), `
			INSERT INTO feeds (title, content, pub_date, link)
			VALUES ($1, $2, $3, $4)
		`, feed.Title, feed.Content, feed.PubDate, feed.Link)
		if err != nil {
			log.Error().Err(err).Msg("failed to store feeds to db")
			return err
		}
	}

	return nil
}
