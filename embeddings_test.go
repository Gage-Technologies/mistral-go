package mistral

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmbeddings(t *testing.T) {
	client := NewMistralClient("", "", 3, time.Second*10)
	res, err := client.Embeddings("mistral-embed", []string{"Embed this sentence.", "As well as this one."})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, len(res.Data), 2)
	assert.Len(t, res.Data[0].Embedding, 1024)
}
