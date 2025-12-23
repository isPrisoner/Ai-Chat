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

// SaveKnowledge 文档入库，自动进行embedding和切分
func SaveKnowledge(title, content, source, namespace string) ([]*models.Knowledge, error) {
	chunks := ChunkText(content, DefaultChunkSize)

	var results []*models.Knowledge
	embeddingModel := GetEmbeddingModelVersion()

	for i, chunk := range chunks {
		vec, err := EmbedText(chunk)
		if err != nil {
			return nil, err
		}

		vecBytes, err := json.Marshal(vec)
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

// RetrieveRelevantDocs 根据查询语句检索相关文档
func RetrieveRelevantDocs(query string, topK int) ([]models.Knowledge, error) {
	if topK <= 0 {
		topK = DefaultTopK
	}

	queryVec, err := EmbedText(query)
	if err != nil {
		return nil, err
	}

	var all []models.Knowledge
	if err := config.DB.Find(&all).Error; err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return []models.Knowledge{}, nil
	}

	type scored struct {
		Doc   models.Knowledge
		Score float64
	}

	var scoredDocs []scored
	for _, d := range all {
		if d.Vector == "" {
			continue
		}
		var vec []float64
		if err := json.Unmarshal([]byte(d.Vector), &vec); err != nil {
			continue
		}
		s := cosineSimilarity(queryVec, vec)
		if s >= MinSimilarityThreshold {
			scoredDocs = append(scoredDocs, scored{Doc: d, Score: s})
		}
	}

	if len(scoredDocs) == 0 {
		return []models.Knowledge{}, nil
	}

	// 排序并截取TopK
	for i := 0; i < len(scoredDocs)-1; i++ {
		maxIdx := i
		for j := i + 1; j < len(scoredDocs); j++ {
			if scoredDocs[j].Score > scoredDocs[maxIdx].Score {
				maxIdx = j
			}
		}
		scoredDocs[i], scoredDocs[maxIdx] = scoredDocs[maxIdx], scoredDocs[i]
	}

	limit := topK
	if len(scoredDocs) < topK {
		limit = len(scoredDocs)
	}

	result := make([]models.Knowledge, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, scoredDocs[i].Doc)
	}
	return result, nil
}

// RetrieveRelevantDocsByNamespace 根据命名空间检索相关文档
func RetrieveRelevantDocsByNamespace(query string, namespace string, topK int) ([]models.Knowledge, error) {
	if topK <= 0 {
		topK = DefaultTopK
	}

	queryVec, err := EmbedText(query)
	if err != nil {
		return nil, err
	}

	var all []models.Knowledge
	queryDB := config.DB
	if namespace != "" {
		queryDB = queryDB.Where("namespace = ?", namespace)
	}
	if err := queryDB.Find(&all).Error; err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return []models.Knowledge{}, nil
	}

	type scored struct {
		Doc   models.Knowledge
		Score float64
	}

	var scoredDocs []scored
	for _, d := range all {
		if d.Vector == "" {
			continue
		}
		var vec []float64
		if err := json.Unmarshal([]byte(d.Vector), &vec); err != nil {
			continue
		}
		s := cosineSimilarity(queryVec, vec)
		if s >= MinSimilarityThreshold {
			scoredDocs = append(scoredDocs, scored{Doc: d, Score: s})
		}
	}

	if len(scoredDocs) == 0 {
		return []models.Knowledge{}, nil
	}

	// 排序并截取TopK
	for i := 0; i < len(scoredDocs)-1; i++ {
		maxIdx := i
		for j := i + 1; j < len(scoredDocs); j++ {
			if scoredDocs[j].Score > scoredDocs[maxIdx].Score {
				maxIdx = j
			}
		}
		scoredDocs[i], scoredDocs[maxIdx] = scoredDocs[maxIdx], scoredDocs[i]
	}

	limit := topK
	if len(scoredDocs) < topK {
		limit = len(scoredDocs)
	}

	result := make([]models.Knowledge, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, scoredDocs[i].Doc)
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
