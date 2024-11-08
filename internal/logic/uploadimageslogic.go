package logic

import (
	"context"
	"github.com/UTC-Six/upload/internal/svc"
	"github.com/UTC-Six/upload/internal/types"
	"github.com/minio/minio-go/v7"
)

type UploadImagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImagesLogic {
	return &UploadImagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadImagesLogic) UploadImages(req types.UploadImagesReq) (*types.UploadResponse, error) {
	bucketName := "mybucket" // 您的 Bucket 名称

	// 检查 Bucket 是否存在
	exists, err := l.svcCtx.MinioClient.BucketExists(l.ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		// 创建 Bucket
		err = l.svcCtx.MinioClient.MakeBucket(l.ctx, bucketName, minio.MakeBucketOptions{
			Region: "us-east-1", // 根据需要更改区域
		})
		if err != nil {
			return nil, err
		}
	}

	// 处理上传逻辑，例如保存文件到 MinIO
	// 您可以根据实际需求完善文件上传逻辑

	return &types.UploadResponse{
		Message: "图片上传成功",
	}, nil
}
