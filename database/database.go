package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Id     int `gorm:"primaryKey; autoIncrement"`
	UserId int
	Type   string
	Url    string
	Time   string
}

var Db *gorm.DB

func ConnectDb() {
	dsn := "host=localhost port=5432 user=postgres password=root dbname=event-tracking-app sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Event{})
	Db = db
}
