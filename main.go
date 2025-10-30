package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello World!")
	
	if uri := os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatalf("You must set your 'MONGODB_URI' environment variable")
	}
}
