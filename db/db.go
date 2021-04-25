package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

// GetDBCollection ...
func GetDBCollection(collectione string) (*mongo.Collection, *mongo.Client, error) {
	//client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("ATLAS")))
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("LOCAL_DB")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	collection := client.Database("RHT").Collection(collectione)
	return collection, client, nil
}
