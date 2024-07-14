package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *API) Feeds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limitStr := vars["limit"]

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	if limit == 0 {
		limit = 10
	}

	feeds, err := api.repository.Feeds.Feeds(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(feeds)
}
