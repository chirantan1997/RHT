package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDBCollection ...
func GetDBCollection(collectione string) (*mongo.Collection, *mongo.Client, error) {
	//client, err := mongo.NewClient(options.Client().ApplyURI(helpers.GetEnvWithKey("ATLAS")))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
