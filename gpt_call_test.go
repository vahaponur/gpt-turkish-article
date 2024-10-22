package gpt_turkish_article

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"testing"
	"time"
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
func TestClient_GenerateArticle(t *testing.T) {
	loadEnv()
	client, err := NewGptClient(os.Getenv("API_KEY"))
	if err != nil {
		t.Error(fmt.Sprintf("error: %s", err))
	}
	topic := "Promosyon Ürünleri(diknot,bloknot,kalem,takvim vs...)"
	resp, err := client.GenerateKeywords(topic)
	if err != nil {
		t.Error(fmt.Sprintf("error calling gpt: %s", err))
	}
	article, err := client.GenerateArticle(topic, resp, "https://matbaago.com/promosyon", 3)
	if err != nil {
		t.Error(fmt.Sprintf("error calling gpt: %s", err))
	}
	t.Logf(article.Title)
	t.Logf(article.Content)
	t.Logf(article.MetaDescription)
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
func writeToTimestampedFile(content string, prefix string) (string, error) {
	timestampsDir := "timestamps"
	if err := os.MkdirAll(timestampsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create timestamps directory: %w", err)
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", prefix, timestamp)

	fullPath := filepath.Join(timestampsDir, filename)
	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return fullPath, nil
}
