package internal

import (
	"time"
)

type VectorDocumentV1 struct {
	Text        string    `bson:"text,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty"`
	EmbeddingV1 []float32 `bson:"embedding"`
}

// creates a full VectorDocumentV1 using `text`
func NewVectorDocumentV1(text string) (*VectorDocumentV1, error) {
	newVec := VectorDocumentV1{}
	newVec.Text = text
	newVec.CreatedAt = time.Now()

	embedRes, err := GetVectorEmbedding(text)
	if err != nil {
		return nil, err
	}
	newVec.EmbeddingV1 = *embedRes.GetVector()

	return &newVec, nil
}