package gpt_turkish_article

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
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
