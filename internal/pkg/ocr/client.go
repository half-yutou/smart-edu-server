package ocr

type Client interface {
	// RecognizeBasic 通用文字识别
	RecognizeBasic(imageURL string) (string, error)
}