package logic

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/UTC-Six/upload/internal/svc"
	"github.com/UTC-Six/upload/internal/types"
	"github.com/minio/minio-go/v7"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImagesLogic {
	return &GetImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetImagesLogic) GetImages(req *types.GetImagesReq) (resp *types.GetImagesResponse, err error) {
	prefix := req.OrderID + "/"
	objectCh := l.svcCtx.MinioClient.ListObjects(l.ctx, l.svcCtx.Config.Minio.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var images []string
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}
		images = append(images, filepath.Join(req.OrderID, filepath.Base(object.Key)))
	}

	return &types.GetImagesResponse{Images: images}, nil
}
