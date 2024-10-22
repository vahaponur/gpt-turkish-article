package gpt_turkish_article

import (
	"github.com/go-resty/resty/v2"
)

const (
	defaultModel      = "gpt-4-turbo"
	openAIAPIURL      = "https://api.openai.com/v1/chat/completions"
	openAIImageAPIURL = "https://api.openai.com/v1/images/generations"
	defaultImageSize  = "512x512"
	defaulyImageModel = "dall-e-2"
)

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}
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

type Client struct {
	APIKey     string
	Model      string
	ImageModel string
	ImageSize  string
	client     *resty.Client
}

func NewGptClient(apiKey string) (*Client, error) {
	userModel := defaultModel
	return &Client{
		APIKey:     apiKey,
		Model:      userModel,
		client:     resty.New(),
		ImageModel: defaulyImageModel,
		ImageSize:  defaultImageSize,
	}, nil
}
