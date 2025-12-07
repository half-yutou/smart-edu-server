package ocr

import "fmt"

type mockClient struct{}

func NewMockClient() Client {
	return &mockClient{}
}

func (c *mockClient) RecognizeBasic(imageURL string) (string, error) {
	return fmt.Sprintf("[Mock OCR] Recognized text from %s", imageURL), nil
}