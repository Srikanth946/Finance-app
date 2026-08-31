package models

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/genai"
)

func NewGeminiModel(ctx context.Context) model.LLM {

	llm, err := gemini.NewModel(ctx, "gemini-3.5-flash-lite", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return llm
}
