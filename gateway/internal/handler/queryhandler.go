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

func QueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewQueryLogic(r.Context(), svcCtx)
		resp, err := l.Query(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
