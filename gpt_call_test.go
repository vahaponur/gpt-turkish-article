package gpt_turkish_article

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"testing"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("error loading env: %s", err))
	}
}
func TestCallGpt(t *testing.T) {
	loadEnv()
	client, err := NewGptClient(os.Getenv("API_KEY"))
	if err != nil {
		t.Error(fmt.Sprintf("error: %s", err))
	}
	resp, err := client.requestGpt("Bana kalileti baskı ile ilgili makale yaz yanıtı json ver baska bisey yazma yanıtına")
	if err != nil {
		t.Error(fmt.Sprintf("error calling gpt: %s", err))
	}
	t.Log(resp)
}
func TestClient_GenerateKeywords(t *testing.T) {
	loadEnv()
	client, err := NewGptClient(os.Getenv("API_KEY"))
	if err != nil {
		t.Error(fmt.Sprintf("error: %s", err))
	}
	resp, err := client.GenerateKeywords("Kartvizit Kalitesinin önemi")
	if err != nil {
		t.Error(fmt.Sprintf("error calling gpt: %s", err))
	}
	t.Log(len(resp))
	for _, keyword := range resp {
		t.Log(keyword)
	}
}

func TestClient_GenerateImageForArticle(t *testing.T) {
	loadEnv()
	client, err := NewGptClient(os.Getenv("API_KEY"))
	if err != nil {
		t.Error(fmt.Sprintf("error: %s", err))
	}
	title := "Düğün davetiyesi"
	keywords, err := client.GenerateKeywords(title)
	if err != nil {
		t.Error(fmt.Sprintf("error:%s", err))
	}
	image, err := client.GenerateImageForArticle("düğün nişan kına", keywords)
	if err != nil {
		t.Error(fmt.Sprintf("error:%s", err))
	}
	t.Log(image)
	base64, err := DownloadImageToBase64(image)
	if err != nil {
		t.Error(fmt.Sprintf("error:%s", err))
	}
	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(fmt.Sprintf("error getting current directory: %s", err))
	}
	outputPath := filepath.Join(currentDir, "tmp", "output.jpg")
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		t.Error(fmt.Sprintf("error creating directory: %s", err))
	}
	err = Base64ToJpeg(base64, outputPath)
	if err != nil {
		t.Error(fmt.Sprintf("error saving image: %s", err))
	}

}
func TestClient_GenerateBulkBlogContent(t *testing.T) {
	loadEnv()
	client, err := NewGptClient(os.Getenv("API_KEY"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	keyword := "Davetiye"
	backlinks := []string{
		"https://lainvito.com/",
		"https://lainvito.com/categories/oval-davetiye",
	}
	topicCount := 3

	t.Logf("Starting bulk generation for keyword: %s", keyword)

	response, err := client.GenerateBulkBlogContent(keyword, backlinks, topicCount)
	if err != nil {
		t.Fatalf("GenerateBulkBlogContent failed: %v", err)
	}

	// Log successful generations
	t.Log("\n=== Successfully Generated Contents ===")
	for i, content := range response.Contents {
		t.Logf("\nContent #%d:", i+1)
		t.Logf("Topic: %s", content.Topic)
		t.Logf("Article Title: %s", content.Article.Title)
		t.Logf("Article Meta Description: %s", content.Article.MetaDescription)
		t.Logf("Article Content: %s", content.Article.Content)
		t.Log("---")
	}

	// Log errors if any
	if len(response.Errors) > 0 {
		t.Log("\n=== Generation Errors ===")
		for i, err := range response.Errors {
			t.Logf("Error #%d: %s", i+1, err)
		}
	}

	// Basic assertions
	if len(response.Contents) == 0 {
		t.Error("No content was generated")
	}

	if len(response.Contents) > topicCount {
		t.Errorf("Generated more contents than requested. Got %d, want <= %d",
			len(response.Contents), topicCount)
	}
}
