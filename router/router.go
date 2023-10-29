package router

import (
	"event-tracking-app/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func Start() {
	router := mux.NewRouter()

	router.HandleFunc("/upsert", handlers.UpsertEvent).Methods("POST")
	router.HandleFunc("/get-events", handlers.GetAllEvents).Methods("GET")

	http.ListenAndServe(":8080", router)
}
