package router

import (
	"net/http"

	"github.com/gorilla/mux"

	handlers "event-tracking-app/handlers"
)

func Start() {
	router := mux.NewRouter()

	router.HandleFunc("/upsert", handlers.UpsertEvent).Methods("POST")

	http.ListenAndServe(":8080", router)
}
