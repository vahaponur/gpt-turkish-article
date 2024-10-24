package gpt_turkish_article

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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
func (c *Client) GenerateArticle(request ArticleRequest) (Article, error) {
	totalBacklinks, err := strconv.Atoi(request.BacklinkCount)
	if err != nil {
		return Article{}, fmt.Errorf("invalid backlink count: %v", err)
	}
	if len(request.Backlinks) == 0 {
		return Article{}, fmt.Errorf("backlinks cannot be empty")
	}
	if totalBacklinks > len(request.Backlinks)*5 { // bunun üstü de artık ayıp aq ebenin amı
		return Article{}, fmt.Errorf("too many backlinks requested for the given URLs")
	}
	keywordsStr := strings.Join(request.Keywords, ", ")
	backlinksStr := strings.Join(request.Backlinks, ", ")

	prompt := fmt.Sprintf(`SEO uyumlu bir makale yaz. Content HTML olsun. HTML taglerini düzgün kullan content içinde.
- Konu: %s
- Anahtar Kelimeler: %s
- Verilen linkleri (%s) toplam %s kez makale içinde backlink olarak kullan. Linkleri mümkün olduğunca eşit dağıt.
- Makale uzunluğu %s ile %s kelime arası olmalı.

Sonucu sadece geçerli bir JSON olarak ver. Ek açıklama veya formatlama ekleme. Sadece JSON ver.
{
  "title": "Makale Başlığı",
  "meta_description": "Meta açıklaması",
  "content": "Makale içeriği"
}`,
		request.Topic,
		keywordsStr,
		backlinksStr,
		request.BacklinkCount,
		request.MinCount,
		request.MaxCount,
	)

	response, err := c.requestGpt(prompt)
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
func (c *Client) GenerateTopicsFromKeyword(keyword string) ([]string, error) {
	if keyword == "" {
		return nil, fmt.Errorf("keyword cannot be empty")
	}
	prompt := fmt.Sprintf(`'%s' anahtar kelimesi ile ilgili 10 adet blog yazısı konusu üret.
Blog konuları SEO açısından ilgi çekici ve özgün olsun.
Blog başlıkları en az 30, en fazla 60 karakter uzunluğunda olsun.
Yanıtında konular dışında HİÇBİR ŞEY yazma ve her konu ayrı satırda olsun çünkü bu yanıt parse edilecek.`, keyword)

	response, err := c.requestGpt(prompt)
	if err != nil {
		return nil, fmt.Errorf("topic generation failed: %w", err)
	}
	topics := strings.Split(response, "\n")
	var cleanedTopics []string
	for _, topic := range topics {
		topic = strings.TrimSpace(topic)
		if topic != "" {
			cleanedTopics = append(cleanedTopics, topic)
		}
	}
	if len(cleanedTopics) == 0 {
		return nil, fmt.Errorf("no topics generated")
	}
	return cleanedTopics, nil
}
func (c *Client) UltimateGenerate(topic string, backlinks []string) (Article, string, error) {
	if topic == "" {
		return Article{}, "", fmt.Errorf("topic cannot be empty")
	}
	if len(backlinks) == 0 {
		return Article{}, "", fmt.Errorf("backlinks cannot be empty")
	}
	keywords, err := c.GenerateKeywords(topic)
	if err != nil {
		return Article{}, "", fmt.Errorf("keyword generation failed: %w", err)
	}
	//gotumden atiyom claude baba iyidir dedi bu rakamlara
	request := ArticleRequest{
		Topic:         topic,
		Keywords:      keywords,
		Backlinks:     backlinks,
		BacklinkCount: strconv.Itoa(len(backlinks) * 2),
		MinCount:      "800",
		MaxCount:      "1200",
	}
	article, err := c.GenerateArticle(request)
	if err != nil {
		return Article{}, "", fmt.Errorf("article generation failed: %w", err)
	}
	time.Sleep(time.Second * 5)
	imageURL, err := c.GenerateImageForArticle(article.Title, keywords)
	if err != nil {
		return Article{}, "", fmt.Errorf("image generation failed: %w", err)
	}
	base64Image, err := DownloadImageToBase64(imageURL)
	if err != nil {
		return Article{}, "", fmt.Errorf("image download and conversion failed: %w", err)
	}
	return article, base64Image, nil
}

type GeneratedContent struct {
	Topic       string   `json:"topic"`
	Article     Article  `json:"article"`
	ImageBase64 string   `json:"image_base64"`
	Errors      []string `json:"errors,omitempty"`
}

type BulkGenerationResponse struct {
	Contents []GeneratedContent `json:"contents"`
	Errors   []string           `json:"errors,omitempty"`
}

func (c *Client) GenerateBulkBlogContent(keyword string, backlinks []string, topicCount int) (*BulkGenerationResponse, error) {
	if keyword == "" {
		return nil, fmt.Errorf("keyword cannot be empty")
	}
	if len(backlinks) == 0 {
		return nil, fmt.Errorf("backlinks cannot be empty")
	}
	if topicCount <= 0 {
		return nil, fmt.Errorf("topic count must be positive")
	}
	topics, err := c.GenerateTopicsFromKeyword(keyword)
	if err != nil {
		return nil, fmt.Errorf("topic generation failed: %w", err)
	}
	if len(topics) > topicCount {
		topics = topics[:topicCount]
	}
	var contents []GeneratedContent
	var errors []string
	results := make(chan struct {
		content GeneratedContent
		err     error
	}, len(topics))

	for _, topic := range topics {
		go func(t string) {
			article, imageBase64, err := c.UltimateGenerate(t, backlinks)
			result := struct {
				content GeneratedContent
				err     error
			}{
				content: GeneratedContent{
					Topic:       t,
					Article:     article,
					ImageBase64: imageBase64,
				},
				err: err,
			}
			results <- result
		}(topic)
	}

	for i := 0; i < len(topics); i++ {
		result := <-results
		if result.err != nil {
			errMsg := fmt.Sprintf("failed for topic '%s': %v", result.content.Topic, result.err)
			errors = append(errors, errMsg)
			continue
		}
		contents = append(contents, result.content)
	}
	response := &BulkGenerationResponse{
		Contents: contents,
		Errors:   errors,
	}
	if len(contents) == 0 {
		return response, fmt.Errorf("all content generations failed")
	}
	return response, nil
}
