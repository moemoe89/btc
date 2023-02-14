package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Transaction struct {
	UserID   int64     `json:"user_id"`
	Amount   float64   `json:"amount"`
	Datetime time.Time `json:"datetime"`
}

const (
	queueName = "transaction"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	defer func() { _ = ch.Close() }()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	num := rand.Float64() * 100
	num = float64(int(num+0.5)) / 100

	trx := &Transaction{
		UserID:   1,
		Amount:   num,
		Datetime: time.Now(),
	}

	body, err := json.Marshal(trx)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			AppId:       uuid.New().String(),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s", body)
}
