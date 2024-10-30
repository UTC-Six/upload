package logic

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/UTC-Six/upload/internal/svc"
	"github.com/UTC-Six/upload/internal/types"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImagesLogic {
	return &UploadImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadImagesLogic) UploadImages(req *types.UploadImagesReq, r *http.Request) (resp *types.UploadResponse, err error) {
	if err := r.ParseMultipartForm(l.svcCtx.Config.Minio.MaxUploadSize); err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}

	files := r.MultipartForm.File["images"]
	if len(files) > 6 {
		return nil, fmt.Errorf("maximum 6 images allowed")
	}

	for i, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}
		defer src.Close()

		filename := fmt.Sprintf("%s/image_%d%s", req.OrderID, i+1, filepath.Ext(file.Filename))
		_, err = l.svcCtx.MinioClient.PutObject(l.ctx, l.svcCtx.Config.Minio.BucketName, filename, src, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
		if err != nil {
			return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
		}

		l.Logger.Infof("Uploaded image: %s", filename)
	}

	return &types.UploadResponse{Message: "Images uploaded successfully"}, nil
}
