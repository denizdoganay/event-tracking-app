package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UpsertEvent(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserId    int    `json:"user_id"`
		EventType string `json:"event_type"`
		Url       string `json:"url"`
		Time      string `json:"time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	fmt.Fprintf(w, "id: %d, type: %s, url: %s, time: %s", payload.UserId, payload.EventType, payload.Url, payload.Time)
}
