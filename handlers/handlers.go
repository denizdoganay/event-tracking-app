package handlers

import (
	"encoding/json"
	"errors"
	"event-tracking-app/database"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
)

func UpsertEvent(w http.ResponseWriter, r *http.Request) {
	var event database.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := handleDynamicParams(&event); err != nil {
		fmt.Printf("Error handling custom params: %v\n", err)
		return
	}

	/*
		if err := database.Db.Create(&event).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Print("sasa1")
			return
		}
	*/

	fmt.Print("asasa2")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}

	defer p.Close()

	topic := "testing"

	serialized, err := json.Marshal(map[string]interface{}{
		"UserId": event.UserId,
		"Type":   event.Type,
		"Url":    event.Url,
		"Time":   event.Time,
		"Param1": event.Param1,
		"Param2": event.Param2,
	})
	if err != nil {
		fmt.Printf("Error marshaling event to JSON: %v\n", err)
		return
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          serialized,
	}
	p.Produce(msg, nil)

	/*
		conn, err := database.GetClickHouseConn()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer conn.Close()

		if len(event.Details) > 0 {
			if err := handleCustomParams(&event); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO user_events (user_id, type, url, time, param1, param2) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = batch.Append(
			event.UserId,
			event.Type,
			event.Url,
			event.Time,
			event.Param1,
			event.Param2,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := batch.Send(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	//fmt.Fprintf(w, "id: %d, type: %s, url: %s, time: %s", event.UserId, event.Type, event.Url, event.Time)
}

func handleDynamicParams(event *database.Event) error {
	var params database.EventParams

	/*
		columnNameQuery := "SELECT param1, param2 FROM events_params WHERE event = ?"
		result := database.Db.Raw(columnNameQuery, event.Type).Scan(&params)
		if result.Error != nil {
			return result.Error
		}
	*/

	// Attempt to retrieve existing params from events_params table
	result := database.Db.Table("events_params").Where("event = ?", event.Type).First(&params)
	if result.Error != nil {
		// If the record doesn't exist, create a new one with dynamic params
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			params = database.EventParams{
				EventName: event.Type,
				Param1:    "", // Replace with a default value or set based on your logic
				Param2:    "", // Replace with a default value or set based on your logic
			}

			for key := range event.Details {
				if params.Param1 == "" {
					params.Param1 = key
				} else if params.Param2 == "" {
					params.Param2 = key
					break
				}
			}

			if err := database.Db.Table("events_params").Create(&params).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	}

	if paramName, ok := event.Details[params.Param1]; ok {
		event.Param1 = fmt.Sprintf("%v", paramName)
	}

	if paramName, ok := event.Details[params.Param2]; ok {
		event.Param2 = fmt.Sprintf("%v", paramName)
	}

	if err := database.Db.Save(event).Error; err != nil {
		fmt.Print("xddd")
		return err
	}

	return nil
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
