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
	Id      int `gorm:"primaryKey; autoIncrement"`
	UserId  int
	Type    string
	Url     string
	Time    string
	Details map[string]interface{} `gorm:"-"`
	Param1  string
	Param2  string
}

type EventParams struct {
	gorm.Model
	EventName string
	Param1    string
	Param2    string
}

var Db *gorm.DB

func ConnectDb() {
	dsn := "host=localhost port=5432 user=postgres password=root dbname=event-tracking-app sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Event{}, &EventParams{})
	Db = db
}

func GetClickHouseConn() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"ht3qo8ftnx.eu-west-1.aws.clickhouse.cloud:9440"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "y5~2wQRrmN8Bs",
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
