package router

import (
	"event-tracking-app/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start() {
	fmt.Print("2")
	router := mux.NewRouter()
	fmt.Print("3")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)
	fmt.Print("4")

	router.HandleFunc("/upsert", handlers.UpsertEvent).Methods("POST")
	fmt.Print("5")
	router.HandleFunc("/get-all-events", handlers.GetAllEvents).Methods("GET")
	fmt.Print("6")
	router.HandleFunc("/get-events", handlers.GetEvents).Methods("POST")
	fmt.Print("7")

	http.Handle("/", corsHandler)
	fmt.Print("8")
	http.ListenAndServe(":8080", router)
	fmt.Print("9")
}
