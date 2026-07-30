package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	milvusindexer "github.com/cloudwego/eino-ext/components/indexer/milvus"
	milvusretriever "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/spf13/viper"
)

// 集合字段定义（与 milvusFloatRow / 自定义 converter 严格对应）
const (
	fieldID     = "id"
	fieldText   = "text"
	fieldVector = "vector"
	fieldSource = "source"
)

// milvusFloatRow 写入 milvus 的行结构，字段名需与自定义 Fields 对齐
type milvusFloatRow struct {
	ID     string    `json:"id" milvus:"name:id"`
	Text   string    `json:"text" milvus:"name:text"`
	Vector []float32 `json:"vector" milvus:"name:vector"`
	Source string    `json:"source" milvus:"name:source"`
}

// AgentService 基于 Eino 框架的知识库服务，组合 Embedding + Indexer + Retriever
type AgentService struct {
	embedder  embedding.Embedder
	indexer   *milvusindexer.Indexer
	retriever *milvusretriever.Retriever

	chunkSize      int
	chunkOverlap   int
	defaultTopK    int
}

var (
	agentOnce     sync.Once
	agentService  *AgentService
	agentInitErr  error
)

// NewAgentService 创建知识库服务（惰性初始化：首次调用时建立 Milvus 连接与组件）。
//
// 注意：Milvus / OpenAI 不可用时会返回 error，但不影响其它业务。
// 调用方（handler）应捕获 error 并返回 500。
func NewAgentService() *AgentService {
	agentOnce.Do(func() {
		agentService, agentInitErr = buildAgentService()
	})
	return agentService
}

// InitErr 返回惰性初始化过程中发生的错误（若有）
func (s *AgentService) InitErr() error {
	return agentInitErr
}

