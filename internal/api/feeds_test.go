package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skillfactory/news_agregator/internal/models"
	"skillfactory/news_agregator/internal/repository"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockFeedsRepository struct{}

func (m *mockFeedsRepository) Feeds(limit int) ([]models.Feeds, error) {
	return []models.Feeds{
		{Id: 1, Title: "Feed 1", Content: "Content 1"},
		{Id: 2, Title: "Feed 2", Content: "Content 2"},
	}, nil
}

func (m *mockFeedsRepository) StoreFeeds([]models.Feeds) error {
	return nil
}

func TestFeedsAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	repo := &mockFeedsRepository{}
	api := &API{
		repository: &repository.Repository{
			Feeds: repo,
		},
	}

	r.GET("/feeds/:limit", api.Feeds)

	limit := 5
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/feeds/"+strconv.Itoa(limit), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))

	var response []models.Feeds
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	expected := []models.Feeds{
		{Id: 1, Title: "Feed 1", Content: "Content 1"},
		{Id: 2, Title: "Feed 2", Content: "Content 2"},
	}

	assert.Equal(t, expected, response)
}
