// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/suiran17/go-zero-rag/gateway/internal/logic"
	"github.com/suiran17/go-zero-rag/gateway/internal/svc"
	"github.com/suiran17/go-zero-rag/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadDocHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadDocReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadDocLogic(r.Context(), svcCtx)
		resp, err := l.UploadDoc(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