func buildAgentService() (*AgentService, error) {
	ctx := context.Background()

	// 1. 创建 OpenAI Embedder
	emb, err := newEmbedder()
	if err != nil {
		return nil, fmt.Errorf("创建 embedder 失败: %w", err)
	}

	// 2. 创建 Milvus client
	cli, err := newMilvusClient()
	if err != nil {
		return nil, fmt.Errorf("创建 milvus client 失败: %w", err)
	}

	dim := viper.GetInt("agent.embedding.dimensions")
	if dim <= 0 {
		dim = 1536
	}
	collection := viper.GetString("agent.milvus.collection")
	if collection == "" {
		collection = "wechat_doc_kb"
	}

	// 3. 创建 Indexer（FloatVector + COSINE，适配 OpenAI float 向量）
	indexer, err := milvusindexer.NewIndexer(ctx, &milvusindexer.IndexerConfig{
		Client:     cli,
		Collection: collection,
		MetricType: milvusindexer.COSINE,
		Embedding:  emb,
		Fields:     floatVectorFields(dim),
		// 自定义 converter：把 schema.Document 转为 milvusFloatRow
		DocumentConverter: func(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]interface{}, error) {
			rows := make([]interface{}, 0, len(docs))
			for idx, doc := range docs {
				fv := make([]float32, len(vectors[idx]))
				for i, x := range vectors[idx] {
					fv[i] = float32(x)
				}
				source, _ := doc.MetaData["source"].(string)
				rows = append(rows, &milvusFloatRow{
					ID:     doc.ID,
					Text:   doc.Content,
					Vector: fv,
					Source: source,
				})
			}
			return rows, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建 indexer 失败: %w", err)
	}

	// 4. 创建 Retriever
	topK := viper.GetInt("agent.search.top_k")
	if topK <= 0 {
		topK = 5
	}
	scoreThreshold := viper.GetFloat64("agent.search.score_threshold")

	retriever, err := milvusretriever.NewRetriever(ctx, &milvusretriever.RetrieverConfig{
		Client:         cli,
		Collection:     collection,
		VectorField:    fieldVector,
		OutputFields:   []string{fieldID, fieldText, fieldSource},
		MetricType:     entity.COSINE,
		TopK:           topK,
		ScoreThreshold: scoreThreshold,
		Embedding:      emb,
		// 自定义 DocumentConverter：把搜索结果映射为 schema.Document
		// （默认 converter 只识别 id/content/metadata，而本集合使用 text/source 字段）
		DocumentConverter: func(ctx context.Context, sr client.SearchResult) ([]*schema.Document, error) {
			count := sr.ResultCount
			if sr.IDs != nil {
				count = sr.IDs.Len()
			}
			docs := make([]*schema.Document, 0, count)

			// 预取字段列，避免在循环内反复遍历
			var (
				idCol, textCol, sourceCol entity.Column
			)
			for _, f := range sr.Fields {
				switch f.Name() {
				case fieldID:
					idCol = f
				case fieldText:
					textCol = f
				case fieldSource:
					sourceCol = f
				}
			}

			for i := 0; i < count; i++ {
				doc := &schema.Document{MetaData: make(map[string]any)}
				if idCol != nil {
					if v, err := idCol.GetAsString(i); err == nil {
						doc.ID = v
					}
				} else if sr.IDs != nil {
					if v, err := sr.IDs.GetAsString(i); err == nil {
						doc.ID = v
					}
				}
				if textCol != nil {
					if v, err := textCol.GetAsString(i); err == nil {
						doc.Content = v
					}
				}
				if sourceCol != nil {
					if v, err := sourceCol.GetAsString(i); err == nil {
						doc.MetaData[fieldSource] = v
					}
				}
				// COSINE 度量下 Scores 为相似度，越大越相关
				if i < len(sr.Scores) {
					doc.WithScore(float64(sr.Scores[i]))
				}
				docs = append(docs, doc)
			}
			return docs, nil
		},
		// 自定义 VectorConverter：把查询向量转为 FloatVector（默认 converter 走 BinaryVector，与索引类型不符）
		VectorConverter: func(ctx context.Context, vectors [][]float64) ([]entity.Vector, error) {
			vec := make([]entity.Vector, 0, len(vectors))
			for _, v := range vectors {
				fv := make([]float32, len(v))
				for i, x := range v {
					fv[i] = float32(x)
				}
				vec = append(vec, entity.FloatVector(fv))
			}
			return vec, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建 retriever 失败: %w", err)
	}

	chunkSize := viper.GetInt("agent.chunk.size")
	if chunkSize <= 0 {
		chunkSize = 500
	}
	chunkOverlap := viper.GetInt("agent.chunk.overlap")
	if chunkOverlap < 0 || chunkOverlap >= chunkSize {
		chunkOverlap = 50
	}

	return &AgentService{
		embedder:     emb,
		indexer:      indexer,
		retriever:    retriever,
		chunkSize:    chunkSize,
		chunkOverlap: chunkOverlap,
		defaultTopK:  topK,
	}, nil
}

// newEmbedder 从配置创建 OpenAI Embedder
func newEmbedder() (embedding.Embedder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := &openai.EmbeddingConfig{
		APIKey: viper.GetString("agent.embedding.api_key"),
		Model:  viper.GetString("agent.embedding.model"),
	}
	if baseURL := viper.GetString("agent.embedding.base_url"); baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if model := cfg.Model; model == "" {
		cfg.Model = "text-embedding-3-small"
	}
	if dim := viper.GetInt("agent.embedding.dimensions"); dim > 0 {
		cfg.Dimensions = &dim
	}
	cfg.Timeout = 30 * time.Second

	emb, err := openai.NewEmbedder(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return emb, nil
}

// newMilvusClient 从配置创建 Milvus client
func newMilvusClient() (client.Client, error) {
	addr := viper.GetString("agent.milvus.addr")
	if addr == "" {
		return nil, fmt.Errorf("agent.milvus.addr 未配置")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		Username: viper.GetString("agent.milvus.username"),
		Password: viper.GetString("agent.milvus.password"),
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// floatVectorFields 定义集合 schema：id(主键) / text / vector(FloatVector) / source
func floatVectorFields(dim int) []*entity.Field {
	return []*entity.Field{
		entity.NewField().
			WithName(fieldID).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256).
			WithIsPrimaryKey(true),
		entity.NewField().
			WithName(fieldText).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(8192),
		entity.NewField().
			WithName(fieldVector).
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(dim)),
		entity.NewField().
			WithName(fieldSource).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(512),
	}
}
