package svc

import (
	"context"
	"time"

	"github.com/UTC-Six/upload/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

type ServiceContext struct {
	Config      config.Config
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("YOUR-ACCESS-KEY", "YOUR-SECRET-KEY", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("初始化 MinIO 客户端失败: %v", err)
	}

	return &ServiceContext{
		Config:      c,
		MinioClient: minioClient,
	}
}

func (s *ServiceContext) CheckMinIOHealth() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.MinioClient.BucketExists(ctx, s.Config.Minio.BucketName)
	if err != nil {
		logx.Errorf("MinIO health check failed: %v", err)
		return err
	}

	logx.Info("MinIO health check passed")
	return nil
}
