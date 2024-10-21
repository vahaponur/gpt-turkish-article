package gpt_turkish_article

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	defaultModel               = "gpt-4-turbo"
	openAIAPIURL               = "https://api.openai.com/v1/chat/completions"
	openAIImageAPIURL          = "https://api.openai.com/v1/images/generations"
	defaultImageSize           = "512x512"
	defaultImageResponseFormat = "b64_json"
)

type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Article struct {
	Title           string `json:"title"`
	MetaDescription string `json:"meta_description"`
	Content         string `json:"content"`
}
type ImageResponse struct {
	Data []struct {
		B64JSON string `json:"b64_json"`
	} `json:"data"`
}
type Client struct {
	APIKey string
	Model  string
	client *resty.Client
}

func NewGptClient(apiKey string, model ...string) (*Client, error) {
	if len(model) > 1 {
		return &Client{}, fmt.Errorf("error: only one model can be used")
	}
	userModel := defaultModel
	if len(model) != 0 {
		userModel = model[0]
	}
	return &Client{
		APIKey: apiKey,
		Model:  userModel,
		client: resty.New(),
	}, nil
}
