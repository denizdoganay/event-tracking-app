package main

import (
	"context"
	"encoding/json"
	"event-tracking-app/database"
	"event-tracking-app/router"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var wg sync.WaitGroup

func main() {
	database.ConnectDb()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "your_consumer_group_id",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}
	defer closeKafkaConsumer(c)

	wg.Add(1)

	go startKafkaConsumer(c)

	router.Start()
	wg.Wait()
}

func closeKafkaConsumer(c *kafka.Consumer) {
	fmt.Println("Closing Kafka consumer")
	c.Close()
}

func startKafkaConsumer(c *kafka.Consumer) {
	defer wg.Done()
	defer closeKafkaConsumer(c)

	/*
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "your_consumer_group_id",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}
	defer c.Close()
	*/

	err := c.SubscribeTopics([]string{"testing"}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			return
		default:
			ev := c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				processKafkaMessage(string(e.Value))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			}
		}
	}
}

func processKafkaMessage(message string) {
	fmt.Println("Received Kafka message:", message)

	var event database.Event
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		fmt.Printf("Error parsing Kafka message: %v\n", err)
		return
	}

	conn, err := database.GetClickHouseConn()
	if err != nil {
		fmt.Printf("Error connecting to ClickHouse: %v\n", err)
		return
	}
	defer conn.Close()

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO user_events (user_id, type, url, time, param1, param2) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("Error preparing batch: %v\n", err)
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
		fmt.Printf("Error appending to batch: %v\n", err)
		return
	}

	err = batch.Send()
	if err != nil {
		fmt.Printf("Error sending batch: %v\n", err)
		return
	}

	fmt.Println("Processed Kafka message and inserted into ClickHouse")
}