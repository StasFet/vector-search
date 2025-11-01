package main

import (
	"context"
	"fmt"
	"log"
	i "mongo_vector_search/internal"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	err := godotenv.Load("./containers/app.env")
	if err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	// fmt.Println("Hello World!")
	client, err := i.ConnectToMongo()
	if err != nil {
		log.Fatalf("error connecting to mongo: %v", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("error disconnecting from db: %v", err)
		}
	}()

	if ping := i.PingDB(client); !ping {
		log.Fatalf("error: cannot ping db")
	}

	// TODO: check if this causes an error if the index is already made
	err = client.Database("vector_db_1").CreateCollection(context.TODO(), "coll")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	collection := client.Database("vector_db_1").Collection("coll")
	if err := i.CreateVectorSearchIndex(context.TODO(), collection); err != nil {
		log.Fatalf("error creating vector search index: %v", err)
	}

	// check all the docs
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("error finding documents: %v", err)
	}

	var searchResults []i.VectorDocumentV1
	if err = cursor.All(context.TODO(), &searchResults); err != nil {
		log.Fatalf("error getting all documents from cursor: %v", err)
	}

	for _, result := range searchResults {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res)[:80])
	}
}
