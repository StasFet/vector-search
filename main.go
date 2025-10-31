package main

import (
	"context"
	"fmt"
	"log"
	"os"
	i "mongo_vector_search/internal"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	// fmt.Println("Hello World!")

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatalf("You must set your 'MONGODB_URI' environment variable")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("The deployment has been successfully pinged! The connection is made.")

	// temporary test
	testString := "Hello World!"
	embedding, err := i.GetVectorEmbedding(testString)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Len of embedding of '%s' : %v\n", testString, len(*embedding.GetVector()))

}
