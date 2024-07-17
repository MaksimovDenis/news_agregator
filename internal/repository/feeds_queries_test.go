package repository

import (
	"context"
	"skillfactory/news_agregator/internal/models"
	feedsdb "skillfactory/news_agregator/internal/repository/feedsDB"
	"strconv"
	"testing"

	"math/rand"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgres(t *testing.T) {
	conStr := "postgres://admin:admin@localhost:5432/newsAgregator?sslmode=disable"
	var log zerolog.Logger
	_, err := NewPostgres(context.Background(), conStr, log)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFeedsQueries_StoreFeeds(t *testing.T) {
	conStr := "postgres://admin:admin@localhost:5432/newsAgregator?sslmode=disable"
	var log zerolog.Logger
	db, err := NewPostgres(context.Background(), conStr, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}
	defer db.Close()

	repo := feedsdb.NewFeedsPostgres(db, log)

	type args struct {
		feeds []models.Feeds
	}

	testTable := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				feeds: []models.Feeds{
					{
						Title:   "testTitle",
						Content: "testContent",
						Link:    strconv.Itoa(rand.Intn(1_000_000_000)),
						PubDate: "testPubDate",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.StoreFeeds(testCase.args.feeds)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !testCase.wantErr {
				var title string
				err = db.QueryRow(context.Background(), "SELECT title FROM feeds WHERE title = $1", testCase.args.feeds[0].Title).Scan(&title)
				assert.NoError(t, err)
				assert.Equal(t, "testTitle", title)
			}
		})
	}
}
