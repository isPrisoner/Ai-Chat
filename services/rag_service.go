package services

import (
	"AiDemo/config"
	"AiDemo/models"
	"encoding/json"
	"errors"
	"math"
	"time"
)

const (
	// DefaultTopK 默认检索文档数量
	DefaultTopK = 3

	// MinSimilarityThreshold 最小相似度阈值
	MinSimilarityThreshold = 0.0
)

// SaveKnowledge 文档入库，自动切分+批量向量化
func SaveKnowledge(title, content, source, namespace string) ([]*models.Knowledge, error) {
	chunks := ChunkText(content, DefaultChunkSize)

	embeddingModel := GetEmbeddingModelVersion()
	vecs, err := EmbedTextBatch(chunks)
	if err != nil {
		return nil, err
	}

	var results []*models.Knowledge

	for i, chunk := range chunks {
		vecBytes, err := json.Marshal(vecs[i])
		if err != nil {
			return nil, err
		}

		chunkTitle := title
		if len(chunks) > 1 {
			chunkTitle = title + " (片段 " + itoa(i+1) + "/" + itoa(len(chunks)) + ")"
		}

		now := time.Now()
		k := &models.Knowledge{
			ID:             generateKnowledgeID(),
			Title:          chunkTitle,
			Content:        chunk,
			Vector:         string(vecBytes),
			Source:         source,
			Namespace:      namespace,
			EmbeddingModel: embeddingModel,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := config.DB.Create(k).Error; err != nil {
			return nil, err
		}

		results = append(results, k)
	}

	return results, nil
}

// RetrieveRelevantDocsWithScores 返回带相似度分数的结果
func RetrieveRelevantDocsWithScores(query, namespace string, topK int) ([]ScoredDoc, error) {
	if topK <= 0 {
		topK = DefaultTopK
	}

	queryVec, err := EmbedText(query)
	if err != nil {
		return nil, err
	}

	return defaultVectorStore.Search(queryVec, namespace, topK)
}

// RetrieveRelevantDocs 根据查询语句检索相关文档
func RetrieveRelevantDocs(query string, topK int) ([]models.Knowledge, error) {
	scored, err := RetrieveRelevantDocsWithScores(query, "", topK)
	if err != nil {
		return nil, err
	}
	result := make([]models.Knowledge, 0, len(scored))
	for _, s := range scored {
		result = append(result, s.Doc)
	}
	return result, nil
}

// RetrieveRelevantDocsByNamespace 根据命名空间检索相关文档
func RetrieveRelevantDocsByNamespace(query string, namespace string, topK int) ([]models.Knowledge, error) {
	scored, err := RetrieveRelevantDocsWithScores(query, namespace, topK)
	if err != nil {
		return nil, err
	}
	result := make([]models.Knowledge, 0, len(scored))
	for _, s := range scored {
		result = append(result, s.Doc)
	}
	return result, nil
}

func generateKnowledgeID() string {
	return "k_" + itoa(int(time.Now().UnixNano()))
}

// 计算余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	var digits []byte
	for i > 0 {
		d := i % 10
		digits = append([]byte{byte('0' + d)}, digits...)
		i /= 10
	}
	return string(digits)
}

var ErrNoKnowledge = errors.New("no knowledge found")
