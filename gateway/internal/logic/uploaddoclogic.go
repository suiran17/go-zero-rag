// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	docpb "github.com/suiran17/go-zero-rag/docservice/pb/pb"
	"github.com/suiran17/go-zero-rag/gateway/internal/svc"
	"github.com/suiran17/go-zero-rag/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadDocLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDocLogic {
	return &UploadDocLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadDocLogic) UploadDoc(req *types.UploadDocReq) (resp *types.UploadDocResp, err error) {
	reply, err := l.svcCtx.DocRpc.UploadDoc(l.ctx, &docpb.UploadDocReq{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &types.UploadDocResp{
		DocId:      reply.DocId,
		ChunkCount: int(reply.ChunkCount),
	}, nil
}
