package api

import (
	"encoding/json"
	"fmt"
	"os"
	"skillfactory/news_agregator/internal/models"
	"skillfactory/news_agregator/pkg/rss"
	"time"
)

type configRss struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func (hdl *API) StartParseUrl() {

	reader, err := os.ReadFile("./cmd/config.json")
	if err != nil {
		hdl.l.Fatal().Err(err).Msg("failed to read config RSS")
	}

	var config configRss

	err = json.Unmarshal(reader, &config)
	if err != nil {
		hdl.l.Fatal().Err(err).Msg("failed to unmarshal config RSS")
	}

	chFeeds := make(chan []models.Feeds)
	chErrors := make(chan error)

	for _, url := range config.URLS {
		go parseURL(url, chFeeds, chErrors, config.Period)
	}

	go func() {
		for feeds := range chFeeds {
			hdl.repository.Feeds.StoreFeeds(feeds)
		}
	}()

	go func() {
		for err := range chErrors {
			hdl.l.Fatal().Err(err)
		}
	}()
}

func parseURL(url string, feeds chan<- []models.Feeds, errs chan<- error, period int) {
	ticker := time.NewTicker(time.Minute * time.Duration(period))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			news, err := rss.Parse(url)
			fmt.Printf("Полчены новости по ссылке %v\n", url)
			if err != nil {
				errs <- err
				continue
			}
			feeds <- news
		}
	}
}
