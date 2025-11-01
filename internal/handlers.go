package internal

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InsertRequestContents struct {
	Text string `json:"text"`
}

type SearchRequestContents struct {
	Text   string `json:"text"`
	Amount int    `json:"amount"`
}

func respondJSON(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, map[string]any{"message": message})
}

func HandleInsert(client *mongo.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vectorColl := client.Database(DatabaseName).Collection(CollectionName)

		// extract the text from the request
		reqContents := InsertRequestContents{}
		if err := ctx.BindJSON(&reqContents); err != nil {
			respondJSON(ctx, http.StatusBadRequest, "could not bind json")
			return
		}

		// create a document
		vectorDoc, err := NewVectorDocumentV1(reqContents.Text)
		if err != nil {
			respondJSON(ctx, http.StatusInternalServerError, "could not generate Document")
			return
		}

		// insert the document into the db
		_, err = InsertDocument(context.TODO(), vectorColl, vectorDoc)
		if err != nil {
			respondJSON(ctx, http.StatusInternalServerError, "could not insert document")
			return
		}

		//respond
		ctx.Status(http.StatusOK)
	}
}

// responds with an arbitrary amount of matches for a string
func HandleSearch(client *mongo.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vectorColl := client.Database(DatabaseName).Collection(CollectionName)

		// extract text and amount
		reqContents := SearchRequestContents{}
		if err := ctx.BindJSON(&reqContents); err != nil {
			respondJSON(ctx, http.StatusBadRequest, "could not bind json")
			return
		}

		// complete search
		matches, err := VectorSearch(context.TODO(), reqContents.Text, 1, *vectorColl)
		if err != nil {
			respondJSON(ctx, http.StatusInternalServerError, "could not complete search")
			return
		}

		// respond
		ctx.JSON(http.StatusOK, map[string]any{
			"status":  "success",
			"matches": matches,
		})
	}
}

// responds with all the documents in the collection
func HandleGetAll(client *mongo.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vectorColl := client.Database(DatabaseName).Collection(CollectionName)

		// fetch all the documents
		documents, err := GetAllDocuments(context.TODO(), *vectorColl)
		if err != nil {
			respondJSON(ctx, http.StatusInternalServerError, "could not get all docuemnts")
			return
		}

		// extract text 
		textOnly := []string{}
		for _, document := range *documents {
			textOnly = append(textOnly, document.Text)
		}

		// respond
		ctx.JSON(http.StatusOK, map[string]any{
			"status": "success",
			"documents": textOnly,
		})
	}
}