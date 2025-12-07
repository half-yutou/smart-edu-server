package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type baiduClient struct {
	apiURL string
	token  string
}

func NewBaiduClient(apiURL, token string) Client {
	return &baiduClient{
		apiURL: apiURL,
		token:  token,
	}
}

type ocrRequest struct {
	File                      string `json:"file"`
	FileType                  int    `json:"fileType"`
	UseDocOrientationClassify bool   `json:"useDocOrientationClassify"`
	UseDocUnwarping           bool   `json:"useDocUnwarping"`
	UseTextlineOrientation    bool   `json:"useTextlineOrientation"`
}

type ocrResponse struct {
	Result struct {
		OCRResults []struct {
			PrunedResult struct {
				RecTexts []string `json:"rec_texts"`
			} `json:"prunedResult"`
		} `json:"ocrResults"`
	} `json:"result"`
}

func (c *baiduClient) RecognizeBasic(imageURL string) (string, error) {
	// 1. 下载图片
	var imgBytes []byte
	var err error

	if strings.HasPrefix(imageURL, "http") {
		resp, err := http.Get(imageURL)
		if err != nil {
			return "", fmt.Errorf("download image failed: %v", err)
		}
		defer resp.Body.Close()
		imgBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("read image body failed: %v", err)
		}
	} else {
		return "", fmt.Errorf("unsupported image url scheme: %s", imageURL)
	}

	// 2. Base64 编码
	b64Str := base64.StdEncoding.EncodeToString(imgBytes)

	// 3. 构造请求
	payload := ocrRequest{
		File:                      b64Str,
		FileType:                  1,
		UseDocOrientationClassify: false,
		UseDocUnwarping:           false,
		UseTextlineOrientation:    false,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ocr api request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ocr api error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 4. 解析结果
	var ocrResp ocrResponse
	if err := json.Unmarshal(bodyBytes, &ocrResp); err != nil {
		return "", fmt.Errorf("parse ocr response failed: %v", err)
	}

	// 5. 提取文本
	var fullText strings.Builder
	for _, res := range ocrResp.Result.OCRResults {
		for _, line := range res.PrunedResult.RecTexts {
			fullText.WriteString(line)
			fullText.WriteString("\n")
		}
	}

	return strings.TrimSpace(fullText.String()), nil
}
