package llm

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

func Create(model_name string, apiKey string, basePrompt string) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	systemPrompt := generateSystemPrompt(basePrompt)

	result, err := client.Models.GenerateContent(
		ctx,
		model_name,
		genai.Text(systemPrompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())

}

func generateSystemPrompt(basePrompt string) string {
	return fmt.Sprintf(`You are a helpful assistant that helps people find information.
	%s`, basePrompt)
}
