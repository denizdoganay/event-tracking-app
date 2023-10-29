package main

import (
	"event-tracking-app/database"
	"event-tracking-app/router"
)

func main() {
	database.ConnectDb()
	router.Start()
}
