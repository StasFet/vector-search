package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
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
				Path:          "embedding",
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

func GetAllDocuments(ctx context.Context, coll mongo.Collection) (*[]VectorDocumentV1, error) {
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var searchResults []VectorDocumentV1
	if err = cursor.All(context.TODO(), &searchResults); err != nil {
		return nil, err
	}

	return &searchResults, nil
}

func VectorSearch(ctx context.Context, text string, coll mongo.Collection) (string, error) {
	embedRes, err := GetVectorEmbedding(text)
	if err != nil {
		return "N/A", err
	}

	queryVector := BSONBinVector(embedRes.GetVector())

	vectorSearchStage := bson.D{
		{"$vectorSearch", bson.D{
			{"index", indexName},
			{"path", "embedding"},
			{"queryVector", queryVector},
			{"numCandidates", 150},
			{"limit", 1},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"text", 1},
			{"created_at", 1},
			{"embedding", 1},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{vectorSearchStage, projectStage})
	if err != nil {
		return "N/A", err
	}

	var results []VectorDocumentV1
	if err = cursor.All(ctx, &results); err != nil {
		return "N/A", err
	}

	return string(results[0].Text), nil
}
