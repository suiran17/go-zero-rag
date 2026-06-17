package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/suiran17/go-zero-rag/aiservice/internal/svc"
	"github.com/suiran17/go-zero-rag/aiservice/pb/pb"
	"github.com/suiran17/go-zero-rag/pkg/llm"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogic {
	return &QueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryLogic) Query(in *pb.QueryReq) (*pb.QueryResp, error) {
	topK := int(in.TopK)
	if topK <= 0 {
		topK = 3
	}

	// 1. Embed the question.
	vecs, err := l.svcCtx.Embedder.Embed([]string{in.Question})
	if err != nil {
		return nil, fmt.Errorf("embed question: %w", err)
	}

	// 2. Retrieve top-k chunks from Qdrant.
	hits, err := l.svcCtx.Qdrant.Search(vecs[0], topK)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	if len(hits) == 0 {
		return &pb.QueryResp{Answer: "知识库中未找到相关内容。", Sources: nil}, nil
	}

	// 3. Build prompt.
	var contextParts []string
	sources := make([]string, len(hits))
	for i, h := range hits {
		contextParts = append(contextParts, fmt.Sprintf("[%d] %s", i+1, h.Text))
		sources[i] = h.Text
	}
	contextText := strings.Join(contextParts, "\n\n")

	messages := []llm.Message{
		{
			Role: "system",
			Content: "你是一个知识问答助手。请根据提供的参考资料回答用户问题，" +
				"答案需忠实于参考资料内容，不要编造。如果参考资料不足以回答，请如实说明。",
		},
		{
			Role: "user",
			Content: fmt.Sprintf("参考资料：\n%s\n\n问题：%s", contextText, in.Question),
		},
	}

	// 4. Call LLM.
	answer, err := l.svcCtx.LLM.Complete(messages)
	if err != nil {
		return nil, fmt.Errorf("llm: %w", err)
	}

	l.Logger.Infof("query=%q hits=%d", in.Question, len(hits))
	return &pb.QueryResp{
		Answer:  answer,
		Sources: sources,
	}, nil
}
