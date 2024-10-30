package svc

import (
	"context"
	"time"

	"github.com/UTC-Six/upload/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config      config.Config
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyID, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSSL,
	})
	if err != nil {
		logx.Errorf("Failed to create MinIO client: %v", err)
		panic(err)
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
