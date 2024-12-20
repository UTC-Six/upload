// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"

	"github.com/UTC-Six/upload/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/images/:order_id",
				Handler: GetImagesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/upload/:order_id",
				Handler: UploadImagesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/create-bucket",
				Handler: CreateBucketHandler(serverCtx),
			},
		},
	)
}
