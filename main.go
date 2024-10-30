package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/UTC-Six/upload/internal/config"
	"github.com/UTC-Six/upload/internal/handler"
	"github.com/UTC-Six/upload/internal/svc"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/image-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleanupOldImages(ctx.MinioClient, c.Minio.BucketName)
		}
	}()

	// 添加 MinIO 健康检查
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
		defer ticker.Stop()

		for range ticker.C {
			if err := ctx.CheckMinIOHealth(); err != nil {
				logx.Errorf("MinIO health check failed: %v", err)
			}
		}
	}()

	server.Start()
}

func cleanupOldImages(minioClient *minio.Client, bucketName string) {
	ctx := context.Background()
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			logx.Errorf("Error listing objects: %v", object.Err)
			continue
		}

		if object.LastModified.Before(oneWeekAgo) {
			err := minioClient.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				logx.Errorf("Failed to remove old object: %s, error: %v", object.Key, err)
			} else {
				logx.Infof("Removed old object: %s", object.Key)
			}
		}
	}
}
