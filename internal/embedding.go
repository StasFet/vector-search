package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Vector []float32

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// returns the embedding of s
func GetVectorEmbedding(s string) (*EmbeddingResponse, error) {
	client := &http.Client{}
	api_key := os.Getenv("OPENAI_API_KEY")

	body := map[string]any{
		"model": "text-embedding-3-small",
		"input": s,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result EmbeddingResponse
	if err := json.Unmarshal(bodyData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (e *EmbeddingResponse) GetVector() *[]float32 {
	return &e.Data[0].Embedding
}

func (v Vector) ToBSONBinVector() bson.Binary {
	vectorFloat := bson.NewVector(v)
	return vectorFloat.Binary()
}
