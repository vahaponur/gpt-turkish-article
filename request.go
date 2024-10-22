package gpt_turkish_article

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type ArticleRequest struct {
	Topic         string   `json:"topic"`
	Keywords      []string `json:"keywords"`
	Backlinks     []string `json:"backlinks"`
	BacklinkCount string   `json:"backlinkCount"`
	MinCount      string   `json:"minCount"`
	MaxCount      string   `json:"maxCount"`
}

func (c *Client) requestGpt(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	var chatResp ChatGPTResponse
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&chatResp).
		Post(openAIAPIURL)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("error: %s", resp.String())
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from chatgpt")
	}

	return chatResp.Choices[0].Message.Content, nil
}
func (c *Client) GenerateKeywords(topic string) ([]string, error) {
	prompt := fmt.Sprintf("'%s' konusu ile ilgili 5 adet odak anahtar kelime üret. Bunlar bir makale için kullanılacak"+
		"Küçük harfli olsun aralarında virgül olsun."+
		" Yanıtında anahtar kelimeler dışında HİÇBİR ŞEY yazma, çünkü bir programın içindesin."+
		" Bu yanıt parse edilecek.", topic)
	response, err := c.requestGpt(prompt)
	if err != nil {
		return nil, err
	}

	keywords := strings.Split(response, ",")
	for i := range keywords {
		keywords[i] = strings.TrimSpace(keywords[i])
	}

	return keywords, nil
}
func cleanJSONResponse(response string) (string, error) {
	//regexin gotunu sikim
	re := regexp.MustCompile("```json\n((?s).+?\n)```")
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		if strings.TrimSpace(response)[0] == '{' {
			return strings.TrimSpace(response), nil
		}
		return "", fmt.Errorf("JSON içeriği bulunamadı")
	}
	jsonContent := matches[1]
	jsonContent = strings.TrimSpace(jsonContent)

	return jsonContent, nil
}
func (c *Client) GenerateArticle(topic string, keywords []string, backlink string, backlinkCount int) (Article, error) {
	keywordsStr := strings.Join(keywords, ", ")
	prompt := fmt.Sprintf(`SEO uyumlu bir makale yaz.Content HTML Olsun.Bak HTML Taglerini kullanman çok önemli bu makale direk backende gidecek html taglerini düzgün kullan content içinde.
- Konu: %s
- Anahtar Kelimeler: %s
- Makale içinde %d kez '%s' adresine backlink ver.
minimum 100 kelime.
Sonucu sadece geçerli bir JSON olarak ver. Ek açıklama veya formatlama ekleme. Sadece JSON ver. 
{
  "title": "Makale Başlığı",
  "meta_description": "Meta açıklaması",
  "content": "Makale içeriği"
}`, topic, keywordsStr, backlinkCount, backlink)

	response, err := c.requestGpt(prompt)
	fmt.Println(response)
	if err != nil {
		return Article{}, err
	}

	cleanedResponse, err := cleanJSONResponse(response)
	if err != nil {
		return Article{}, err
	}

	var article Article
	err = json.Unmarshal([]byte(cleanedResponse), &article)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}
func (c *Client) GenerateImageForArticle(title string, keywords []string) (string, error) {
	prompt := title
	if len(keywords) > 0 {
		prompt = fmt.Sprintf("%s başlığı ve bu keywordlere %s uygun blog post resmi üret", title, strings.Join(keywords, ", "))
	}
	reqBody := map[string]interface{}{
		"model":  c.ImageModel,
		"prompt": prompt,
		"n":      1,
		"size":   defaultImageSize,
	}
	var imgResp ImageResponse
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&imgResp).
		Post(openAIImageAPIURL)

	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("API error: %s", resp.String())
	}
	if len(imgResp.Data) == 0 {
		return "", fmt.Errorf("no images generated")
	}
	return imgResp.Data[0].URL, nil
}
