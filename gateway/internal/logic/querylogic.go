// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	aipb "github.com/suiran17/go-zero-rag/aiservice/pb/pb"
	"github.com/suiran17/go-zero-rag/gateway/internal/svc"
	"github.com/suiran17/go-zero-rag/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogic {
	return &QueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryLogic) Query(req *types.QueryReq) (resp *types.QueryResp, err error) {
	reply, err := l.svcCtx.AiRpc.Query(l.ctx, &aipb.QueryReq{
		Question: req.Question,
		TopK:     int32(req.TopK),
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryResp{
		Answer:  reply.Answer,
		Sources: reply.Sources,
	}, nil
}
