package utils

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB initializes MongoDB connection
func ConnectDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 
	30 * time.Second)

	client, err := mongo.Connect(ctx,options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	DB = client.Database("fiber_assignment")
	log.Println("Connected to MongoDB")
}
