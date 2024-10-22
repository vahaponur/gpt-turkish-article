package gpt_turkish_article

import (
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

// Returns jpeg
func DownloadImageToBase64(imageURL string) (string, error) {
	client := resty.New()
	resp, err := client.R().Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("failed to download image, status: %d, response: %s", resp.StatusCode(), resp.String())
	}
	base64Str := base64.StdEncoding.EncodeToString(resp.Body())
	return base64Str, nil
}
func Base64ToJpeg(base64Str string, outputPath string) error {
	imageData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %w", err)
	}
	err = os.WriteFile(outputPath, imageData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}
	return nil
}
