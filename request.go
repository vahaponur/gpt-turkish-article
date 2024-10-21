package gpt_turkish_article

import (
	"fmt"
	"strings"
)

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
