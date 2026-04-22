package logic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suiran17/go-zero-rag/docservice/internal/chunk"
	"github.com/suiran17/go-zero-rag/pkg/qdrantcli"
	"github.com/suiran17/go-zero-rag/docservice/internal/svc"
	"github.com/suiran17/go-zero-rag/docservice/pb/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDocLogic {
	return &UploadDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadDocLogic) UploadDoc(in *pb.UploadDocReq) (*pb.UploadDocResp, error) {
	docID := uuid.New().String()

	// 1. Split content into chunks.
	chunks := chunk.Split(in.Content)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("empty content after chunking")
	}

	// 2. Embed all chunks via Qwen.
	vectors, err := l.svcCtx.Qwen.Embed(chunks)
	if err != nil {
		return nil, fmt.Errorf("embed: %w", err)
	}

	// 3. Ensure Qdrant collection exists.
	if err := l.svcCtx.Qdrant.EnsureCollection(); err != nil {
		return nil, fmt.Errorf("ensure collection: %w", err)
	}

	// 4. Build points and upsert.
	points := make([]qdrantcli.Point, len(chunks))
	for i, text := range chunks {
		points[i] = qdrantcli.Point{
			ID:       uuid.New().String(),
			Vector:   vectors[i],
			DocID:    docID,
			Text:     text,
			ChunkIdx: i,
		}
	}
	if err := l.svcCtx.Qdrant.Upsert(points); err != nil {
		return nil, fmt.Errorf("upsert: %w", err)
	}

	l.Logger.Infof("uploaded doc %s: %d chunks", docID, len(chunks))
	return &pb.UploadDocResp{
		DocId:      docID,
		ChunkCount: int32(len(chunks)),
	}, nil
}
