package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AI struct {
		BaiduOCR struct {
			APIURL string `json:"api_url"`
			Token  string `json:"token"`
		} `json:"baidu_ocr"`

		DeepSeek struct {
			APIKey  string `json:"api_key"`
			BaseURL string `json:"base_url"`
		} `json:"deepseek"`
	} `json:"ai"`

	OSS struct {
		Endpoint        string `json:"endpoint"`
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		BucketName      string `json:"bucket_name"`
		UseSSL          bool   `json:"use_ssl"`
	} `json:"oss"`

	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
}

var Global Config

func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		// 如果文件不存在，我们就不加载，使用默认零值（Mock模式）
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&Global)
}
