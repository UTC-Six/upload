package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/UTC-Six/upload/internal/svc"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type CreateBucketReq struct {
	BucketName string `json:"bucket_name"`
}

type CreateBucketResp struct {
	Message string `json:"message"`
}

func CreateBucketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateBucketReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 检查 Bucket 是否存在
		exists, err := svcCtx.MinioClient.BucketExists(r.Context(), req.BucketName)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if !exists {
			// 创建 Bucket
			err = svcCtx.MinioClient.MakeBucket(r.Context(), req.BucketName, minio.MakeBucketOptions{
				Region: "us-east-1",
			})
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			httpx.OkJson(w, CreateBucketResp{Message: "Bucket 创建成功"})
			return
		}

		httpx.OkJson(w, CreateBucketResp{Message: "Bucket 已存在"})
	}
}
