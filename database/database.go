package database

import (
	//"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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

func GetClickHouseConn() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"pq0xtq6316.eu-west-1.aws.clickhouse.cloud:9440"},
		Auth: clickhouse.Auth{
			Database: "user_events",
			Username: "default",
			Password: "rY4V8~rS6jcCK",
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "an-example-go-client", Version: "0.1"},
			},
		},

		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		panic(err)
	}
	return conn, nil
}
