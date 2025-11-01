package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type vectorDefinitionField struct {
	Type          string `bson:"type"`
	Path          string `bson:"path"`
	NumDimensions int    `bson:"numDimensions"`
	Similarity    string `bson:"similarity"`
	Quantization  string `bson:"quantization"`
}

type vectorDefinition struct {
	Fields []vectorDefinitionField `bson:"fields"`
}

const indexName = "vector_search_index"

func CreateVectorSearchIndex(ctx context.Context, coll *mongo.Collection) error {
	opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")

	vectorSearchIndexModel := mongo.SearchIndexModel{
		Definition: vectorDefinition{
			Fields: []vectorDefinitionField{{
				Type:          "vector",
				Path:          "plot_embedding",
				NumDimensions: 1536,
				Similarity:    "dotProduct",
				Quantization:  "scalar",
			}},
		},
		Options: opts,
	}

	_, err := coll.SearchIndexes().CreateOne(ctx, vectorSearchIndexModel)
	if err != nil {
		return err
	}
	return nil
}

func InsertDocument(ctx context.Context, coll *mongo.Collection, document any) (*mongo.InsertOneResult, error) {
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}
