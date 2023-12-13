package mistral

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbeddings(t *testing.T) {
	client := NewMistralClientDefault("")
	res, err := client.Embeddings("mistral-embed", []string{"Embed this sentence.", "As well as this one."})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, len(res.Data), 2)
	assert.Len(t, res.Data[0].Embedding, 1024)
}
