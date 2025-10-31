package main

import (
	"log"
	i "mongo_vector_search/internal"
)

func main() {
	// fmt.Println("Hello World!")
	client, err := i.ConnectToMongo()
	if err != nil {
		log.Fatalf("error connecting to mongo: %v", err)
	}

	if ping := i.PingDB(client); !ping {
		log.Fatalf("error: cannot ping db")
	}

}
