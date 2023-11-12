package handlers

import (
	"context"
	"encoding/json"
	"event-tracking-app/database"
	"fmt"
	"net/http"
)

func UpsertEvent(w http.ResponseWriter, r *http.Request) {
	var event database.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := database.Db.Create(&event).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	conn, err := database.GetClickHouseConn()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO events (user_id, type, url, time) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = batch.Append(
		event.UserId,
		event.Type,
		event.Url,
		event.Time,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := batch.Send(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "id: %d, type: %s, url: %s, time: %s", event.UserId, event.Type, event.Url, event.Time)
}

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	var events []database.Event

	if err := database.Db.Find(&events).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	type UserEvent struct {
		Id     int    `json:"event_id"`
		UserId int    `json:"user_id"`
		Type   string `json:"event_type"`
	}

	var userEvents []UserEvent

	for _, event := range events {
		userEvents = append(userEvents, UserEvent{
			Id:     event.Id,
			UserId: event.UserId,
			Type:   event.Type,
		})
	}

	response, err := json.Marshal(userEvents)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	var events []database.Event

	type UserEvent struct {
		Id     int    `json:"event_id"`
		UserId int    `json:"user_id"`
		Type   string `json:"event_type"`
	}

	var requestPayload struct {
		UserId int    `json:"user_id"`
		Type   string `json:"event_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	fmt.Printf("%+v\n", requestPayload)
	if err := database.Db.Where("user_id = ?", requestPayload.UserId).Where("type = ?", requestPayload.Type).Find(&events).Error; err != nil {
		return
	}

	var userEvents []UserEvent

	for _, event := range events {
		userEvents = append(userEvents, UserEvent{
			Id:     event.Id,
			UserId: event.UserId,
			Type:   event.Type,
		})
	}

	response, err := json.Marshal(events)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json)")
	w.Write(response)
}
