package main

import (
	"context"
	"fmt"
	"log"
	i "mongo_vector_search/internal"

	"github.com/joho/godotenv"
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

	// disconnect at the end
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("error disconnecting from db: %v", err)
		}
	}()

	// test ping
	if ping := i.PingDB(client); !ping {
		log.Fatalf("error: cannot ping db")
	}
	collection := client.Database("vector_db_1").Collection("coll")

	allDocs, err := i.GetAllDocuments(context.TODO(), *collection)
	if err != nil {
		log.Fatalf("error getting all documents: %v", err)
	}

	fmt.Println("All entries:")
	for _, doc := range *allDocs {
		fmt.Printf("\t %s\n", doc.Text)
	}

	searchString := "Stool"

	closest, err := i.VectorSearch(context.TODO(), searchString, *collection)
	if err != nil {
		log.Fatalf("error conducting vector search: %v", err)
	}

	fmt.Printf("Closest text to \"%s\": \n%s\n", searchString, closest)
}
