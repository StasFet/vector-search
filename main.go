package main

import (
	"context"
	"log"
	i "mongo_vector_search/internal"
	"os"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./containers/app.env")
	if err != nil {
		log.Fatalf("error loading env: %v", err)
	}

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

	// create API endpoints
	ginClient := gin.Default()
	ginClient.Use(i.ExtractDatabaseCollectionNames())

	apiGroup := ginClient.Group("/api/:database/:collection/")
	{
		apiGroup.GET("/vectorstore/", i.HandleGetAll(client))
		apiGroup.POST("/vectorstore/", i.HandleInsert(client))
		apiGroup.POST("/vectorsearch/", i.HandleSearch(client))
	}

	port := os.Getenv("PORT")
	if utf8.RuneCountInString(port) == 0 {
		port = "3000"
	}
	ginClient.Run(":" + port)
}
