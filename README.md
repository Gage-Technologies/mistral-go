# Mistral Go Client

The Mistral Go Client is a comprehensive Golang library designed to interface with the Mistral AI API, providing developers with a robust set of tools to integrate advanced AI-powered features into their applications. This client supports a variety of functionalities, including Chat Completions, Chat Completions Streaming, and Embeddings, allowing for seamless interaction with Mistral's powerful language models.

## Features

- **Chat Completions**: Generate conversational responses and complete dialogue prompts using Mistral's language models.
- **Chat Completions Streaming**: Establish a real-time stream of chat completions, ideal for applications requiring continuous interaction.
- **Embeddings**: Obtain numerical vector representations of text, enabling semantic search, clustering, and other machine learning applications.

## Getting Started

To begin using the Mistral Go Client in your project, ensure you have Go installed on your system. This client library is compatible with Go 1.20 and higher.

### Installation

To install the Mistral Go Client, run the following command:

```bash
go get github.com/gage-technologies/mistral-go
```

### Usage

To use the client in your Go application, you need to import the package and initialize a new client instance with your API key.

```go
package main

import (
	"log"

	"github.com/gage-technologies/mistral-go"
)

func main() {
	// If api key is empty it will load from MISTRAL_API_KEY env var
	client := mistral.NewMistralClientDefault("your-api-key")

	// Example: Using Chat Completions
	chatRes, err := client.Chat("mistral-tiny", []mistral.ChatMessage{{Content: "Hello, world!", Role: mistral.RoleUser}}, nil)
	if err != nil {
		log.Fatalf("Error getting chat completion: %v", err)
	}
	log.Printf("Chat completion: %+v\n", chatRes)

	// Example: Using Chat Completions Stream
	chatResChan, err := client.ChatStream("mistral-tiny", []mistral.ChatMessage{{Content: "Hello, world!", Role: mistral.RoleUser}}, nil)
	if err != nil {
		log.Fatalf("Error getting chat completion stream: %v", err)
	}

	for chatResChunk := range chatResChan {
		if chatResChunk.Error != nil {
			log.Fatalf("Error while streaming response: %v", chatResChunk.Error)
		}
		log.Printf("Chat completion stream part: %+v\n", chatResChunk)
	}

	// Example: Using Embeddings
	embsRes, err := client.Embeddings("mistral-embed", []string{"Embed this sentence.", "As well as this one."})
	if err != nil {
		log.Fatalf("Error getting embeddings: %v", err)
	}

	log.Printf("Embeddings response: %+v\n", embsRes)
}
```

## Documentation

For detailed documentation on the Mistral AI API and the available endpoints, please refer to the [Mistral AI API Documentation](https://docs.mistral.ai).

## Contributing

Contributions are welcome! If you would like to contribute to the project, please fork the repository and submit a pull request with your changes.

## License

The Mistral Go Client is open-sourced software licensed under the [MIT license](LICENSE).

## Support

If you encounter any issues or require assistance, please file an issue on the GitHub repository issue tracker.
