package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Dial gRPC server connection.
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Create the client connection.
	client := rpc.NewBTCServiceClient(conn)

	// Call CreateTransaction RPC.
	transaction, err := client.CreateTransaction(context.Background(), &rpc.CreateTransactionRequest{
		UserId:   1,
		Datetime: timestamppb.New(time.Now()),
		Amount:   float32(rand.Int()),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("transaction created: %v\n", transaction)

	// Call ListTransaction RPC.
	transactions, err := client.ListTransaction(context.Background(), &rpc.ListTransactionRequest{
		UserId:        1,
		StartDatetime: timestamppb.New(time.Now().Add(-1 * time.Hour)),
		EndDatetime:   timestamppb.New(time.Now().Add(1 * time.Hour)),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("getting transactions: %v\n", transactions)

	// Call GetUserBalance RPC.
	balance, err := client.GetUserBalance(context.Background(), &rpc.GetUserBalanceRequest{
		UserId: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("getting user balance: %v\n", balance)
}
