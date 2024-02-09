package router

import (
	"event-tracking-app/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start() {
	router := mux.NewRouter()
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	router.HandleFunc("/upsert", handlers.UpsertEvent).Methods("POST")
	router.HandleFunc("/get-all-events", handlers.GetAllEvents).Methods("GET")
	router.HandleFunc("/get-events", handlers.GetEvents).Methods("POST")

	http.Handle("/", corsHandler)
	http.ListenAndServe(":8081", router)
}
