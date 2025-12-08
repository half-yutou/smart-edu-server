package oss

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"smarteduhub/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

func Init() {
	cfg := config.Global.OSS
	// 默认配置 (适配本地 Docker 环境)
	if cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:9000"
		cfg.AccessKeyID = "minio"
		cfg.SecretAccessKey = "minio123"
		cfg.BucketName = "test"
		cfg.UseSSL = false

		// 更新全局配置，以便后续使用
		config.Global.OSS = cfg
	}

	var err error
	client, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Printf("Failed to create minio client: %v", err)
		return
	}

	// 确保 Bucket 存在
	err = client.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), cfg.BucketName)
		if errBucketExists == nil && exists {
			// Bucket 已存在
		} else {
			log.Printf("Failed to create bucket: %v", err)
			return
		}
	} else {
		log.Printf("Successfully created bucket %s", cfg.BucketName)
	}

	// 设置 Bucket 策略为公开只读 (Public Read)
	policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, cfg.BucketName)
	err = client.SetBucketPolicy(context.Background(), cfg.BucketName, policy)
	if err != nil {
		log.Printf("Failed to set bucket policy: %v", err)
	}
}

func UploadFile(file *multipart.FileHeader) (string, error) {
	if client == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 生成文件名: 原始文件名(去除扩展名)_短时间戳.扩展名
	ext := filepath.Ext(file.Filename)
	nameWithoutExt := strings.TrimSuffix(file.Filename, ext)
	// 简单的清理文件名，防止特殊字符问题 (可选)
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")

	// 使用年月日作为前缀目录，避免单个文件夹下文件过多
	dateDir := time.Now().Format("20060102")
	// 使用短时间戳 (毫秒级后6位)
	shortTimestamp := time.Now().UnixMicro() % 1000000

	objectName := fmt.Sprintf("%s/%s_%d%s", dateDir, nameWithoutExt, shortTimestamp, ext)

	contentType := file.Header.Get("Content-Type")
	// 如果 Content-Type 为空，尝试猜测
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传
	_, err = client.PutObject(context.Background(),
		config.Global.OSS.BucketName, objectName, src, file.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		return "", err
	}

	// 构造返回 URL
	protocol := "http"
	if config.Global.OSS.UseSSL {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://%s/%s/%s",
		protocol, config.Global.OSS.Endpoint, config.Global.OSS.BucketName, objectName)
	return url, nil
}
