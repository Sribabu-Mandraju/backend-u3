package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DBinstance func
func DBinstance() *mongo.Client {
    connectionString := "mongodb+srv://sribabu:63037sribabu@atlascluster.k6u2oy9.mongodb.net/?retryWrites=true&w=majority"

	clientOptions := options.Client().ApplyURI(connectionString);
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	client,err := mongo.Connect(ctx,clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
	}
    fmt.Println("Connected to MongoDB!")

    return client
}

//Client Database instance
var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

    var collection *mongo.Collection = client.Database("u3technologies").Collection(collectionName)

    return collection
}