package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb+srv://admin:admin@cluster0.vm2ymw9.mongodb.net/"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Connecting to MongoDB...")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// List databases
	databases, err := client.ListDatabaseNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to list databases: %v", err)
	}

	fmt.Println("Available databases:")
	for _, db := range databases {
		fmt.Printf("  - %s\n", db)
	}
}
