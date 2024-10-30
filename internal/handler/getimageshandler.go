package handler

import (
	"net/http"

	"github.com/UTC-Six/upload/internal/logic"
	"github.com/UTC-Six/upload/internal/svc"
	"github.com/UTC-Six/upload/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetImagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetImagesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetImagesLogic(r.Context(), svcCtx)
		resp, err := l.GetImages(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
