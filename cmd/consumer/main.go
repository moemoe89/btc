package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxRetries = 3

type Transaction struct {
	UserID   int64     `json:"user_id"`
	Amount   float64   `json:"amount"`
	Datetime time.Time `json:"datetime"`
}

const (
	queueName = "transaction"
)

func main() { //nolint: funlen
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_HOST"))
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	m := make(map[string]int)

	go func() {
		for d := range msgs {
			log.Printf(" [x] Received %s", d.Body)

			m[d.AppId]++

			if m[d.AppId] > maxRetries {
				log.Printf("Received max retry for %s", d.Body)

				delete(m, d.AppId)

				_ = d.Ack(true)
			}

			var trx *Transaction

			err = json.Unmarshal(d.Body, &trx)
			if err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)

				_ = d.Nack(false, true)

				continue
			}

			err = createTransaction(trx)
			if err != nil {
				log.Printf("Failed to create transcation: %v", err)

				_ = d.Nack(false, true)

				continue
			}

			_ = d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func createTransaction(trx *Transaction) error {
	// Dial gRPC server connection.
	conn, err := grpc.Dial(os.Getenv("GRPC_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// Create the client connection.
	client := rpc.NewBTCServiceClient(conn)

	// Call CreateTransaction RPC.
	transaction, err := client.CreateTransaction(context.Background(), &rpc.CreateTransactionRequest{
		UserId:   trx.UserID,
		Datetime: timestamppb.New(trx.Datetime),
		Amount:   trx.Amount,
	})
	if err != nil {
		return err
	}

	log.Printf("transaction created: %v\n", transaction)

	return nil
}
